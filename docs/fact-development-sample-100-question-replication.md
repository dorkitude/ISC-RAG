# Additional-question retrieval replication

**Status:** development evidence only. It tests whether the earlier fact
retrieval lead survives a new set of questions over the same frozen 100
documents. It does not establish document-disjoint generalization or answer
accuracy.

## Protocol

The new set contains 100 source-grounded questions at
[`questions.jsonl`](../experiments/fact-development-sample-100-expansion-20260725/questions.jsonl),
one per frozen document and 20 per question family. Claude Fable 5 generated
each item in a fresh transcript-only session. It was told the prior question
only to avoid the same claim; it saw neither representations nor retrieval
results. Exact quote and cited-turn validation occurred before admission.

The comparison embeds the original question text and searches the same pgvector
index, with no query-rewriting or answer agent. Both conditions use the same
embedding model (`nomic-ai/nomic-embed-text-v1.5`) and retrieve from the same
100 documents:

- facts: `facts/turnwise-opinion-v1`, fact plus validated quote embedding;
- control: `fixed/fact-dev-sample-100-fixed-v1`, 4,000-character chunks with
  400-character overlap.

The score is gold-turn overlap recall. It is a retrieval screen, not
fact-level answer support or answer correctness.

## Results

| Depth | Facts | Fixed chunks | Difference |
| ---: | ---: | ---: | ---: |
| @1 | 82/100 | 64/100 | +18 points |
| @3 | 94/100 | 83/100 | +11 points |
| @5 | 95/100 | 88/100 | +7 points |

Run IDs, facts then fixed: @1
`7e4aea46-6a7d-453a-9080-a49386073332` /
`e5581de8-e412-491d-b25e-66effd17ff0c`; @3
`6a481ef7-91f1-4a75-9129-ddf6bf65434a` /
`b97cc8a1-0b29-4b21-a68d-0a405789f2bd`; @5
`a618d1bf-7906-477b-b509-4ddde1a78631` /
`417a4347-91d8-4647-a8bf-d45099a99cf6`.

At @1, the question-family breakdown is:

| Family | Facts | Fixed | Fact-only | Fixed-only |
| --- | ---: | ---: | ---: | ---: |
| Attribution | 18/20 | 14/20 | 5 | 1 |
| Decontextualization | 18/20 | 17/20 | 2 | 1 |
| Direct control | 17/20 | 13/20 | 6 | 2 |
| Opinion/belief/judgment | 13/20 | 13/20 | 5 | 5 |
| Pronoun/coreference | 16/20 | 7/20 | 11 | 2 |

Pooled with the prior, separately generated 100-question set, top-1 retrieval
is 165/200 for facts versus 134/200 for fixed chunks. Facts are the only hit
for 55 questions; fixed chunks are the only hit for 24. This replicates the
retrieval advantage, particularly for coreference and direct/attribution
questions, while opinions remain a tie in this second set.

## Interpretation and next boundary

This supports the claim that the current fact representation improves
retrieval of the gold turn under original-question semantic search. It does
**not** prove more accurate answers: a gold turn can contain an answer-irrelevant
fact. The next controlled study should use a document-disjoint test set,
fact-level answer-support labels, token-budget packing, structured citations,
and blind answer adjudication before making an accuracy or cost claim.
