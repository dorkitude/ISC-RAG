# Second-question replication set

**Status:** development-only retrieval replication; not a document-disjoint or
paper-facing test set.

This set adds one new Claude Fable 5 question for each document in the frozen
100-document list. Fable saw only the canonical transcript, document ID, a
requested question family, and the prior question text solely to select a
different claim/topic. It never saw facts, chunks, retrieval outputs, or the
earlier retrieval results.

Question-family requests rotate evenly, yielding 20 each of direct control,
pronoun/coreference, attribution, decontextualization, and
opinion/belief/judgment. Before admission, the generator script verifies that
the stated document matches the target document, the inclusive cited turn range
exists, and the supplied supporting quote occurs exactly in a cited canonical
turn. The tracked `questions.jsonl` is exported from those validated PostgreSQL
records.

Generation is resumable and idempotent via
`scripts/generate-fable-question-expansion.sh`. This is a question-level
replication on the same development documents, not evidence of generalization
to previously unseen documents.
