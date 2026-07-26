# ISC-RAG preliminary findings and abstract

**Timestamp:** 2026-07-25 (America/Los_Angeles)

**Status:** development evidence on a frozen set of 100 MediaSum transcripts; not yet a paper result.

## Abstract

ISC-RAG tests whether an index of provenance-preserving, self-contained claims can improve retrieval-augmented question answering over human dialogue. Starting from transcripts, the system extracts high-recall factual and opinion claims, validates each against an exact source quote, attaches speaker and document metadata, and rewrites claims so resolvable pronouns and local references name their antecedents. We compare this fact-plus-quote representation with conventional transcript chunking using the same query embedding, pgvector retrieval, answer model, and frozen 100-document MediaSum development set. Unlike Dense X Retrieval's propositions, which serve mainly as finer-grained index keys into raw text, each ISC-RAG fact is intended as compiled, provenance-validated read-time context in its own right.

Across two 100-question transcript-derived sets, top-1 retrieval reaches the gold source turn for 165/200 questions with facts, versus 134/200 with fixed chunks. In a controlled 100-question answer-generation frontier using four context budgets and all three chunk controls, facts are more accurate at every tested budget. At the matched 512-token budget, facts are 81% correct versus 45–46% for fixed, turn-aware, and semantic chunks at comparable answer-model cost; 256-token facts are 74% correct, exceeding every 512-token chunk condition while using 38–41% fewer provider-reported answer-model tokens. The advantage is strongest for attribution, coreference, decontextualization, and specific claims buried in long exchanges. These are development findings only: the data are still same-document development data and Kimi judgments await human calibration before replication on document-disjoint held-out transcripts.

## Preliminary findings

### Representation

The current fact representation contains a self-contained claim, its exact supporting quote, speaker attribution where available, and joinable document/source metadata. The first pass extracts explicit factual assertions **and** opinions, beliefs, judgments, preferences, and predictions. A second pass rewrites each retained claim using local transcript context, resolving pronouns only when the antecedent is knowable. Fact embeddings include both the rewritten fact and the verified quote; retrieved evidence returns both.

## How to read the results

There are **two different evaluations** below. They should not be conflated.

1. **Retrieval `@k`** asks whether search found the evidence *before* an answer
   model tries to answer. `@1` means the single highest-ranked pgvector result
   overlaps the gold source turn; `@3` means at least one of the first three
   results does; and `@5` means at least one of the first five does. It is not
   “the first tool call answered correctly.” Higher `k` normally improves the
   chance of finding evidence but gives the answer model more material to read.
2. **The answer frontier** asks whether an answer model can answer correctly
   from a bounded set of retrieved sources. Each row is one complete RAG
   condition over the same 100 questions. `Budget` is the maximum retrieved
   context passed to the answer model, estimated with `cl100k_base`; it does
   not include the answer instruction or answer completion. `Mean total tokens`
   is the provider-reported prompt plus completion usage, which is the relevant
   query-cost measure.

In the answer table, `Correct` is a blind Kimi semantic comparison of the
question, reference answer, and generated answer. `Cited support` is a separate
blind Kimi judgment that the source IDs chosen by the answer model directly
support its answer. These are related but intentionally distinct: an answer can
match the reference yet be unsupported by its cited source, or be supported but
incomplete. Kimi is a development judge, not human ground truth.

Read comparisons in two ways: compare rows at the **same budget** for a
representation-quality result; compare the smallest fact budget that reaches a
given quality with chunk rows at larger budgets for an **accuracy–cost frontier**
result. Do not compare raw `@k` retrieval percentages directly with answer
percentages.

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

### Relation to Dense X Retrieval and contextual chunking

