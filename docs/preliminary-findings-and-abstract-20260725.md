# ISC-RAG preliminary findings and abstract

**Timestamp:** 2026-07-25 (America/Los_Angeles)

**Status:** development evidence on a frozen set of 100 MediaSum transcripts; not yet a paper result.

## Abstract

ISC-RAG tests whether an index of provenance-preserving, self-contained claims can improve retrieval-augmented question answering over human dialogue. Starting from transcripts, the system extracts high-recall factual and opinion claims, validates each against an exact source quote, attaches speaker and document metadata, and rewrites claims so resolvable pronouns and local references name their antecedents. We compare this fact-plus-quote representation with conventional transcript chunking using the same query embedding, pgvector retrieval, answer model, and frozen 100-document MediaSum development set.

Preliminary results suggest that rephrased facts occupy a better accuracy–context-cost frontier than fixed chunks. Across two 100-question transcript-derived sets, top-1 retrieval reaches the gold source turn for 165/200 questions with facts, versus 134/200 with fixed chunks. In a 100-question answer-generation pilot, 512-character fact contexts were judged correct on 64% of questions, exceeding fixed chunks at a four-times-larger 2,048-character context budget (45%), while consuming 50.4% fewer provider-reported answer-model tokens. The apparent advantage is strongest for attribution, coreference, decontextualization, and specific claims buried in long exchanges. These are development findings only: the remaining work is to compare turn-aware and semantic chunk controls, use token-aware packing and source-grounded scoring, calibrate independent/human judges, and replicate on document-disjoint held-out data.

## Preliminary findings

### Representation

The current fact representation contains a self-contained claim, its exact supporting quote, speaker attribution where available, and joinable document/source metadata. The first pass extracts explicit factual assertions **and** opinions, beliefs, judgments, preferences, and predictions. A second pass rewrites each retained claim using local transcript context, resolving pronouns only when the antecedent is knowable. Fact embeddings include both the rewritten fact and the verified quote; retrieved evidence returns both.

### Retrieval signal

Using raw question embedding and the same pgvector setup for both representations, facts retrieve a gold source turn more often than matched fixed chunks.

| Question set | Facts @1 | Fixed chunks @1 | Facts @3 | Fixed chunks @3 | Facts @5 | Fixed chunks @5 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Initial 100 | 83% | 70% | 92% | 82% | 95% | 87% |
| Question-level replication, 100 | 82% | 64% | 94% | 83% | 95% | 88% |
| Pooled | 82.5% | 67.0% | — | — | — | — |

The top-1 gap is 15.5 percentage points pooled (165/200 fact hits versus 134/200 fixed-chunk hits). This particular metric is deliberately conservative in its interpretation: gold-turn overlap can credit a retrieved fact from the correct turn even when it is a different claim. Future results should use fact-level evidence links or validated quote support.

### Why facts appear to help

The clearest qualitative cases are not merely shorter passages. The fact index makes transcript structure explicit:

- **Attribution:** a claim is returned with its normalized speaker, rather than leaving an answer model to distinguish a respondent from a nearby questioner in a long turn.
- **Coreference and decontextualization:** a claim such as “she ultimately removed that tweet” is retrieved as “Senator Maria Chappelle Nadal ultimately removed the tweet” when the antecedent is unambiguous.
- **Claim salience:** a compact, specific proposition—such as a number, belief, or prediction—can match its question without competing with unrelated broadcast vocabulary from the surrounding conversation.
- **Opinions:** retaining explicit judgments and beliefs expands the answerable surface beyond strictly objective assertions.

### Answer accuracy and query-token pilot

For a bounded answer-level pilot, all conditions use the same Fireworks query embedding and DeepSeek V4 Flash answer model. The host retrieves top-20 units and packs ranked evidence into a declared character budget. Provider-reported prompt and completion tokens are recorded for every answer. A fresh, blind Flash session judged semantic correctness from only the question, reference answer, and generated answer; it did not see the representation or retrieved evidence.

| Condition | Context cap | Correct | Mean prompt tokens | Mean completion tokens | Mean total tokens |
| --- | ---: | ---: | ---: | ---: | ---: |
| Facts + verified quote | 2,048 characters | 85/100 | 597.5 | 207.1 | 804.6 |
| Fixed chunks | 2,048 characters | 45/100 | 514.9 | 200.1 | 715.0 |
| Facts + verified quote | 512 characters | 64/100 | 180.6 | 174.3 | 354.9 |

At an equal 2,048-character cap, facts score 85% versus 45% for fixed chunks, though their provenance and quote fields add modest prompt-token overhead. The more useful operational comparison is the frontier: 512-character fact contexts score 64%, above the 45% fixed-chunk result at 2,048 characters, with 354.9 versus 715.0 mean total answer-model tokens—a 50.4% reduction. In paired questions, the 512-character fact condition alone is correct on 44 questions, while fixed chunks alone are correct on 4.

### What this does and does not establish

The present evidence supports a focused hypothesis: a one-time, provenance-preserving extraction and rephrasing step may let a RAG system answer more dialogue questions from substantially less context, particularly questions involving who said something, what a pronoun refers to, and a narrow claim embedded in a long conversation. It does **not** yet establish a publishable general result or a dollar-denominated return on the ISC cost.

Before making that claim, we will: (1) run the same answer/cost frontier against turn-aware and semantic chunk controls; (2) replace character packing with model-token-aware packing; (3) score answer correctness and citation support with independently calibrated and human-reviewed judgments; and (4) repeat on untouched, document-disjoint data. Only then should we estimate the extraction/indexing break-even against query-time savings.

## Reproducibility pointers

- Answer/token pilot: [`fact-development-answer-cost-pilot-20260726.md`](fact-development-answer-cost-pilot-20260726.md)
- Qualitative fact-win audit and metric caveats: [`fact-development-sample-100-fact-wins.md`](fact-development-sample-100-fact-wins.md)
- Frozen development sample and boundary: [`README.md`](../README.md)
