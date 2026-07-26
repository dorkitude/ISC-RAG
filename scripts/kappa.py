#!/usr/bin/env python3
"""Inter-rater agreement for the calibration study (issue I5c).

Two kappas, reported separately:
1. Rater-rater, PRE-adjudication (Cohen's kappa; the independence evidence).
2. Post-adjudication consensus vs the LLM judge (the calibration figure).

Usage:
  uv run --with scikit-learn scripts/kappa.py \
      --rater-a a.csv --rater-b b.csv --key worksheet-key.jsonl \
      [--consensus consensus.csv] [--disagreements out.csv]

Rating CSVs use the worksheet columns (item_id, rater_id, j1_correct, j1_note,
j2_supported, j2_note, ...). Values must be exactly Yes / No / Unclear.
"Unclear" is kept as its own category for the rater-rater kappa and excluded
from the LLM-calibration comparison (the LLM screen is binary).
"""
import argparse
import csv
import json
import sys

from sklearn.metrics import cohen_kappa_score


def load_ratings(path: str) -> dict[str, dict[str, str]]:
    out = {}
    with open(path, newline="", encoding="utf-8") as f:
        for row in csv.DictReader(f):
            item = row["item_id"].strip()
            if not item:
                continue
            out[item] = {k: (row.get(k) or "").strip() for k in
                         ("j1_correct", "j2_supported")}
    return out


def kappa_for(a, b, field):
    items = sorted(set(a) & set(b))
    xs = [a[i][field] for i in items]
    ys = [b[i][field] for i in items]
    pairs = [(x, y) for x, y in zip(xs, ys) if x and y]
    if not pairs:
        return None, 0, []
    xs, ys = zip(*pairs)
    rated = [i for i in items if a[i][field] and b[i][field]]
    disagreements = [i for i in rated if a[i][field] != b[i][field]]
    return cohen_kappa_score(xs, ys), len(pairs), disagreements


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--rater-a", required=True)
    ap.add_argument("--rater-b", required=True)
    ap.add_argument("--key", required=True)
    ap.add_argument("--consensus", help="post-adjudication consensus CSV (same columns)")
    ap.add_argument("--disagreements", help="write disagreement item list CSV here")
    args = ap.parse_args()

    a, b = load_ratings(args.rater_a), load_ratings(args.rater_b)
    key = {}
    for line in open(args.key, encoding="utf-8"):
        rec = json.loads(line)
        key[f"{rec['item']:03d}"] = rec

    print("## 1. Rater-rater agreement (pre-adjudication)")
    all_disagreements = set()
    for field, label in (("j1_correct", "J1 correctness"), ("j2_supported", "J2 source support")):
        k, n, dis = kappa_for(a, b, field)
        all_disagreements.update(dis)
        print(f"- {label}: Cohen's kappa = {k:.3f} over {n} doubly-rated items; "
              f"{len(dis)} disagreement(s)" if k is not None else f"- {label}: no ratings")

    if args.disagreements:
        with open(args.disagreements, "w", newline="", encoding="utf-8") as f:
            w = csv.writer(f)
            w.writerow(["item_id", "a_j1", "b_j1", "a_j2", "b_j2"])
            for i in sorted(all_disagreements):
                w.writerow([i, a[i]["j1_correct"], b[i]["j1_correct"],
                            a[i]["j2_supported"], b[i]["j2_supported"]])
        print(f"- disagreement list written to {args.disagreements}")

    if args.consensus:
        c = load_ratings(args.consensus)
        print("\n## 2. Consensus vs LLM judge (post-adjudication calibration)")
        for field, llm_field, positive in (
            ("j1_correct", "kimi_correctness", "correct"),
            ("j2_supported", "kimi_source_support", "supported"),
        ):
            pairs = []
            for i, rec in c.items():
                human = rec[field]
                if human not in ("Yes", "No") or i not in key:
                    continue  # Unclear excluded: the LLM screen is binary
                llm = "Yes" if key[i][llm_field] == positive else "No"
                pairs.append((human, llm))
            if not pairs:
                print(f"- {field}: no comparable items")
                continue
            hs, ls = zip(*pairs)
            k = cohen_kappa_score(hs, ls)
            agree = sum(1 for h, l in pairs if h == l)
            print(f"- {field} vs {llm_field}: kappa = {k:.3f}, "
                  f"raw agreement {agree}/{len(pairs)}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
