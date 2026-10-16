#!/usr/bin/env python3
"""Reference-answer metrics for ISC-RAG answer runs — no LLM in the loop.

Every judgment reported so far comes from Kimi K2.6, which also wrote the
fact payloads. This script scores the same stored answers against the
human-authored `questions.reference_answer` spans using the standard
SQuAD/EnronQA string metrics, so the facts-vs-chunk ordering can be checked
with no model involved and no API cost.

Three metrics, all computed on SQuAD-normalised text (lowercase, articles and
punctuation stripped, whitespace collapsed):

  em        normalised exact match (strict; generative answers score low)
  contains  the normalised reference occurs as a substring of the answer
  f1        token-level F1 between answer and reference

`em` and `contains` are paired binary outcomes, so they get the same test as
scripts/stats-mcnemar.py: exact McNemar (binomial on the discordant pairs),
discordant odds ratio, risk difference and Wilson intervals. `f1` is
continuous, so it gets a paired bootstrap percentile interval on the mean
difference plus a Wilcoxon signed-rank test.

Usage:
  uv run --with scipy --with psycopg2-binary scripts/stats-reference-metrics.py \
      --runs-tsv experiments/heldout-xl-runs.tsv [--reference facts] \
      [--metric em|contains|f1] [--per-run] [--examples 5] [--seed 20261016]

  # verify the scoring functions with no database:
  uv run --with scipy scripts/stats-reference-metrics.py --selftest

The TSV columns are: label <TAB> budget <TAB> run_id. Within each budget the
reference label is compared against every other label.
"""
import argparse
import math
import os
import random
import re
import string
import sys
from collections import Counter, defaultdict

_PUNCT = str.maketrans("", "", string.punctuation)
_ARTICLES = re.compile(r"\b(a|an|the)\b")

CONTINUOUS = {"f1", "precision", "recall"}


def normalise(s: str) -> str:
    """SQuAD v1.1 normalisation: lowercase, drop punctuation and articles."""
    if s is None:
        return ""
    s = s.lower().translate(_PUNCT)
    s = _ARTICLES.sub(" ", s)
    return " ".join(s.split())


def token_prf(pred: str, gold: str) -> tuple[float, float, float]:
    """Token-level precision, recall and F1.

    Measured on this corpus: reference answers average 22.7 words while system
    answers average 6.6-9.5, so every run is recall-capped and F1 rewards
    verbosity independently of correctness. Precision is the length-robust
    direction. All five metrics are always reported.
    """
    p, g = normalise(pred).split(), normalise(gold).split()
    if not p or not g:
        v = float(p == g)
        return (v, v, v)
    common = Counter(p) & Counter(g)
    overlap = sum(common.values())
    if overlap == 0:
        return (0.0, 0.0, 0.0)
    precision = overlap / len(p)
    recall = overlap / len(g)
    return (precision, recall, 2 * precision * recall / (precision + recall))


def token_f1(pred: str, gold: str) -> float:
    return token_prf(pred, gold)[2]


def score(pred: str, gold: str) -> dict:
    np_, ng = normalise(pred), normalise(gold)
    prec, rec, f1 = token_prf(pred, gold)
    return {
        "em": float(np_ == ng),
        "contains": float(bool(ng) and ng in np_),
        "precision": prec,
        "recall": rec,
        "f1": f1,
        "pred_words": float(len(np_.split())),
        "gold_words": float(len(ng.split())),
    }


def wilson(successes: float, n: int, z: float = 1.959964) -> tuple[float, float]:
    if n == 0:
        return (0.0, 0.0)
    p = successes / n
    denom = 1 + z * z / n
    centre = p + z * z / (2 * n)
    margin = z * math.sqrt(p * (1 - p) / n + z * z / (4 * n * n))
    return (max(0.0, (centre - margin) / denom), min(1.0, (centre + margin) / denom))


def load_scores(cur, run_id: str) -> dict[str, dict]:
    """external_id -> {em, contains, f1} for every answered question with a
    non-null, non-empty reference answer."""
    cur.execute(
        """
        SELECT q.external_id, a.text, q.reference_answer
        FROM answers a JOIN questions q ON q.id = a.question_id
        WHERE a.experiment_run_id = %s
        """,
        (run_id,),
    )
    out = {}
    for external_id, text, reference in cur.fetchall():
        if reference is None or not reference.strip():
            continue
        out[external_id] = score(text or "", reference)
    return out


def paired_bootstrap(diffs: list[float], iters: int, rng: random.Random) -> tuple[float, float]:
    n = len(diffs)
    if n == 0:
        return (0.0, 0.0)
    means = []
    for _ in range(iters):
        means.append(sum(diffs[rng.randrange(n)] for _ in range(n)) / n)
    means.sort()
    lo = means[int(0.025 * iters)]
    hi = means[min(iters - 1, int(0.975 * iters))]
    return (lo, hi)


