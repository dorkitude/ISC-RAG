# Canonical questions v1

This is a 100-question, source-grounded retrieval and answer-evaluation set for
ISC-RAG. Every item is drawn from an explicit statement in a MediaSum **test**
transcript. `corpus_split: "canonical"` identifies the question set, rather
than the underlying MediaSum split.

## Selection rules

- One-hop questions only: named entities, dates, numbers, definitions, simple
  lists, explicit attribution, and direct yes/no statements.
- Each answer is stated in the cited turn itself; the set does not use MediaSum
  titles or summaries as evidence.
- Evidence is one exact transcript turn, with its zero-based MediaSum turn
  indices recorded in `turn_start` and `turn_end`.
- Questions avoid opinion, inference, prediction, claims requiring outside
  knowledge, multi-document synthesis, and answers dependent on a speaker's
  identity unless the identity is stated in the same turn.
- The set spans 25 distinct NPR documents, with four questions per document.
  It deliberately samples news, science, politics, economics, sport, and
  lifestyle interviews rather than repeatedly drawing from one conversation.

## Deterministic answer evaluation

The stored `reference_answer` is the canonical target. For exact-match scoring,
normalize both candidate and target by Unicode NFKC normalization, lowercasing,
trimming, collapsing whitespace, and removing terminal punctuation. For
numeric answers, also normalize commas, currency symbols, and number words only
when they denote the identical stated value (for example, `6,000` and `6000`).
For list answers, normalize each comma- or `and`-separated item and compare as
an unordered set only when the question type is `list`.

Do not credit a semantically similar but differently qualified answer by
default: qualifiers such as `about`, `more than`, `at least`, ranges, units,
and negation are part of the reference answer. A separate, reported
lenient-normalization score may allow a leading `about`/`approximately` for
numeric answers only when the evidence itself uses that approximation.

Retrieval evidence recall should use the stored turn span directly. Generated
answers should be scored separately from evidence recall so that a correct
guess without retrieval does not count as a retrieval success.

## Coverage and caveats

The questions cover `NPR-12`, `NPR-454`, `NPR-759`, `NPR-1046`, `NPR-1466`,
`NPR-1683`, `NPR-2081`, `NPR-2099`, `NPR-2109`, `NPR-2349`, `NPR-2824`,
`NPR-2862`, `NPR-3013`, `NPR-3378`, `NPR-3499`, `NPR-3504`, `NPR-3528`,
`NPR-419`, `NPR-426`, `NPR-640`, `NPR-668`, `NPR-771`, `NPR-813`, `NPR-997`,
and `NPR-1026` in the pinned MediaSum release.

This is a deliberately deterministic canonical set, not a representative sample
of all dialogue comprehension. Its concentration of short, explicit factual
claims makes it suitable for comparing retrieval representations, but it does
not measure multi-hop reasoning, temporal reconciliation, speaker-resolution
across turns, or nuanced open-ended answers. Before reporting final paper
results, freeze this file by checksum and keep it disjoint from any data used
to tune chunking, prompting, extraction, or answer normalization.
