#!/usr/bin/env python3
"""Verify the committed human-calibration key files against the raw database.

The 80-item worksheet and the 10-item practice set each ship a *-key.jsonl that
records the automated screen's verdict for every item. scripts/kappa.py treats
those recorded verdicts as the machine side of the calibration, so if they have
drifted from the database the resulting kappa is meaningless.

This script re-reads every verdict from answers.usage and compares. It must be
run, and must pass, before any human label exists.

Usage:
  DATABASE_URL=postgres:///isc_rag_audit uv run --with psycopg2-binary \
      verify_calibration_keys.py --repo ~/firebase/ISC-RAG
"""
import argparse, json, os, sys
from collections import Counter, defaultdict
import psycopg2

CAL = "experiments/i5a-calibration-20260726"
LEDGER = "experiments/frontier-xl-runs.tsv"
CKEY = "judge_kimi_frontier_xl_v1"
SKEY = "source_support_kimi_frontier_xl_v1"


def load_jsonl(path):
    with open(path, encoding="utf-8") as f:
        return [json.loads(l) for l in f if l.strip()]


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--repo", default=os.path.expanduser("~/firebase/ISC-RAG"))
    args = ap.parse_args()
    repo = os.path.expanduser(args.repo)

    dsn = os.environ.get("DATABASE_URL")
    if not dsn:
        print("DATABASE_URL not set", file=sys.stderr)
        return 1

    # run_id -> (representation, budget) from the committed ledger
    ledger = {}
    with open(os.path.join(repo, LEDGER), encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if line and not line.startswith("#"):
                label, budget, run_id = line.split("\t")
                ledger[run_id] = (label, budget)

    conn = psycopg2.connect(dsn)
    cur = conn.cursor()
    problems = []
    total = 0

    for name, fn in (("worksheet", "worksheet-key.jsonl"), ("practice", "practice-key.jsonl")):
        rows = load_jsonl(os.path.join(repo, CAL, fn))
        print(f"\n== {name}: {len(rows)} items ==")
        cells = Counter()
        verdicts = Counter()
        seen_answers = set()

        for r in rows:
            total += 1
            item = r.get("item")
            aid = r["answer_id"]
            if aid in seen_answers:
                problems.append(f"{name} item {item}: duplicate answer_id {aid}")
            seen_answers.add(aid)

            cur.execute(
                """
                SELECT a.experiment_run_id::text, q.external_id,
                       a.usage->%s->>'verdict', a.usage->%s->>'verdict'
                FROM answers a JOIN questions q ON q.id = a.question_id
                WHERE a.id = %s
                """,
                (CKEY, SKEY, aid),
            )
            got = cur.fetchone()
            if got is None:
                problems.append(f"{name} item {item}: answer_id {aid} not in database")
                continue
            db_run, db_qid, db_correct, db_support = got

            if db_run != r["run_id"]:
                problems.append(f"{name} item {item}: run_id {r['run_id']} vs db {db_run}")
            if db_qid != r["question_external_id"]:
                problems.append(f"{name} item {item}: question {r['question_external_id']} vs db {db_qid}")
            if db_correct != r["kimi_correctness"]:
                problems.append(f"{name} item {item}: correctness {r['kimi_correctness']} vs db {db_correct}")
            if db_support != r["kimi_source_support"]:
                problems.append(f"{name} item {item}: support {r['kimi_source_support']} vs db {db_support}")

            # the recorded representation/budget must match the committed ledger
            led = ledger.get(db_run)
            if led is None:
                problems.append(f"{name} item {item}: run {db_run} absent from {LEDGER}")
            else:
                rep, bud = led
                if r.get("representation") not in (None, rep):
                    problems.append(f"{name} item {item}: representation {r['representation']} vs ledger {rep}")
                if r.get("context_token_budget") not in (None, bud):
                    problems.append(f"{name} item {item}: budget {r['context_token_budget']} vs ledger {bud}")
                cells[(rep, bud)] += 1
            verdicts[(r.get("representation"), r["kimi_correctness"])] += 1

        print("  cells (from the committed ledger):")
        for k, v in sorted(cells.items()):
            print(f"    {k[0]:<10} @ {k[1]:<6} {v:>3}")
        print("  screen verdict balance:")
        for k, v in sorted(verdicts.items(), key=str):
            print(f"    {str(k[0]):<10} {k[1]:<10} {v:>3}")

    # question overlap between the two instruments
    ws = {r["question_external_id"] for r in load_jsonl(os.path.join(repo, CAL, "worksheet-key.jsonl"))}
    pr = {r["question_external_id"] for r in load_jsonl(os.path.join(repo, CAL, "practice-key.jsonl"))}
    print(f"\n== instruments ==")
    print(f"  questions shared by practice and worksheet: {len(ws & pr)} of {len(pr)}")
    if ws & pr:
        print("   ", ", ".join(sorted(ws & pr)))

    print("\n" + "=" * 70)
    print(f"checked {total} items")
    if problems:
        print(f"VERIFY: FAIL — {len(problems)} problem(s)")
        for p in problems:
            print("  -", p)
        return 1
    print("VERIFY: PASS — every recorded verdict matches the database")
    return 0


if __name__ == "__main__":
    sys.exit(main())
