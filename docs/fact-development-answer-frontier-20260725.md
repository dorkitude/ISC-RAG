# Frozen-100 answer accuracy–cost frontier

**Status:** controlled development result on the frozen 100-document MediaSum
sample. This is stronger than the earlier character-budget pilot but is still
not a paper result: questions are from the same frozen documents, and the
independent automated judge still needs human calibration.

## Question

At a fixed answer-model context budget, do provenance-preserving, rephrased
facts answer dialogue questions more accurately than fixed, turn-aware, or
semantic transcript chunks? At a fixed accuracy level, do facts require fewer
answer-model tokens?

## Protocol

- **Documents and questions:** the frozen 100 MediaSum transcripts and the
  independently Fable-authored, source-validated 100-question replication set
  `fact-dev-sample-100-expansion-20260725` (20 questions in each of
  attribution, decontextualization, direct control, opinion/belief/judgment,
  and pronoun/coreference).
- **Conditions:** facts plus verified quote (`facts/turnwise-opinion-v1`),
  4,000-character fixed chunks, turn-aware chunks, and semantic chunks. All
  chunk controls are the versions built from the same 100 documents.
- **Retrieval:** original question text, the same
  `nomic-ai/nomic-embed-text-v1.5` Fireworks embedding, local pgvector search,
  and top-20 retrieval for every condition.
- **Packing:** ranked retrieved sources are greedily packed to one of 256,
  512, 1,024, or 2,048 **estimated context tokens**, using the documented
  `cl100k_base` tokenizer. This is model-token-aware packing, but it is an
  estimate rather than DeepSeek's exact tokenizer. Fireworks-reported prompt,
  completion, and total tokens are the cost measurements.
- **Answering:** `accounts/fireworks/models/deepseek-v4-flash`, temperature
  zero, constrained JSON output with an answer and IDs of the retrieved sources
  it claims support the answer. Source labels are representation-neutral.
- **Scoring:** a fresh `accounts/fireworks/models/kimi-k2p6` session judged
  semantic answer correctness from the question, reference, and generated
  answer only. A separate fresh Kimi session judged whether the model-selected
  cited source(s) directly support the generated answer. Neither judge sees
  representation, budget, retrieval rank, or run ID.

The answer and judge models are different. The Kimi judgments are still an
automated development screen, not human adjudication.

## Results

`Correct` is Kimi's blind semantic-correctness verdict; `supported` is its
separate source-support verdict. `Total tokens` is the actual mean
provider-reported answer-model prompt plus completion usage, not the packing
estimate.

| Representation | Budget | Correct | Supported | Mean prompt | Mean completion | Mean total |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Facts + quote | 256 | 74/100 | 77/100 | 363.4 | 26.4 | 389.8 |
| Fixed chunks | 256 | 23/100 | 24/100 | 401.6 | 17.6 | 419.3 |
| Turn-aware chunks | 256 | 21/100 | 23/100 | 401.9 | 18.4 | 420.3 |
| Semantic chunks | 256 | 28/100 | 29/100 | 387.0 | 17.8 | 404.9 |
| Facts + quote | 512 | 81/100 | 83/100 | 628.6 | 25.8 | 654.4 |
| Fixed chunks | 512 | 46/100 | 44/100 | 643.4 | 20.2 | 663.6 |
| Turn-aware chunks | 512 | 45/100 | 46/100 | 643.5 | 20.4 | 663.9 |
| Semantic chunks | 512 | 45/100 | 44/100 | 610.8 | 20.4 | 631.2 |
| Facts + quote | 1,024 | 87/100 | 86/100 | 1,131.9 | 26.9 | 1,158.8 |
| Fixed chunks | 1,024 | 60/100 | 59/100 | 883.6 | 22.5 | 906.1 |
| Turn-aware chunks | 1,024 | 56/100 | 57/100 | 898.1 | 23.0 | 921.1 |
| Semantic chunks | 1,024 | 61/100 | 61/100 | 846.4 | 23.8 | 870.1 |
| Facts + quote | 2,048 | 86/100 | 87/100 | 1,381.9 | 25.6 | 1,407.5 |
| Fixed chunks | 2,048 | 64/100 | 64/100 | 1,779.1 | 22.2 | 1,801.3 |
| Turn-aware chunks | 2,048 | 59/100 | 62/100 | 1,783.6 | 22.9 | 1,806.5 |
| Semantic chunks | 2,048 | 63/100 | 63/100 | 1,878.6 | 22.5 | 1,901.1 |

