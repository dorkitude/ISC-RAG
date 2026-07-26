# Decontextualization-question retrieval probe

**Status:** exploratory development diagnostic, run 2026-07-25. This is not a
held-out paper result or an answer-quality comparison.

## Question set

This probe turns the proposed failure modes into a concrete, source-grounded
development check. The tracked
[`questions.jsonl`](../experiments/decontextualized-question-development-20260725/questions.jsonl)
contains 15 questions over three CNN MediaSum test documents that already had
both comparison representations locally indexed (`CNN-10004`, `CNN-100082`,
and `CNN-100095`):

| Stratum | Count | Intended pressure |
| --- | ---: | --- |
| `pronoun_resolution` | 5 | Resolve a pronoun or abbreviated referent to a named entity. |
| `attribution` | 5 | Identify the person who made a claim rather than merely the claim. |
| `decontextualization` | 5 | Retrieve a claim after its dialogue-local wording has been turned into a standalone question. |

Each question has one source document and a canonical gold turn range. The
questions were written after inspecting these transcripts, so they are
explicitly **development data** and must never enter a held-out evaluation.
Their SHA-256 is
`64885374c1f26ba5e9559f5ce83e610161df75765fe4198df80fe81cec55c6c4`.

## Method

This is retrieval only: no Q&A model, reranking, or answer judge participates.
For every question, the runner embeds its text with the same Fireworks Nomic
v1.5 embedding model, queries pgvector at `top_k=5`, and marks a hit only when
a retrieved source unit covers the preannotated document and turn range.

The two local, document-matched conditions are:

| Condition | Source units |
| --- | --- |
| Fixed chunks | `fixed/pilot-100x100-v1` |
| Named rephrased facts | `facts/pilot-named-v3` |

The fact condition exposes its fact text for retrieval and retains an exact
quote/speaker field for return. Fixed chunks are the existing chunk projection;
they currently do not serialize speaker labels into their searchable text.
That is a real limitation of this development comparison, rather than evidence
that speaker-aware chunks cannot close the gap.

## Retrieval results

| Condition | Pronoun | Attribution | Decontextualization | All questions |
| --- | ---: | ---: | ---: | ---: |
| Fixed chunks | 4 / 5 (80.0%) | 4 / 5 (80.0%) | 2 / 5 (40.0%) | 10 / 15 (66.7%) |
| Named rephrased facts | 4 / 5 (80.0%) | 4 / 5 (80.0%) | 5 / 5 (100.0%) | 13 / 15 (86.7%) |

Facts uniquely retrieved four gold spans, three of them decontextualization
items. Fixed chunks uniquely retrieved one pronoun-resolution item. Both
missed the same attribution item. Thus this tiny probe gives a plausible
mechanistic signal for the rephrasing hypothesis, but it does *not* establish a
general fact advantage and it does not yet show that an answering model uses
the evidence correctly.

The most legible fact win is the Mitchell Wade question. The fact index
retrieves the standalone text, “The contractor is Mitchell Wade and his company
is MZM Incorporated,” with the exact source quote. The corresponding fixed
unit is a much longer multi-turn chunk, and did not cover the gold turn in its
top five in this run.

## Reproduce

```sh
shasum -a 256 experiments/decontextualized-question-development-20260725/questions.jsonl
go run ./experiments/decontextualized-question-development-20260725 --top-k 5
jq '.' artifacts/decontextualized-question-development-20260725/retrieval.json
```

The JSON report is ignored because it is a rerunnable local embedding/retrieval
artifact. The runner is intentionally limited to the two existing development
projections. Before drawing a paper conclusion, create a document-disjoint,
pre-registered set; add speaker-labeled chunk variants; measure rank and answer
quality; and use the identical frozen corpus slice for every representation.
