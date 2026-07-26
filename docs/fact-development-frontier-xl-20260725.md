# Frozen-100 extended answer frontier (frontier-XL): large budgets, top-k 300

**Status:** controlled development result on the frozen 100-document MediaSum
sample and the `fact-dev-sample-100-expansion-20260725` question set. It
removes the candidate-supply confound of the earlier 256–2,048-token frontier
but remains development evidence: same-document data, automated Kimi screens,
no human calibration yet.

## Question

The original frontier retrieved only the top 20 units before packing. Twenty
facts average roughly 1,300 estimated tokens, so at the 2,048-token budget the
fact condition ran out of candidates rather than budget. This experiment asks:
when the token budget — not candidate supply — is the binding constraint at
every cell, do facts remain more accurate and cheaper than chunks, and what
happens to every representation as budgets grow to 16,384 tokens?

## Protocol

Identical to the [original frontier](fact-development-answer-frontier-20260725.md)
except:

- **Retrieval depth:** top-300 for every condition (300 facts ≈ 19k estimated
  tokens; 300 fixed chunks ≈ 210k), so greedy `cl100k_base` packing truncation
  binds everywhere.
- **Budgets:** 2,048 / 4,096 / 8,192 / 16,384 estimated context tokens.
  DeepSeek V4 Flash accepts 1M input tokens, so the model window is never the
  constraint.
- **Judging:** the new reproducible `isc-rag evaluate judge` command (blind
  Kimi K2.6 correctness and source-support screens with prompts frozen in
  `internal/evaluation/judge.go`), keys `judge_kimi_frontier_xl_v1` and
  `source_support_kimi_frontier_xl_v1`. 1,600/1,600 answers judged without
  failures.

### Retrieval-depth defect found and fixed mid-experiment

An Opus spot-check of the first completed runs found the fact condition packing
at most ~40 units (~2,770 tokens) regardless of budget, with 96/100 answers
byte-identical between the 2,048 and 16,384 fact cells. The cause was
pgvector's HNSW `ef_search` default of 40, which silently caps result count
below the SQL `LIMIT`. `internal/retrieval/retrieve.go` now sets
`SET LOCAL hnsw.ef_search = clamp(top_k, 40, 1000)` inside the query
transaction. The four fact runs were regenerated with the fix and verified
budget-bound (mean estimated context 2,015 / 4,063 / 8,073 / 12,606 against
the four caps; the 16,384 mean is below cap because some questions exhaust the
per-question fact pool). All chunk runs and all previously recorded
experiments used depths at or below 40 and are unaffected.

## Results: correct and supported answers

`Correct` is the blind Kimi semantic-correctness verdict; `Corr+Supp`
additionally requires the separate blind source-support verdict. `Mean total
tokens` is provider-reported prompt plus completion usage.

| Representation | Budget | Correct | Corr+Supp | Mean total tokens |
| --- | ---: | ---: | ---: | ---: |
| Facts + quote | 2,048 | 92/100 | 86/100 | 2,233.3 |
| Facts + quote | 4,096 | 87/100 | 83/100 | 4,332.6 |
| Facts + quote | 8,192 | 88/100 | 81/100 | 8,445.9 |
| Facts + quote | 16,384 | 86/100 | 81/100 | 13,094.1 |
| Fixed chunks | 2,048 | 62/100 | 59/100 | 1,801.2 |
| Fixed chunks | 4,096 | 70/100 | 67/100 | 3,946.8 |
| Fixed chunks | 8,192 | 68/100 | 66/100 | 8,002.2 |
| Fixed chunks | 16,384 | 75/100 | 73/100 | 16,243.4 |
| Turn-aware chunks | 2,048 | 59/100 | 55/100 | 1,804.9 |
| Turn-aware chunks | 4,096 | 59/100 | 55/100 | 3,872.7 |
| Turn-aware chunks | 8,192 | 64/100 | 61/100 | 8,003.9 |
| Turn-aware chunks | 16,384 | 75/100 | 71/100 | 16,271.2 |
| Semantic chunks | 2,048 | 63/100 | 63/100 | 1,901.2 |
| Semantic chunks | 4,096 | 61/100 | 61/100 | 3,964.0 |
| Semantic chunks | 8,192 | 65/100 | 65/100 | 8,049.9 |
| Semantic chunks | 16,384 | 70/100 | 70/100 | 16,361.4 |

### Paired budget response (2,048 → 16,384, correct-and-supported)

| Representation | Regressed (ok at 2k, not at 16k) | Gained (not ok at 2k, ok at 16k) |
| --- | ---: | ---: |
| Facts + quote | 7 | 2 |
| Fixed chunks | 1 | 15 |
| Turn-aware chunks | 4 | 20 |
| Semantic chunks | 4 | 11 |

### Decomposition: coverage versus reading accuracy

Each cell splits into whether any packed source overlapped the gold turn range
(**gold packed**) and blind correctness among only those questions
(**read-acc**), separating retrieval coverage from reading quality:

