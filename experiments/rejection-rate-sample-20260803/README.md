# Validation-gate replay result

The paper-facing figures are derived from the captured replay observation in
[`observed-result.json`](observed-result.json), not from prose in the project
README. Run:

```sh
python3 experiments/rejection-rate-sample-20260803/summarize.py \
  --check experiments/rejection-rate-sample-20260803/summary.json
```

The deterministic output reports 29 gate rejections among 2,724 candidates
(1.06%) and an estimated ingest cost of USD 0.06378218 per document, reported
as USD 0.064. The cost applies the price snapshot recorded in the observation
to provider-reported prompt and completion tokens.

## Provenance and limits

- The observation is the terminal summary from the read-only replay executed
  on 2026-08-03. It was recovered from the original local execution record;
  its totals match commit `3d65aa5`, which first recorded the replay.
- The sample is the first 20 held-out documents in production order. Neither
  figure is a corpus-wide measurement. The 500-document token and dollar
  values in `summary.json` are explicitly extrapolations.
- The production extraction predated rejection persistence. Migration 013 is
  present in the audited database, but `fact_rejections` contains zero rows.
  Therefore, the 1.06% rate must always be attributed to this replay.
- Re-running `main.go` would make fresh model calls and need not reproduce the
  same candidates. `summarize.py` is deterministic over the captured
  observation and makes no model, database, or network calls.
