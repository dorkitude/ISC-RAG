# Support-repair / rejection audit pass (frozen-config pipeline stage)

A fresh independent subagent applies the locked rubric's binary tests strictly
and DROPS any card that fails, before the final blind evaluation. This is a
structural rejection stage (analogous to the production pipeline's post-extraction validation),
not a prompt tweak, added because the gate's "support ≥98% AND zero unsupported"
is a knife-edge that a single over-reaching card fails. It trades a little
coverage (C3 over-harvests, leaving headroom above the 75% floor) for reliable
support/standalone/atomicity.

It is INDEPENDENT of the final `fact_evaluator` (different fresh subagent) and
is instructed to be at least as strict, so the surviving set should contain
nothing the evaluator would flag. It removes cards only; it never edits text or
provenance (keeping facts immutable per the contract), so it cannot fabricate
support.

## Prompt (paths substituted per run)

You are a strict rejection auditor. Read only:
- Rubric: /Users/.../experiments/claude-fable-fact-study-3/rubric.md
- Validated facts (multi-span: index, fact_text, span1_quote, span2_quote_or_-, span3_quote_or_-, span1_line, speaker): {VALIDATED}
- Transcript (`SPEAKER<TAB>text`, line=turn): {LABELED}

The fact's evidence is the UNION of its span1/span2/span3 quotes (already verified
verbatim), read with the span1_line speaker. For EACH card apply the rubric's
five binary tests exactly:
- support: every element of fact_text entailed by the union of its spans; outside
  transcript may only disqualify, never rescue; unresolved-by-quote entity = fail.
- attribution, qualifiers, standalone, atomicity as defined in the rubric.
KEEP a card only if it passes ALL FIVE. DROP it otherwise. Also drop a card whose
fact_text exactly duplicates an earlier KEPT card. Do NOT edit any surviving card.

Write survivors verbatim (same 7 columns, re-indexed from 1) to {KEPT}. Write a
{AUDITREPORT} listing total, kept, dropped, and `DROP <orig_index> <failed
criterion> <one-line reason>` per drop. Final message: kept/dropped counts only.
