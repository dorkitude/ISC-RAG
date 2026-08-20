# Raw-data audit record — 2026-08-20

Every accuracy figure in the ISC-RAG paper is regenerated here from the raw
database and compared, cell by cell, against the tables committed on 2026-07-26.
The audit was requested because the committed tables were reachable only through
`experiment_run_id` values held in `experiments/*.tsv`; the per-question verdicts
behind them had never left Postgres, so the numbers could not be checked
independently of the run that produced them.

**Result: 56 of 56 rows reproduce exactly. Gate closed.**

## Provenance

| Field | Value |
|---|---|
| Snapshot | `20260819-1620_isc_rag_full.dump` (pg_dump custom format, migration 013 applied), released under tag `sqldump-20260819` at commit `9f5c09d` |
| Produced by | Kyle Wild, 2026-08-19 |
| SHA-256 | `8efe7f1d0a4693571593a48f35e01d1cce616e822a7119951506fc216cbff04b` (verified before restore) |
| Restored into | scratch database `isc_rag_audit`, PostgreSQL 17.11 |
| Analysis code | `scripts/stats-mcnemar.py`, run unmodified |
| Audited by | Yusuke Takahashi, 2026-08-20 |

`pgvector` was not installed. The `vector` type is used only by the `embeddings`
table (migration 002), which the audit does not read; that table and its data
were the only objects that failed to restore. No other restore error occurred.

## Database contents

| Table | Rows |
|---|---|
| `questions` | 10,950 |
| `answers` | 28,968 |
| distinct `experiment_run_id` | 152 |

All 56 run IDs referenced by the five committed ledgers (`cr-comparison-runs.tsv`,
`frontier-xl-runs.tsv`, `heldout-cr-runs.tsv`, `heldout-matched-runs.tsv`,
`heldout-xl-runs.tsv`) are present in the snapshot. This check had not been
performed before; it rules out the possibility that a run cited by the paper is
missing from the data of record.

## Comparison

| Committed file | Section | Rows | Result |
|---|---|---|---|
| `frontier-xl-mcnemar-20260726.md` | metric `correct` | 12 | identical |
| `frontier-xl-mcnemar-20260726.md` | metric `corr_supp` | 12 | identical |
| `frontier-xl-mcnemar-20260726.md` | facts vs contextual-retrieval stack | 4 | identical |
| `heldout-500-mcnemar-20260726.md` | matched budgets 256–2048 | 12 | identical |
| `heldout-500-mcnemar-20260726.md` | XL budgets 2048–16384 | 12 | identical |
| `heldout-500-mcnemar-20260726.md` | facts vs contextual-retrieval stack | 4 | identical |
| **Total** | | **56** | **identical** |

Comparison was exact: every cell of every row, including confidence-interval
bounds and p-values as printed. No tolerance was applied.

Two results are worth restating because they had been recorded inconsistently
elsewhere. First, the McNemar test of facts against the full contextual-retrieval
stack on the frozen 100-document sample *was* executed on 2026-07-26 (92/100 vs
88/100, discordant 8/4, exact p = 0.388, risk difference +0.040); a handover note
of 2026-07-27 described it as outstanding, which was incorrect. Second, on the
held-out 500 sample the sign reverses: the full stack reaches 439/499 against
425/499 at 2,048 tokens (p = 0.202) and 436/499 against 417/499 at 16,384
(p = 0.076). Neither difference is significant, but the paper may not claim an
accuracy win over that baseline.

## Method note

The comparison script initially reported a row-count mismatch in
`frontier-xl-mcnemar-20260726.md`. The cause was in the script, not the data: the
file carries three tables under two `#` headings, the third sitting under a `##`
subheading added on request. The parser counted its four rows as part of the
preceding section. The script now treats `##` as a section boundary; the four
rows were compared separately and are identical. No committed number changed.

## Not covered by this audit

- The verdicts recorded in `experiments/i5a-calibration-20260726/*-key.jsonl`
  have not yet been checked against the database. `scripts/kappa.py` uses them as
  the machine side of the calibration, so they must be verified before any human
  label exists. See `verify_calibration_keys.py`.
- Cost figures are not re-derived here; only accuracy comparisons.
- The `embeddings` table was not restored and is out of scope.

## Reproduction

```sh
brew install postgresql@17 uv
brew services start postgresql@17
./20260820_ISCRAG_RawDataAudit_Homebrew_EN_v6.sh path/to/dump
```
