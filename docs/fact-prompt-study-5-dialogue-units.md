# Study #5: dialogue units and deterministic role metadata

**Status:** exploratory, one development document, not a paper condition.

## Question

Can making dialogue units and roles explicit cure the recurring fact-rewriting
failure on a transcript where a single host reads several other people's words?
This is deliberately a difficult provenance case: the physical utterance
speaker is the host, while the asserted source is sometimes a named letter
writer.  A system that stores only one `speaker` field must choose which role
it means.

## Frozen material and protocol

- **Document:** MediaSum validation record `NPR-1209`, distinct from the
  documents used by Studies 1--4 and from canonical question labels.
- **Ledger:** 16 source-only claims, frozen before model calls in
  [the Study #5 ledger](../experiments/fact-prompt-study-5-dialogue-units/claim-ledger.md).
- **Model/configuration:** `accounts/fireworks/models/deepseek-v4-pro`,
  temperature zero and no reasoning (the project CLI's fixed setting).
- **Evidence handling:** all source transcript, provider inputs/outputs, fact
  rows, and metric extracts are local ignored artifacts under
  `artifacts/fact-prompt-study-5-dialogue-units/`.  Tracked prompt concept
  specifications live in `experiments/fact-prompt-study-5-dialogue-units/prompts/`.
- **Important implementation disclosure:** the current shared Study-2 CLI has
  fixed candidate implementations, so provider executions below use the
  closest existing implementations rather than silently treating the tracked
  concept text as executable code.  Each exact effective prompt and SHA-256 is
  persisted in `study_two_runs`; the hashes below identify it.  This is
  a prompt/representation diagnostic, not a fully independent new pipeline.

## Evaluation rubric

I inspected retained cards against the frozen ledger and source units.  Support
means the returned evidence bundle supports the card; attribution requires both
roles to be correct when a host reads a letter: **utterance speaker** and
**asserted source**.  The existing schema contains only `speaker_label`, so a
card cannot pass this stricter two-role criterion merely by attaching the host
label.  Standalone means usable without the surrounding turn; capsule fragments
are not standalone.  Coverage counts ledger claims represented by a card that
also passes support, role attribution, and qualifier preservation.

The numbers are a single-run audit, rounded to whole-card fractions.  They are
diagnostic—not statistically comparable quality estimates.

| Run / tracked concept | Closest executed candidate | Cards | Support | Two-role attribution | Qualifier | Standalone | Atomic | Near duplicate | Ledger coverage | Tokens in/out; latency | Finding |
| --- | --- | ---: | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Answer-only | `direct` | 30 | 25/30 | 13/30 | 27/30 | 25/30 | 25/30 | 2/30 | 13/16 | 788/1259; 13.95s | High recall, but it turns reported opinions into bare world claims and loses the read-letter role. |
| Q+A attributed speech act | `dialogue-act-pairs` | 13 | 12/13 | 0/13 | 12/13 | 10/13 | 9/13 | 0/13 | 10/16 | 1164/1151; 26.85s | Better report wording, but every persisted deterministic speaker is the host; it cannot represent the named author separately. |
| Evidence-capsule-first | `quote-capsules-exact` | 13 | 13/13 | 0/13 | 13/13 | 0/13 | 10/13 | 0/13 | 0/16 | 1013/526; 6.57s | Exact provenance survives, but short quotes are fragments and the card credits the host rather than the quoted author. |
| Declarative-only filter | `prefilter-nli` | 22 | 20/22 | 0/22 | 20/22 | 21/22 | 20/22 | 1/22 | 0/16 | 3105/1692; 23.50s | The model explicitly produces `host said writer said ...`, demonstrating the needed two-layer representation, but the one speaker column cannot store it. |
| Q+A, question/context not fact | `dialogue-act-predicate-units` | 28 | 23/28 | 0/28 | 24/28 | 18/28 | 20/28 | 3/28 | 0/16 | 1943/1942; 19.52s | Predicate splitting restores recall but reintroduces fragments, duplicate material, and source-role collapse. |

`speaker_label` was mechanically `NEAL CONAN, host` for every output because
the MediaSum record marks the physical narration that way.  That is not a
provider error: it exposes a representational error in treating physical
utterance speaker as the author/claim-holder of quoted material.

## Evidence-driven refinement runs

The best seed for **semantic safety** was the capsule condition (no unsupported
rewriting), but it failed standalone usefulness and claim-holder attribution.
Rather than tweak wording blindly, the two follow-ups test the two observed
mechanisms:

1. **`substantive-render`** (13 cards; `2630/739` tokens; 2.81s; hash
   `413fb6bfcff53de851bbb5098ab6b6672c06336ce41485fdee0dd8909a828277`):
   source units are prefiltered before rendering.  It reduced output volume but
   did not add an asserted-source field, so quoted-letter claims remain unable
   to pass the two-role rubric.  This is a representation failure, not a
   missing instruction.
2. **`quote-repair`** (22 cards; `3014/1770` tokens; 22.44s; hash
   `9b2426feef5052e64ebc47bbf8629501f4bf8110d70edb1b9e4ca74392659639`):
   an evidence-first repair pass repairs unresolved spans.  It improved some
   local support/atomicity issues but still stores only the host speaker, so it
   cannot repair the defining provenance failure.

For reproducibility, the five seed effective-prompt hashes are: `direct`
`b95027d33569519941baa03a5f18c9faaa974bd9794ec608b44e17b405df6f40`,
`dialogue-act-pairs`
`0f3c8c2fc2ffaa5d93c3dd8172e9a80f9690f6d0e800873d3d32246a62e1f572`,
`quote-capsules-exact`
`b0e4b56c6b6245d92e4b1983bf2d3b85e58b2bec382b84e659c5d4ae146cfe9b`,
`prefilter-nli`
`70c30bd5bcf334c9e6ece2edfe5be7f5adfdf976b9c1a237a7259ac962967429`,
and `dialogue-act-predicate-units`
`146e038706bdf3f4725f30893b6a083217548d67c5f0d745b17015bf4f3c6e8f`.

## Honest stop decision

No condition is promotable.  The decisive result is structural: a RAG fact
record for dialogue needs at least `utterance_speaker`, `asserted_source` (or
`quoted_author`), and a speech-act/quotation relation.  Prompt changes cannot
recover an author role that the storage/validator overwrites with one generic
speaker field.  The next experiment should first add that two-role source model
and then compare source-faithful evidence capsules against retrieval-only
rephrasing; it should not run more prose prompt hill climbs in this schema.