### Interpretation

Facts lead every chunk strategy at every tested budget. At the same nominal
512-token budget, facts are correct on **81%** of questions, compared with
45–46% for each chunk baseline, while actual total answer-model usage is
similar (654.4 tokens for facts versus 631.2–663.9 for chunks). At 256 tokens,
facts reach **74%**, exceeding the best 512-token chunk result (46%) while
using 38.2–41.3% fewer actual answer-model tokens.

The accuracy–cost frontier is more striking at adjacent budgets: 512-token
facts (81%, 654.4 tokens) exceed every 1,024-token chunk control (56–61%,
870.1–921.1 tokens); 1,024-token facts (87%, 1,158.8 tokens) exceed every
2,048-token chunk control (59–64%, 1,801.3–1,901.1 tokens). Thus, on this
development set, the fact representation is both more accurate and cheaper at
these matched frontier points. The slight 87% → 86% fact change from 1,024 to
2,048 tokens is consistent with extra retrieved context sometimes adding a
distractor; it should not be overinterpreted without repeated generations.

At 512 tokens, facts are especially strong in direct control (18/20),
attribution (17/20), decontextualization (17/20), and coreference (16/20).
Opinion/belief/judgment is still a fact lead (13/20 versus 7–10/20 controls),
but is the least strong of the five families.

## Run ledger

| Representation | 256 | 512 | 1,024 | 2,048 |
| --- | --- | --- | --- | --- |
| Facts | `76e177e5-b41c-4817-9ec2-9f7efc295922` | `28d47a4d-69f5-45ea-82ed-fe41429cc861` | `dfa0dd13-d115-4b0e-b453-8ffdd354a66b` | `12225f72-ed38-422b-a163-dd5e6599308f` |
| Fixed | `8fc06526-982a-4c7b-9518-1ea213f5fa05` | `5cee9df7-76e8-4ebd-943f-22356adf1e3f` | `bc241efa-dd77-41cd-8d47-f250b0507c10` | `af62b994-b63b-428c-a386-e884f2c93a52` |
| Turn | `a0e9f6c9-7031-4995-a709-5c3db18ac137` | `24e5e888-259b-4f04-ab4c-b8a0023cc0f5` | `14bc49ae-9929-4814-88dd-382813bf2b69` | `2796a54b-241a-4dc3-94c6-d7f77d0b670b` |
| Semantic | `d189aa2b-04be-4270-8b61-5fe9eb8bc86b` | `f38ca0d9-6d21-4481-b8ed-56c3311af5a9` | `f6c48244-fccf-49f4-9c1e-c3794d6cf648` | `4ef714a6-728a-4763-a99b-22b31ab1961f` |

An earlier fixed-256 run aborted after seven answers because the answer model
returned an empty `answer` field. It remains marked `failed` in PostgreSQL and
is excluded. The replacement fixed-256 run above completed all 100 items under
the clarified “Insufficient source evidence.” protocol.

## Human calibration package

[`human-calibration-worksheet.md`](../experiments/fact-development-answer-frontier-20260725/human-calibration-worksheet.md)
contains 80 blinded review items, one deterministic item per representation ×
budget × question-family stratum. It preferentially samples questions whose
Kimi correctness outcome differs across representations at the same budget.
The companion `human-calibration-worksheet-key.jsonl` maps item numbers to
conditions and Kimi verdicts and should be withheld from reviewers. No human
labels have been collected yet.

## Next boundary

1. Obtain independent human labels for the 80 blinded items and compare them
   with Kimi's correctness and support judgments.
2. Freeze the answer prompt, packing implementation, answer model, judge
   protocol, and primary budget before touching new documents.
3. Only after the frozen-100 development methodology is accepted, run the
   untouched document-disjoint test and compute extraction/indexing break-even
   from current provider prices and measured token usage.
