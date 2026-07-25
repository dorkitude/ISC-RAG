# Two-pass DeepSeek V4 Flash check

**Status:** exploratory engineering check, 2026-07-24. This is not a paper
result and does not promote a primary fact configuration.

## Question

Can Fireworks `accounts/fireworks/models/deepseek-v4-flash` improve the
self-containedness of existing source-validated dialogue claims when given the
claim, its exact quote, and a deterministic +/- two-turn, speaker-labeled
context window?

## Fixed material

The trial used the first five MediaSum test documents selected by the existing
CLI. They contain 95 `pilot-named-v3` claims. The first pass was frozen: these
are the existing isolated claims, not new Flash extraction. The only model work
was the second-pass contextual rewrite. Each completion was requested as JSON
with temperature zero and `reasoning_effort=none`.

Raw transcript-derived results are intentionally ignored and checkpointed after
every completion:

- `artifacts/two-pass-rephrasing/flash-5docs.json` — conservative prompt
- `artifacts/two-pass-rephrasing/flash-resolve-references-5docs.json` — role-
  reference prompt

## Conditions

1. **Conservative.** Resolve a pronoun only when a nearby explicit antecedent
   supports it; otherwise preserve the wording.
2. **Resolve references.** The same prompt plus permission to resolve a unique
   local role description (for example, “the president”) when explicitly
   identified.

Both runs used the same model, claims, quotes, context windows, output limit,
and deterministic document selection.

## Mechanical outcomes

| Measure | Conservative | Resolve references |
| --- | ---: | ---: |
| Completed rewrites | 95 / 95 | 95 / 95 |
| Provider/JSON failures | 0 | 0 |
| Output differs from frozen claim | 21 | 20 |
| `the speaker` / `the host` / `the guest` | 0 | 0 |
| Personal-pronoun surface matches | 20 | 23 |
| Selected ambiguous role-phrase matches | 10 | 10 |

The last two rows are screening counts, not accuracy scores: some remaining
pronouns and role phrases are legitimate because the supplied evidence does not
uniquely resolve them.

## Qualitative finding

The conservative condition made useful, source-supported repairs such as
resolving “his son, Bashar Assad” to “Bashar Assad,” “they” to “U.S. troops,”
and “my dad” to the named speaker's dad. The expanded condition sometimes named
President Bush or Congressman Randy Cunningham correctly, but it also reverted
some conservative pronoun repairs and added source attribution inconsistently.
It therefore did not lower either screening count and is not the preferred
prompt.

Use the conservative Flash prompt as the candidate for the next *independent,
source-blind* audit. Do not treat this comparison as proof of factual fidelity:
it has no human/independent correctness labels, and second-pass rewrites must
still pass quote/provenance validation before they become a retrieval condition.

## Reproduction

```sh
isc-rag two-pass-test --documents 5 --prompt-variant conservative \
  --output artifacts/two-pass-rephrasing/flash-5docs.json

isc-rag two-pass-test --documents 5 --prompt-variant resolve-references \
  --output artifacts/two-pass-rephrasing/flash-resolve-references-5docs.json
```
