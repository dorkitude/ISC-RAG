# Held-out 500 protocol and progress (I6a / I6b)

**Status:** protocol frozen 2026-07-26 BEFORE any pipeline execution on these
documents, per the held-out discipline. This document is the running record;
the Progress section is appended as stages complete. Nothing in the pipeline
may be tuned against these documents; deltas versus development numbers will
be reported as observed, including any that shrink.

## Sample selection (frozen)

500 documents, deterministically selected, sized 5× the minimum spec to
support tighter intervals:

```sql
WITH touched AS (
  SELECT document_id FROM source_units
  UNION SELECT document_id FROM gold_evidence
  UNION SELECT document_id FROM facts
  UNION SELECT document_id FROM study_two_runs
  UNION SELECT document_id FROM fact_rephrasing_study_runs
)
SELECT d.external_id FROM documents d
WHERE d.corpus_split='test' AND d.id NOT IN (SELECT document_id FROM touched)
ORDER BY md5('heldout-500-20260726:' || d.external_id)
LIMIT 500;
```

The exclusion set covers every document ever touched by any representation,
question, fact extraction, or isolated study. Result committed as
[`examples/heldout-500-20260726.txt`](../examples/heldout-500-20260726.txt).

**Verification:** 500/500 rows; pairwise overlap with
`fact-development-sample-100-20260725.txt` = 0 and with
`fact-development-depth-pilot-5-20260725.txt` = 0.

## Frozen pipeline configuration

Identical to the frozen-100 development pipeline:

- **Representations:** `build --representation all --version heldout500-v1
  --document-list examples/heldout-500-20260726.txt --split test` (defaults:
  4,000-byte units, 400-byte overlap, 2,000-byte semantic minimum, 10,000-byte
  fact inputs).
- **Fact extraction:** `extract facts --version heldout500-turnwise-opinion-v1
  --input-version heldout500-v1 --document-list … --workers 8` with the
  unchanged high-recall claim-harvester prompt and
  `accounts/fireworks/models/kimi-k2p6` as the extraction model — the model
  recorded for the frozen-100 extraction in `pipeline_runs.model`, which is
  the authoritative source. Exact-quote validation gate unchanged.

  **Correction log:** two false starts occurred before the correct run, both
  fully purged with zero downstream use. (1) The CLI-default `gpt-oss-120b`
  failed on its first request (rejects the extractor's no-reasoning option)
  and wrote nothing. (2) A Flash run based on an incorrect inference wrote
  1,816 fact units before `pipeline_runs.model` was consulted; all were
  deleted and the run is marked superseded in `pipeline_runs`. Rule adopted:
  **`pipeline_runs.model` is the protocol authority for every model choice;
  never infer a model from CLI defaults or surrounding usage.**
- **Embeddings:** `nomic-ai/nomic-embed-text-v1.5` for all four
  representations, unchanged.
- **Questions:** one per document (500 total; 100 per question type), authored
  by fresh transcript-only `claude -p --model fable` sessions using the frozen
  expansion-script prompt, mechanically validated against Postgres (document
  exists, cited turn range exists, supporting quote occurs verbatim in a cited
  turn) before insertion. Split `heldout-500-20260726`, IDs `heldout500-NNN`.
- **I6b runs (after all artifacts exist, single shot):**
  1. Budget-matched study: 4 representations × 256/512/1,024/2,048 estimated
     tokens, top-20 retrieval (replicating the original protocol as run).
  2. Frontier-XL: 4 representations × 2,048/4,096/8,192/16,384, top-300
     (post-`ef_search`-fix protocol).
  - Answer model DeepSeek V4 Flash, temperature zero, structured citations;
    greedy `cl100k_base` packing; blind Kimi K2.6 screens via
    `evaluate judge`, keys `judge_kimi_heldout500_v1` /
    `source_support_kimi_heldout500_v1`.
  - Coverage / reading-accuracy decomposition reproduced for the XL cells.

## Progress

- [x] 2026-07-26 — Sample frozen and committed; overlap checks passed.
- [ ] Representations built (`heldout500-v1`)
- [ ] Chunk embeddings complete
- [ ] Fact extraction complete (`heldout500-turnwise-opinion-v1`)
- [ ] Fact embeddings complete
- [ ] 500 questions authored and validated (split `heldout-500-20260726`)
- [ ] Budget-matched runs complete (16 runs)
- [ ] Frontier-XL runs complete (16 runs)
- [ ] Judging complete
- [ ] Results recorded; deltas vs development reported
