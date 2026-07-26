# Fact-development sample-100 question set

**Status:** development-only; never use as a held-out or paper test set.

This set contains one source-grounded question for each document in
`examples/fact-development-sample-100-20260725.txt`, the frozen deterministic
100-document fact-development sample. It is balanced by construction: 20 each
of direct controls, pronoun/coreference, attribution, decontextualization, and
opinion/belief/judgment.

Claude Fable 5 (`claude -p --model fable`) authored one candidate per
document in fresh sessions from canonical Postgres turns only. Raw completions
and their supporting quotes remain in ignored artifacts. Before this tracked
file was created, every row was mechanically checked against Postgres: its
document exists, every ordinal in the inclusive cited range exists, and the
raw candidate's supporting quote occurs exactly in at least one cited turn.

The tracked JSONL intentionally keeps the question, reference answer, and
canonical evidence range, but omits transcript quotations. This question set
may be used to improve the fact pipeline and retrieval configuration for this
frozen sample. Create a document-disjoint, untouched held-out set before any
paper-facing result is reported.