| Representation | Gold packed @2k/4k/8k/16k | Read-acc given gold @2k/4k/8k/16k |
| --- | --- | --- |
| Facts + quote | 100 / 100 / 100 / 100 | 92.0 / 87.0 / 88.0 / 86.0 |
| Fixed chunks | 80 / 88 / 90 / 95 | 76.3 / 78.4 / 75.6 / 78.9 |
| Turn-aware chunks | 75 / 82 / 91 / 95 | 77.3 / 72.0 / 70.3 / 78.9 |
| Semantic chunks | 78 / 85 / 88 / 95 | 78.2 / 70.6 / 72.7 / 73.7 |

This decomposition addresses the needle-in-a-haystack objection — that large
chunk budgets merely measure the answer model's haystack weakness. Fixed-chunk
reading accuracy is flat (~76–79%) from 2,048 to 16,384 tokens: the model finds
gold in a large packed context about as reliably as in a small one, and chunk
budget gains are therefore mostly *coverage* gains (80%→95% gold packed).
The dominant effect is instead a roughly 14–16-point reading-accuracy gap
between compiled facts and raw transcript that is present already at the
2,048-token budget, where haystack pressure is minimal. Genuine long-context
degradation appears as a secondary effect: the fact conditions' own 92→86
decline occurs at fixed 100% coverage and is therefore pure reading
degradation with volume, and the turn/semantic mid-budget dips recover only at
16,384.

## Findings

1. **Facts dominate on both axes.** The 2,048-token fact cell (92/100 correct,
   86/100 correct-and-supported, 2,233 mean tokens) exceeds every chunk cell at
   every budget. The best any chunk condition reaches is 75/100 correct — at
   roughly 16,300 tokens, 7.3× the fact cost. Facts' *worst* cell (86/100)
   still beats chunks' best.
2. **The cheapest fact budget is the best fact budget.** With supply no longer
   the constraint, fact accuracy declines as budget grows (92 → 87 → 88 → 86),
   with 7 paired regressions against 2 gains from 2,048 to 16,384. Compiled
   context is not immune to long-context degradation; its optimum simply sits
   near 2k tokens rather than beyond the tested range. The prescription is not
   "facts tolerate large contexts better" but "facts make large contexts
   unnecessary."
3. **Chunks buy accuracy with budget, expensively and non-monotonically.**
   Turn-aware chunks gain 20 questions from 2,048 to 16,384 and fixed gains
   15 — the evidence is present in raw text at depth — but fixed dips at
   8,192 (68 < 70 at 4,096), and a qualitative spot-check verified individual
   correct→insufficient flips with the gold evidence still packed:
   within-experiment lost-in-the-middle cases.
4. **Behavioral health checks passed.** Packing never exceeded any budget; no
   hallucination on genuinely unanswerable questions (all conditions abstain);
   one systematic budget- and representation-invariant coreference error
   ("whose father led Indonesia" → all conditions answer the father, not the
   asked-for child) reflects question difficulty, not representation.

## Boundaries

Same-document development data; Kimi is an automated development judge pending
the human-calibration worksheet; single answer model (DeepSeek V4 Flash) —
a model-generality follow-up (Kimi K2.6 and GLM 5.2 as answer models with
cross-judging) is running as of this record's timestamp and will be recorded
separately.

## Reproducibility

Final run IDs (representation @ budget):

- facts (post-fix): a00fb789-a14d-4d9a-9285-61d4af91f005 @2048,
  519b1c35-d314-4c3b-82a4-f60e7cbaa088 @4096,
  5cf07cfc-a135-4080-b7e3-92c19c58d992 @8192,
  adbfce87-cac8-4aed-80e8-d3262a094b41 @16384
- fixed: 737b9df9-158c-45c8-bf16-6ca2dbc3bdaf @2048,
  85a4d5ef-7d2b-4a8a-9f05-0bd4b283c865 @4096,
  57834359-fbc2-4b7e-a776-f74df6b53df6 @8192,
  b13cc9e8-1615-40ce-a13d-00021ee64647 @16384
- turn: f92ee367-79d6-4a88-b891-4f2305eb9d22 @2048,
  ada7f9b2-cc08-4524-a280-fcdfd2a37e0f @4096,
  99b978b9-966b-44a9-96b2-ca3fe4857eac @8192,
  a72590e0-ecc7-40ab-b47b-de381f8b948f @16384
- semantic: e42a506f-01a7-4253-8db8-f287fc3f504c @2048,
  0192d427-8800-45ed-8e22-94caecef6e03 @4096,
  4e1fe17f-deb1-4132-87d5-bdd056aa0d5d @8192,
  f2cbaf1f-afcd-43c9-8aa1-2e7ece16dcbd @16384

The four superseded, supply-starved fact runs remain in the database without
`judge_kimi_frontier_xl_v1` keys and are excluded from this ledger. Judging is
re-runnable with `isc-rag evaluate judge`; aggregate tables reproduce with the
SQL in this document's history and `evaluate depth-summary`.
