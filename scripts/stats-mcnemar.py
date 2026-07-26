#!/usr/bin/env python3
"""Paired significance analysis for ISC-RAG answer runs (issue S1).

The same questions are answered under every condition, so accuracy
comparisons are paired binary data: McNemar's test (exact binomial on the
discordant pairs), with the discordant odds ratio, the risk difference, and
Wilson score intervals on every reported proportion.

Usage:
  uv run --with scipy --with psycopg2-binary scripts/stats-mcnemar.py \
      --runs-tsv experiments/frontier-xl-runs.tsv \
      --correctness-key judge_kimi_frontier_xl_v1 \
      --support-key source_support_kimi_frontier_xl_v1 \
      [--reference facts] [--metric correct|corr_supp]

The TSV columns are: label <TAB> budget <TAB> run_id. Within each budget the
reference label is compared against every other label.
"""
import argparse
import math
import os
import sys
from collections import defaultdict

import psycopg2
from scipy import stats


def wilson(successes: int, n: int, z: float = 1.959964) -> tuple[float, float]:
    if n == 0:
        return (0.0, 0.0)
    p = successes / n
    denom = 1 + z * z / n
    centre = p + z * z / (2 * n)
    margin = z * math.sqrt(p * (1 - p) / n + z * z / (4 * n * n))
    return ((centre - margin) / denom, (centre + margin) / denom)


def load_outcomes(cur, run_id: str, ckey: str, skey: str, metric: str) -> dict[str, bool]:
    cur.execute(
        """
        SELECT q.external_id,
               a.usage->%s->>'verdict' = 'correct' AS correct,
               a.usage->%s->>'verdict' = 'supported' AS supported
        FROM answers a JOIN questions q ON q.id = a.question_id
        WHERE a.experiment_run_id = %s
        """,
        (ckey, skey, run_id),
    )
    out = {}
    for external_id, correct, supported in cur.fetchall():
        ok = bool(correct) if metric == "correct" else bool(correct) and bool(supported)
        out[external_id] = ok
    return out


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--runs-tsv", required=True)
    ap.add_argument("--correctness-key", required=True)
    ap.add_argument("--support-key", required=True)
    ap.add_argument("--reference", default="facts")
    ap.add_argument("--metric", choices=["correct", "corr_supp"], default="correct")
    args = ap.parse_args()

    dsn = os.environ.get("DATABASE_URL")
    if not dsn:
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

    print(f"# McNemar paired analysis — metric: {args.metric}, "
          f"correctness key: {args.correctness_key}\n")
    rows = []

    for budget in sorted(groups, key=lambda b: int(b)):
        runs = dict(groups[budget])
        if args.reference not in runs:
            continue
        ref = load_outcomes(cur, runs[args.reference], args.correctness_key,
                            args.support_key, args.metric)
        for label, run_id in groups[budget]:
            if label == args.reference:
                continue
            other = load_outcomes(cur, run_id, args.correctness_key,
                                  args.support_key, args.metric)
            shared = sorted(set(ref) & set(other))
            if not shared:
                continue
            n = len(shared)
            ref_ok = sum(ref[q] for q in shared)
            oth_ok = sum(other[q] for q in shared)
            b = sum(1 for q in shared if ref[q] and not other[q])
            c = sum(1 for q in shared if other[q] and not ref[q])
            # Exact McNemar: two-sided binomial test on the discordant pairs.
            p = stats.binomtest(min(b, c), b + c, 0.5).pvalue if b + c else 1.0
            odds = (b / c) if c else math.inf
            rd = (ref_ok - oth_ok) / n
            rlo, rhi = wilson(ref_ok, n)
            olo, ohi = wilson(oth_ok, n)
            rows.append((budget, f"{args.reference} vs {label}", ref_ok, oth_ok,
                         n, rlo, rhi, olo, ohi, b, c, p, odds, rd))

    # Family-wise corrections across every comparison in this invocation.
    m = len(rows)
    order = sorted(range(m), key=lambda i: rows[i][11])
    holm = [0.0] * m
    running = 0.0
    for rank, i in enumerate(order):
        running = max(running, min(1.0, (m - rank) * rows[i][11]))
        holm[i] = running
    bh = [0.0] * m
    prev = 1.0
    for rank in range(m - 1, -1, -1):
        i = order[rank]
        prev = min(prev, min(1.0, rows[i][11] * m / (rank + 1)))
        bh[i] = prev

    print("| Budget | Comparison | Ref acc [95% CI] | Other acc [95% CI] | "
          "Discordant b/c | Exact p | Holm p | BH q | Odds ratio | Risk diff |")
    print("|---|---|---|---|---|---|---|---|---|---|")
    for i, (budget, cmp_, ref_ok, oth_ok, n, rlo, rhi, olo, ohi, b, c, p, odds, rd) in enumerate(rows):
        print(f"| {budget} | {cmp_} | "
              f"{ref_ok}/{n} [{rlo:.3f},{rhi:.3f}] | "
              f"{oth_ok}/{n} [{olo:.3f},{ohi:.3f}] | "
              f"{b}/{c} | {p:.2e} | {holm[i]:.2e} | {bh[i]:.2e} | "
              f"{'inf' if odds == math.inf else f'{odds:.2f}'} | {rd:+.3f} |")
    print(f"\n{m} comparisons in the family; Holm controls FWER, "
          "Benjamini-Hochberg controls FDR.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