The nearest published work is Dense X Retrieval ([Chen et al., EMNLP
2024](https://aclanthology.org/2024.emnlp-main.845/); [arXiv:2312.06648](https://arxiv.org/abs/2312.06648)),
which shows that indexing Wikipedia by atomic "propositions" outperforms
passage- and sentence-level indexing, with downstream QA gains under a fixed
reading-token budget. ISC-RAG shares the granularity thesis but differs in what
the atomic unit *is for*.

In Dense X, the proposition is primarily a better **index key**: retrieval can
match a question against a compact factoid, but the system may then hand the
reader the surrounding passage, and the propositions themselves are generated
by a distilled model without source verification. The reader still reconstructs
meaning from raw text at query time — in our terms, query-time semantic
reconstruction (QSR) with a sharper pointer into it.

ISC-RAG's fact is intended as the **compiled payload itself**, not a pointer.
The semantic work — resolving pronouns, attributing the claim to a normalized
speaker, isolating the assertion from surrounding dialogue — is performed once
at ingest, validated against an exact source quote, and stored with speaker and
turn-span provenance. At query time the answer model reads the compiled claim
and its verified quote directly; it does not re-derive them from a transcript
window. This is why the frontier below is the central measurement: if the
compiled unit is doing its job, accuracy should hold or improve as read-time
context shrinks, which is what the 256- and 512-token fact rows test. Dense X
compiles the index; ISC-RAG aims to compile the substrate.

Anthropic's [Contextual Retrieval](https://www.anthropic.com/engineering/contextual-retrieval)
(2024) is the same distinction seen from the chunking side. It pays an
ingest-time LLM cost to prepend 50–100 tokens of situating context to each
chunk before embedding and BM25 indexing, substantially reducing retrieval
failures. That is ingest-time semantic work, and we regard it as a partial form
of ISC — but the compiled artifact is again the pointer: the retrieval unit
remains a raw chunk that has grown *larger*, the reader still reconstructs the
answer from raw transcript text, and the generated context is not validated
against the source. Both Dense X and Contextual Retrieval improve where
retrieval lands; neither changes what the answer model must read.

More recent proposition work extends Dense X toward multi-hop graph reasoning
rather than toward compiled read-time context:
[SiReRAG](https://arxiv.org/abs/2412.06206) (ICLR 2025) groups propositions and
entities into relatedness trees; [PropRAG](https://arxiv.org/abs/2504.18070)
(EMNLP 2025) runs beam search over proposition paths; and
[AtomicRAG](https://arxiv.org/abs/2604.20844) (2026) replaces chunks with
"knowledge atoms" in an atom–entity graph, arguing explicitly that chunks
rigidly group multiple atomic facts together. In each, atomic units organize
retrieval and reasoning structure; none validates units against exact source
spans, and none measures answer accuracy under shrinking read-time context
budgets, which is the measurement this document centers.

The dialogue setting also raises compilation problems Wikipedia prose does not:
speaker-versus-asserted-source attribution, cross-turn pronoun chains, and
explicit opinions and predictions, which the extraction passes above address
directly.

### Controlled answer accuracy and query-token frontier

All 16 conditions use the same Fireworks query embedding, local pgvector top-20 search, and DeepSeek V4 Flash answer model. Sources are packed greedily to 256, 512, 1,024, or 2,048 `cl100k_base`-estimated context tokens. The answer model must return structured citations; Fireworks reports actual input and output tokens. A fresh Kimi K2.6 session judges semantic correctness without seeing the representation or retrieved sources, while a separate Kimi session judges whether the cited source directly supports the answer.

| Representation | Budget | Correct | Cited support | Mean total tokens |
| --- | ---: | ---: | ---: | ---: |
| Facts + quote | 256 | 74/100 | 77/100 | 389.8 |
| Fixed chunks | 256 | 23/100 | 24/100 | 419.3 |
| Turn-aware chunks | 256 | 21/100 | 23/100 | 420.3 |
| Semantic chunks | 256 | 28/100 | 29/100 | 404.9 |
| Facts + quote | 512 | 81/100 | 83/100 | 654.4 |
| Fixed chunks | 512 | 46/100 | 44/100 | 663.6 |
| Turn-aware chunks | 512 | 45/100 | 46/100 | 663.9 |
| Semantic chunks | 512 | 45/100 | 44/100 | 631.2 |
| Facts + quote | 1,024 | 87/100 | 86/100 | 1,158.8 |
| Fixed chunks | 1,024 | 60/100 | 59/100 | 906.1 |
| Turn-aware chunks | 1,024 | 56/100 | 57/100 | 921.1 |
| Semantic chunks | 1,024 | 61/100 | 61/100 | 870.1 |
| Facts + quote | 2,048 | 86/100 | 87/100 | 1,407.5 |
| Fixed chunks | 2,048 | 64/100 | 64/100 | 1,801.3 |
| Turn-aware chunks | 2,048 | 59/100 | 62/100 | 1,806.5 |
| Semantic chunks | 2,048 | 63/100 | 63/100 | 1,901.1 |

The clearest matched comparison is 512 tokens: facts are 81% correct, 35–36 points above all three chunk controls, at nearly identical actual answer-model cost. The frontier comparison is stronger: 256-token facts are 74% correct, above every 512-token chunk condition (45–46%), while consuming 38–41% fewer total answer-model tokens. Similarly, 1,024-token facts are 87% correct and remain cheaper than every 2,048-token chunk condition (59–64% correct). This is evidence for an accuracy–cost frontier, not merely a retrieval difference.

### What this does and does not establish

The present evidence supports a focused hypothesis: a one-time, provenance-preserving extraction and rephrasing step may let a RAG system answer more dialogue questions from substantially less context, particularly questions involving who said something, what a pronoun refers to, and a narrow claim embedded in a long conversation. It does **not** yet establish a publishable general result or a dollar-denominated return on the ISC cost.

Before making a paper claim, we will: (1) obtain independent human labels for the blinded 80-item calibration worksheet and compare them with the Kimi screens; (2) freeze the prompt, packing, and judgment protocol; and (3) repeat on untouched, document-disjoint data. Only then should we estimate the extraction/indexing break-even against query-time savings.

## Reproducibility pointers

- Full controlled frontier: [`fact-development-answer-frontier-20260725.md`](fact-development-answer-frontier-20260725.md)
- Blinded human-calibration worksheet: [`human-calibration-worksheet.md`](../experiments/fact-development-answer-frontier-20260725/human-calibration-worksheet.md)
- Qualitative fact-win audit and metric caveats: [`fact-development-sample-100-fact-wins.md`](fact-development-sample-100-fact-wins.md)
- Frozen development sample and boundary: [`README.md`](../README.md)
