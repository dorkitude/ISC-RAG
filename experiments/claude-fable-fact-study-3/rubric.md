# Locked evaluation rubric — Claude Fable fact study #3

You are evaluating one candidate fact extraction, blind. Do not guess which
configuration produced it, do not compare to any other run, and do not use any
knowledge outside the files you were given. Judge only what is on the page.

## Inputs you were given

1. `validated.tsv` — one fact per line:
   `index, fact_text, primary_quote, antecedent_quote_or_-, primary_line, antecedent_line_or_-, speaker`
   (tab-separated). The quotes are already verified verbatim substrings of the
   cited transcript lines; the speaker is attached deterministically from the
   primary line. You do not need to re-verify substring presence.
2. `labeled.tsv` — the full transcript, one turn per line, `speaker<TAB>text`.
   Line numbers are turn numbers.
3. The frozen claim ledger for this document (`L<n>, turn_line(s), claim`).
4. This rubric.

## Per-fact binary judgments (1 = pass, 0 = fail)

1. **support** — Every element of `fact_text` (entities, predicate, numbers,
   dates, places, causal links) is stated in or strictly entailed by the
   fact's own quoted evidence (`primary_quote` plus `antecedent_quote` if
   present), read with the canonical speaker labels of their lines. The
   surrounding transcript may only *disqualify* a fact (e.g. the quote is a
   question, a hypothetical, sarcasm, or a misreading of who says what); it may
   never rescue an element missing from the quotes. Outside knowledge never
   counts as support.
2. **attribution** — The fact correctly preserves who asserted or reported the
   claim. If the fact text names a source or speaker, it must match the
   canonical speaker or the reported source inside the quote. If the claim is a
   stance, prediction, estimate, or first-person report, the fact text itself
   must attribute it to the correct person; a bare world-fact phrasing fails
   attribution for such claims. Neutral world facts asserted by a reporter in
   their own voice may rely on the attached speaker metadata alone.
3. **qualifiers** — All negation, modality and hedges ("may", "expected to",
   "I think"), temporal scope, quantities, approximations ("about 10 percent"),
   conditions, and reported-source framing ("officials said") present in the
   evidence are preserved in the fact text, with no strengthening, weakening,
   or dropped scope.
4. **standalone** — The fact text is fully understandable with the transcript
   hidden. Fail if it contains unresolved pronouns or deictics ("he", "this
   report", "there", "that decision"), relative time without an anchor
   ("today", "next week" with no date context), dangling comparatives ("more
   than before"), or depends on an unstated question.
5. **atomicity** — Exactly one predicate-level claim. Compound facts joining
   independent claims fail. A single claim with its own integral qualifier or
   attribution wrapper still passes.

## Duplicates

For each fact, mark `duplicate_of` = the index of the earliest fact making the
same claim about the same subject (near-duplicates and paraphrases count), or
`-` if it is the first/unique.

## Ledger coverage

For every ledger claim `L<n>`, mark covered = 1 iff at least one validated fact
expresses that claim's core content with compatible attribution and qualifiers
(a partial or weakened rendering does not count). List the covering fact
indices.

## Required output format (exactly these three blocks)

```
FACTS
<index>\t<support>\t<attribution>\t<qualifiers>\t<standalone>\t<atomicity>\t<duplicate_of>\t<short reason if any 0, else ->
...one line per fact...
COVERAGE
<ledger id>\t<covered 0/1>\t<covering fact indices comma-separated or ->
...one line per ledger claim...
SUMMARY
validated=<n> support_fail=<n> attribution_fail=<n> qualifiers_fail=<n> standalone_fail=<n> atomicity_fail=<n> duplicates=<n> ledger_total=<n> ledger_covered=<n>
```

Be strict. When uncertain on any binary judgment, score 0 and say why in the
reason column. Do not add commentary outside the three blocks.
