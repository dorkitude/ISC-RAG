# Top-k answer-depth pilot: five frozen MediaSum documents

**Status:** five-question feasibility pilot, not a paper result. Its purpose is
to validate the depth-saturation protocol and expose how retrieval-unit size
interacts with top-k and query cost before we scale it to the frozen 100.

## Question

For a fixed question and retrieval representation, what is the smallest number
of ranked sources (`top_k`) needed for a correct, source-supported answer? Does
that depth translate to lower query-token use for facts than for chunks?

## Protocol

The pilot uses five documents from the already frozen 100-document development
set, recorded in
[`fact-development-depth-pilot-5-20260725.txt`](../examples/fact-development-depth-pilot-5-20260725.txt).
They contribute one existing Fable-authored, source-validated question from
each family: attribution, decontextualization, direct control,
opinion/belief/judgment, and pronoun/coreference.

For facts, fixed chunks, turn-aware chunks, and semantic chunks, the evaluator
retrieves the same ranked results and runs the DeepSeek V4 Flash answer model at
`top_k = 1, 3, 5, 10, 20, 30`. There is **no context packing cap** in this
pilot: every source in the specified top-k is passed to the answer model. The
model returns structured source IDs. Separate blind Kimi K2.6 sessions judge
answer correctness and whether the cited source directly supports the answer.
A success below means both judgments are positive.

This deliberately measures depth saturation, but it makes raw `k` an unfair
cost proxy: one fixed chunk is far larger than one fact. We therefore report
actual provider-reported answer-model tokens alongside k.

## Answerable-and-supported curve

| Representation | k=1 | k=3 | k=5 | k=10 | k=20 | k=30 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Facts | 2/5 | 3/5 | 3/5 | 3/5 | 4/5 | 3/5 |
| Fixed chunks | 3/5 | 3/5 | 3/5 | 3/5 | 4/5 | 4/5 |
| Turn-aware chunks | 3/5 | 3/5 | 3/5 | 4/5 | 4/5 | 4/5 |
| Semantic chunks | 2/5 | 3/5 | 3/5 | 3/5 | 4/5 | 4/5 |

The curves are not strictly monotonic. More context can introduce distractors,
and this is one generation per cell rather than repeated sampling. The key
analysis is therefore each question's **minimum** successful k, with `never`
meaning that no correct, supported answer occurred through k=30.

| Question family | Facts | Fixed | Turn | Semantic |
| --- | ---: | ---: | ---: | ---: |
| Pronoun/coreference | 1 | 20 | 10 | 20 |
| Attribution | 10 | 1 | 1 | 3 |
| Decontextualization | 1 | never | never | never |
| Opinion/belief/judgment | never | 1 | 1 | 1 |
| Direct control | 3 | 1 | 1 | 1 |

This tiny sample reveals the intended mechanism in two cases: the rephrased
fact condition resolves the coreference question at k=1 while chunks need
k=10–20, and it answers the decontextualization question at k=1 while no chunk
condition produces a correct, supported answer through k=30. It also reveals
counterexamples: the attribution, opinion, and direct-control questions are
easier for a top-1 chunk in this sample. The right conclusion is not that
facts universally need lower k; it is that the full study should analyze the
distribution of first successful k by question family.

## Cost makes the difference visible

| Representation | Mean total tokens @1 | @3 | @10 | @20 | @30 |
| --- | ---: | ---: | ---: | ---: | ---: |
| Facts | 229 | 354 | 775 | 1,415 | 1,920 |
| Fixed chunks | 865 | 2,131 | 6,846 | 14,067 | 21,175 |
| Turn-aware chunks | 690 | 1,984 | 6,910 | 13,721 | 20,976 |
| Semantic chunks | 631 | 1,722 | 5,190 | 10,072 | 15,519 |

Even when a chunk succeeds at k=1, it starts at roughly 3–4× the token cost of
a top-1 fact in this pilot. This is why the scaled experiment must report both
minimum successful k and actual tokens at that point. A chunk can have a lower
k yet still be materially more expensive to provide to an answer model.

## Reproducible CLI support

`isc-rag evaluate run` now accepts `--document-list FILE`, selecting only
questions whose gold evidence belongs to the listed document IDs. The new
`isc-rag evaluate depth-summary --runs UUID,...` command reports both the
success curve and per-question first successful k from stored answer and judge
records.

## Next scaled protocol

On the frozen 100 questions, run the same six depths and calculate: (1)
coverage at each k; (2) median/mean first successful k among answerable
questions; (3) fraction never answerable by k=30; (4) actual tokens at first
success; and (5) paired results by question family. Repeat generations or add
human calibration before treating non-monotonic differences as substantive.
