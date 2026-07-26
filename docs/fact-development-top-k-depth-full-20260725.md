# Full frozen-100 top-k answer-depth curve

**Status:** controlled development result on the frozen 100-document,
100-question MediaSum set. It is a full answer-level depth curve, but not a
paper result: the documents are development data and DeepSeek Flash is both the
answer model and the automated evaluator at the user's requested setting.

## Protocol

For every question, each representation retrieves exactly the top
`k ∈ {1, 3, 5, 10, 20, 30}` units from the same pgvector index. There is no
context packing cap: all retrieved units at that depth are provided to the
answer model. This makes the manipulation a literal top-k depth curve, while
also exposing its actual token cost.

All 24 conditions use the same question text, embedding model, retrieval
configuration, and `accounts/fireworks/models/deepseek-v4-flash` answer model.
The answer must return structured source IDs. Fresh DeepSeek Flash sessions
then independently screen (a) semantic correctness against the reference answer
and (b) whether the answer-selected source directly supports the answer. A
cell below is **answerable-and-supported** only when both screens are positive.

The full run contains 2,400 answers and 4,800 Flash judgments, all completed
without judgment failures. A transient provider timeout was recovered using
bounded client retries and `evaluate run --resume-run`; partial or duplicate
runs are excluded from this ledger.

## Curve: correct and source-supported answers

| Representation | k=1 | k=3 | k=5 | k=10 | k=20 | k=30 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Facts + quote | 47/100 | 68/100 | 72/100 | 76/100 | 79/100 | 79/100 |
| Fixed chunks | 49/100 | 60/100 | 62/100 | 62/100 | 66/100 | 61/100 |
| Turn-aware chunks | 48/100 | 53/100 | 54/100 | 63/100 | 63/100 | 69/100 |
| Semantic chunks | 49/100 | 53/100 | 57/100 | 58/100 | 59/100 | 65/100 |

At top-1, chunk coverage is marginally higher, but it is not a cost-neutral
comparison because one chunk is much larger than one fact. At top-3 and every
higher depth, facts have the highest answerable-and-supported coverage. Facts
saturate at 79/100 by k=20; adding ten more facts does not help. Fixed chunks
fall from 66/100 at k=20 to 61/100 at k=30, an example of additional context
distracting the answer model rather than monotonically improving it.

## Actual answer-model cost

| Representation | Mean total tokens @1 | @3 | @5 | @10 | @20 | @30 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Facts + quote | 233 | 366 | 496 | 817 | 1,414 | 1,892 |
| Fixed chunks | 889 | 2,311 | 3,734 | 7,342 | 14,465 | 21,599 |
| Turn-aware chunks | 903 | 2,348 | 3,757 | 7,369 | 14,306 | 21,328 |
| Semantic chunks | 729 | 1,892 | 3,084 | 6,036 | 11,697 | 17,212 |

Facts are cheaper at every depth. The salient points are:

- At **k=1**, facts are 47% answerable-and-supported at 233 mean tokens;
  chunks are 48–49% at 729–903 tokens.
- At **k=3**, facts are **68%** at 366 tokens, versus 53–60% at
  1,892–2,348 tokens for chunks.
- At **k=20**, facts are **79%** at 1,414 tokens, versus 59–66% at
  11,697–14,465 tokens for chunks.

Thus a fixed k is not a meaningful efficiency comparison by itself. The fact
representation delivers higher coverage from a much smaller unit: at k=20 it
is roughly 8–10× cheaper in actual answer-model tokens while answering 13–20
more questions than the chunk controls.

## Minimum successful k

For each question/representation pair, the first successful k is the smallest
tested depth at which both Flash screens were positive. Among questions that
ever succeed by k=30:

| Representation | Answerable by k=30 | Mean minimum k | Median minimum k | Never by k=30 |
| --- | ---: | ---: | ---: | ---: |
| Facts + quote | 84/100 | 3.21 | 1 | 16/100 |
| Fixed chunks | 76/100 | 3.25 | 1 | 24/100 |
| Turn-aware chunks | 75/100 | 4.49 | 1 | 25/100 |
| Semantic chunks | 73/100 | 4.14 | 1 | 27/100 |

The median is 1 for every representation because many questions are simple
direct lookups. The more informative measures are coverage and cost: facts
make eight more questions answerable than fixed chunks and 9–11 more than the
other chunk controls by k=30, while avoiding the massive context growth needed
to reach those depths with chunks.

## Interpretation and boundaries

This is evidence for the proposed mechanism, with a useful qualification. Fact
cards do not universally need fewer *units* than chunks: some direct questions
are answered from one chunk. They do provide a better joint
**coverage–depth–token** curve on this frozen development set. The benefit
appears after top-1, where facts retrieve compact, answer-ready propositions
while chunk contexts rapidly become expensive and distractible.

Because DeepSeek Flash judged DeepSeek Flash outputs here, the absolute
correctness values should be treated as a self-model development screen. The
next confirmatory run should preserve this exact protocol but use human
calibration and/or a different judge on a document-disjoint test set. Do not
use this result to calculate a general extraction-cost break-even yet.

## CLI support

- `isc-rag evaluate run --resume-run UUID` resumes a failed, configuration-
  matching run without regenerating durable answers.
- `isc-rag evaluate depth-summary --runs UUID,... --correctness-key KEY
  --support-key KEY` reports the stored curve and first-success depths.

The selected complete run IDs are recorded in PostgreSQL. Facts are
`c09be366-1289-4585-83b1-1436bba02f5e`,
`2da5b00e-8f98-428e-be22-1a5bb7173205`,
`68e765eb-4053-433a-8d63-b8f238f2b5e0`,
`474a41c9-2a76-49c1-89e6-126b22b32db1`,
`356d25e4-d974-48ab-8a52-50bde9468f16`, and
`4dd13da9-5773-451f-bd3f-28777aba9cb5` in k order; the complete control-run
ledger can be regenerated from the CLI and experiment tables.