def selftest() -> int:
    """Known SQuAD-style cases. No database, no network."""
    cases = [
        # (pred, gold, em, contains, f1 rounded to 3dp)
        ("About a quarter.", "about a quarter", 1.0, 1.0, 1.000),
        ("The answer is about a quarter.", "About a quarter.", 0.0, 1.0, 0.667),
        ("Michelle Singletary.", "Michelle Singletary", 1.0, 1.0, 1.000),
        # `contains` is a substring test, so a gold prefix matches: deliberate, and
        # the reason `f1` is reported alongside it.
        ("Restaurants and retailers", "restaurants and retail", 0.0, 1.0, 0.667),
        ("Retailers only", "restaurants and retail", 0.0, 0.0, 0.000),
        ("Paris", "London", 0.0, 0.0, 0.000),
        ("", "London", 0.0, 0.0, 0.000),
    ]
    failures = 0
    print("| pred | gold | em | contains | f1 | expected f1 | ok |")
    print("|---|---|---|---|---|---|---|")
    for pred, gold, em, contains, f1 in cases:
        s = score(pred, gold)
        ok = (s["em"] == em and s["contains"] == contains
              and abs(s["f1"] - f1) < 5e-4)
        failures += 0 if ok else 1
        print(f"| {pred!r} | {gold!r} | {s['em']:.0f} | {s['contains']:.0f} | "
              f"{s['f1']:.3f} | {f1:.3f} | {'PASS' if ok else 'FAIL'} |")
    print(f"\n{len(cases) - failures}/{len(cases)} self-test cases pass.")
    return 0 if failures == 0 else 1


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--runs-tsv")
    ap.add_argument("--reference", default="facts")
    ap.add_argument("--metric",
                    choices=["em", "contains", "f1", "precision", "recall"],
                    default="contains")
    ap.add_argument("--per-run", action="store_true",
                    help="also print all five metrics for every run in the ledger")
    ap.add_argument("--examples", type=int, default=0,
                    help="print N answer/reference pairs from the reference run")
    ap.add_argument("--bootstrap", type=int, default=10000)
    ap.add_argument("--seed", type=int, default=20261016)
    ap.add_argument("--selftest", action="store_true")
    args = ap.parse_args()

    if args.selftest:
        return selftest()
    if not args.runs_tsv:
        print("--runs-tsv is required unless --selftest", file=sys.stderr)
        return 2

    import psycopg2
    from scipy import stats

    dsn = os.environ.get("DATABASE_URL")
    if not dsn and os.path.exists(".env"):
        for line in open(".env", encoding="utf-8"):
            if line.startswith("DATABASE_URL="):
                dsn = line.split("=", 1)[1].strip()
    if not dsn:
        print("DATABASE_URL not found", file=sys.stderr)
        return 1

    groups: dict[str, list[tuple[str, str]]] = defaultdict(list)
    for line in open(args.runs_tsv, encoding="utf-8"):
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        label, budget, run_id = line.split("\t")
        groups[budget].append((label, run_id))

    conn = psycopg2.connect(dsn)
    cur = conn.cursor()
    rng = random.Random(args.seed)

    print(f"# Reference-answer metrics — no judge model involved\n")
    print(f"Ledger: `{args.runs_tsv}`  ·  reference label: `{args.reference}`  ·  "
          f"paired metric: `{args.metric}`  ·  seed {args.seed}\n")

    cached: dict[str, dict[str, dict]] = {}

    def get(run_id: str) -> dict[str, dict]:
        if run_id not in cached:
            cached[run_id] = load_scores(cur, run_id)
        return cached[run_id]

    if args.examples:
        ref_runs = [r for b in groups for (l, r) in groups[b] if l == args.reference]
        if ref_runs:
            cur.execute(
                """
                SELECT q.external_id, q.reference_answer, a.text
                FROM answers a JOIN questions q ON q.id = a.question_id
                WHERE a.experiment_run_id = %s AND q.reference_answer IS NOT NULL
                ORDER BY q.external_id LIMIT %s
                """,
                (ref_runs[0], args.examples),
            )
            print("## Sample answer/reference pairs (normalisation sanity check)\n")
            for ext, ref, txt in cur.fetchall():
                s = score(txt or "", ref)
                print(f"- `{ext}`  em={s['em']:.0f} contains={s['contains']:.0f} "
                      f"f1={s['f1']:.3f}\n  - reference: {ref!r}\n  - answer: {txt!r}")
            print()

    if args.per_run:
        print("## Per-run means\n")
        print("| Budget | Label | n | EM | Contains | Prec | Rec | token-F1 "
              "| ans words | ref words |")
        print("|---|---|---|---|---|---|---|---|---|---|")
        for budget in sorted(groups, key=lambda b: int(b)):
            for label, run_id in groups[budget]:
                s = get(run_id)
                n = len(s)
                if n == 0:
                    print(f"| {budget} | {label} | 0 | — | — | — | — | — | — | — |")
                    continue
                mean = lambda k: sum(v[k] for v in s.values()) / n
                print(f"| {budget} | {label} | {n} | {mean('em'):.3f} | "
                      f"{mean('contains'):.3f} | {mean('precision'):.3f} | "
                      f"{mean('recall'):.3f} | {mean('f1'):.3f} | "
                      f"{mean('pred_words'):.1f} | {mean('gold_words'):.1f} |")
        print()

    rows = []
    for budget in sorted(groups, key=lambda b: int(b)):
        runs = dict(groups[budget])
        if args.reference not in runs:
            continue
        ref = get(runs[args.reference])
        for label, run_id in groups[budget]:
            if label == args.reference:
                continue
            other = get(run_id)
            shared = sorted(set(ref) & set(other))
            if not shared:
                continue
            n = len(shared)
            if args.metric in CONTINUOUS:
                mk = args.metric
                diffs = [ref[q][mk] - other[q][mk] for q in shared]
                rm = sum(ref[q][mk] for q in shared) / n
                om = sum(other[q][mk] for q in shared) / n
                lo, hi = paired_bootstrap(diffs, args.bootstrap, rng)
                nz = [d for d in diffs if d != 0.0]
                p = (stats.wilcoxon(nz).pvalue if len(nz) >= 1 else 1.0)
                rows.append((budget, f"{args.reference} vs {label}", n,
                             rm, om, rm - om, lo, hi, p))
            else:
                k = args.metric
                ref_ok = int(sum(ref[q][k] for q in shared))
                oth_ok = int(sum(other[q][k] for q in shared))
                b = sum(1 for q in shared if ref[q][k] and not other[q][k])
                c = sum(1 for q in shared if other[q][k] and not ref[q][k])
                p = stats.binomtest(min(b, c), b + c, 0.5).pvalue if b + c else 1.0
                odds = (b / c) if c else math.inf
                rd = (ref_ok - oth_ok) / n
                rlo, rhi = wilson(ref_ok, n)
                olo, ohi = wilson(oth_ok, n)
                rows.append((budget, f"{args.reference} vs {label}", ref_ok, oth_ok,
                             n, rlo, rhi, olo, ohi, b, c, p, odds, rd))

    if not rows:
        print("No comparable pairs found.", file=sys.stderr)
        return 1

    # Family-wise corrections across every comparison in this invocation,
    # identical in form to scripts/stats-mcnemar.py.
    pidx = 8 if args.metric in CONTINUOUS else 11
    m = len(rows)
    order = sorted(range(m), key=lambda i: rows[i][pidx])
    holm = [0.0] * m
    running = 0.0
    for rank, i in enumerate(order):
        running = max(running, min(1.0, (m - rank) * rows[i][pidx]))
        holm[i] = running
    bh = [0.0] * m
    prev = 1.0
    for rank in range(m - 1, -1, -1):
        i = order[rank]
        prev = min(prev, min(1.0, rows[i][pidx] * m / (rank + 1)))
        bh[i] = prev

    print(f"## Paired comparisons — {args.metric}\n")
    if args.metric in CONTINUOUS:
        print(f"| Budget | Comparison | n | Ref {args.metric} | "
              f"Other {args.metric} | Mean diff "
              "[95% bootstrap] | Wilcoxon p | Holm p | BH q |")
        print("|---|---|---|---|---|---|---|---|---|")
        for i, (budget, cmp_, n, rm, om, d, lo, hi, p) in enumerate(rows):
            print(f"| {budget} | {cmp_} | {n} | {rm:.3f} | {om:.3f} | "
                  f"{d:+.3f} [{lo:+.3f},{hi:+.3f}] | {p:.2e} | "
                  f"{holm[i]:.2e} | {bh[i]:.2e} |")
    else:
        print("| Budget | Comparison | Ref [95% CI] | Other [95% CI] | "
              "Discordant b/c | Exact p | Holm p | BH q | Odds ratio | Risk diff |")
        print("|---|---|---|---|---|---|---|---|---|---|")
        for i, (budget, cmp_, ref_ok, oth_ok, n, rlo, rhi, olo, ohi,
                b, c, p, odds, rd) in enumerate(rows):
            print(f"| {budget} | {cmp_} | "
                  f"{ref_ok}/{n} [{rlo:.3f},{rhi:.3f}] | "
                  f"{oth_ok}/{n} [{olo:.3f},{ohi:.3f}] | "
                  f"{b}/{c} | {p:.2e} | {holm[i]:.2e} | {bh[i]:.2e} | "
                  f"{'inf' if odds == math.inf else f'{odds:.2f}'} | {rd:+.3f} |")
    print(f"\n{m} comparisons in the family; Holm controls FWER, "
          "Benjamini-Hochberg controls FDR.")
    print("\nScored against human-authored `questions.reference_answer`. "
          "No judge model, no API call, no network access.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
