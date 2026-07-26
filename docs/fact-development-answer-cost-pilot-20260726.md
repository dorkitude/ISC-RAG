# Answer-accuracy and query-token pilot

**Status:** frozen-document development pilot. It is answer-level evidence over
the 100-question second-question replication set, but is not a paper result:
the documents are not held out and the blind judge is a fresh DeepSeek Flash
session rather than human-adjudicated or cross-model calibrated.

## Protocol

Questions are the 100 transcript-only Fable questions in
`fact-dev-sample-100-expansion-20260725`. The host embeds the original question
and retrieves the top 20 units using the same Fireworks embedding model. It
packs ranked evidence into a declared character cap, then asks DeepSeek V4
Flash to answer from sources only. Fireworks reports actual input/output tokens
for every answer; they are stored in `answers.usage`.

A separate, blinded Flash session sees only question, reference answer, and
generated answer, and returns a semantic correct/incorrect decision. It does
not see the representation, retrieved text, or run IDs. This removes the old
literal-match failure mode but remains a development screen, not an independent
human judgment.

The answer model and query embedding are identical across all conditions. Query
embedding cost therefore cancels in the comparison; retrieval runs locally in
the same pgvector database.

## Completed cost/accuracy points

| Condition | Context cap | Correct | Mean prompt tokens | Mean completion tokens | Mean total tokens |
| --- | ---: | ---: | ---: | ---: | ---: |
| Facts + quote | 2,048 chars | 85/100 | 597.5 | 207.1 | 804.6 |
| Fixed chunks | 2,048 chars | 45/100 | 514.9 | 200.1 | 715.0 |
| Facts + quote | 512 chars | 64/100 | 180.6 | 174.3 | 354.9 |

Run IDs are facts/2,048 `d998bf41-0cb9-4291-9806-52e8f2f8aa93`, fixed/2,048
`92a23195-e90c-4dc0-8dc7-b7c25d20f77a`, and facts/512
`6b0106e8-8ec8-473e-9989-b225bfc5ff0b`.

At the same cap, facts are substantially more accurate but carry modestly more
prompt overhead because every fact includes a label, provenance, and verified
quote. The operational point is the frontier: **512-character facts are 64%
correct, exceeding the 45% fixed-chunk result at a four-times-larger cap, while
using 354.9 versus 715.0 total answer-model tokens (50.4% fewer).** Both prompt
and completion tokens are lower for that fact condition, so it is cheaper under
any positive input/output pricing for the same answer model.

By family, facts at the 512-character cap are correct on 13/20 attribution,
15/20 decontextualization, 16/20 direct-control, 8/20 opinion, and 12/20
coreference questions. The low-budget fact condition improves on the fixed
2,048-character baseline in 44 paired questions; fixed alone is correct in 4.

## What remains before the claim can be promoted

1. Repeat this answer/cost frontier for turn-aware and semantic chunks, not
   only fixed chunks.
2. Add at least one different blind judge and calibrate it against human labels.
3. Replace character-cap packing with model-token-aware packing, while retaining
   provider-reported usage as the cost measure.
4. Repeat on document-disjoint, untouched data. Only then calculate the
   one-time ISC extraction/indexing break-even in dollars.
