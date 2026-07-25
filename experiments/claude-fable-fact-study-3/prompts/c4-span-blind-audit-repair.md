# Config C4 — "Span-Blind Audit and Repair" (C4-SBAR)

Proposed by a fresh study_designer subagent (agentId a9280f273fc6ed5b3) after
the r03 audit. Full designer spec (verbatim, includes transcript-quoting
failure analysis) is retained in ignored
`artifacts/claude-fable-fact-study-3/design/c4-designer-spec.md`; this tracked
file contains the frozen executable configuration with no transcript text.

Structural rationale: C1–C3 enforced every quality axis inside one generation
prompt. The audit shows support is the wrong axis to enforce at generation time
(the extractor knows the transcript, so it cannot perceive an ungrounded
referent; C3 proved that generation-time groundability constraints degrade the
generator). C4 leaves the best generator (C2) byte-identical and enforces
support post-hoc with an information-asymmetric auditor that sees ONLY what the
rubric's support judgment sees — the fact's own quoted spans plus speaker
labels, never the transcript — so an ungrounded element is invisible-as-grounded
to it exactly as it is to the blind evaluator. Audit-failed cards go to a
targeted repair stage (transcript access; ordered moves REANCHOR → RECUT →
REWRITE → DROP; drop always beats a deictic), repairs are re-validated
mechanically and re-audited span-blind, and the final TSV is assembled by a
mechanical keep-iff-PASS rule. Repair blast radius is bounded to audit-failed
cards, protecting the axes C2 already passes.

## Pipeline (frozen; all stage agents fresh, no study history)

| Stage | Agent | Receives | Emits |
|---|---|---|---|
| S1 Extract | fresh Fable extractor | `{LABELED}` only; prompt = `prompts/c2-contract-antecedent.md` **byte-identical** | `raw.tsv` (4 cols per C2) |
| S2 Validate | mechanical validator (A1, grep -Fq only) | raw.tsv, `{TURNS}`, `{SPEAKERS}` | `stage1.tsv` (7-col validated schema, index = extraction order) |
| S3 Blind audit | fresh Fable auditor (Prompt A) | stage1.tsv + `{SPEAKERS}` only — **no transcript** | `audit1.tsv`: `index<TAB>PASS\|FAIL<TAB>missing_elements_or_-` |
| S4 Repair | fresh Fable repairer (Prompt R) | `{LABELED}`, `{TURNS}`, `failed.tsv` (FAIL rows of stage1.tsv + audit reasons) | `repairs.tsv`: REPAIR/DROP decision lines |
| S4b Validate repairs | mechanical validator | repairs.tsv REPAIR rows, `{TURNS}`, `{SPEAKERS}` | `repairs.validated.tsv` (grep-failing rows become DROP) |
| S5 Blind re-audit | fresh Fable auditor (Prompt A, fresh instance) | repairs.validated.tsv + `{SPEAKERS}` only | `audit2.tsv` |
| S6 Assemble | mechanical validator | all above | final `validated.tsv` + validation-report.txt |

Frozen skip rule: if S3 returns zero FAILs, S4–S5 are skipped.

## Deterministic validation (mechanical only; A1 execution)

1. primary_quote verbatim substring (grep -Fq, exit status only) of its cited
   turn line, ≤240 chars; else drop.
2. antecedent_quote (when not `-`) verbatim substring of one line in
   `[primary_line − 4, primary_line]`, ≤240 chars; else drop (first match wins
   for the recorded antecedent line).
3. Speaker attached from the primary line's canonical label; never from model
   output.
4. Exact-duplicate fact_text dropped, first kept.
5. Assembly (mechanical lookup only): a stage1 row survives iff its index has
   PASS in audit1.tsv; a repaired row substitutes at its original index iff it
   has REPAIR in repairs.tsv, survived S4b re-grep, and has PASS in audit2.tsv;
   all other indices deleted. Then rules 1–4 re-run once over the assembled
   file; reindex 1..n preserving original order.

Final validated.tsv schema is identical to C2/r02, so the evaluator's input
contract is unchanged.

## Prompt A — span-blind support auditor (S3, S5; placeholders {CARDS}, {SPEAKERS}, {OUTFILE})

```
You are auditing fact cards for SUPPORT ONLY. You deliberately do NOT have the
source transcript. Read EXACTLY TWO files and nothing else:
- {CARDS} — one card per line, tab-separated:
  index, fact_text, primary_quote, antecedent_quote_or_-, primary_line,
  antecedent_line_or_-, speaker
- {SPEAKERS} — line N is the canonical speaker label of turn line N.

For each card, judge ONE question. Break fact_text into its elements: every
entity, predicate, number, date, place, descriptor, and causal link. The card
PASSES iff every element is stated in or strictly entailed by the card's own
quoted evidence — primary_quote plus antecedent_quote if present — read
together with the speaker labels of their lines (primary_line's label is the
speaker column; look up antecedent_line's label in {SPEAKERS}). Otherwise it
FAILS.

Rules:
- The quotes are the entire universe of evidence. Outside knowledge,
  plausibility, and what "must have been meant" never count. If a content word
  of fact_text (a name, a topic word, a descriptor, a possessive) is neither
  present in nor strictly entailed by the quoted text, FAIL and name it.
- Paraphrase is fine: an element counts as supported if the quotes state it in
  different words ("the plan was rejected" supports "the proposal was turned
  down"). Do not demand lexical identity — demand presence or strict
  entailment of the element.
- Speaker labels are evidence: "I" in a quote is that line's labeled speaker;
  "you" in a two-person exchange is the other line's labeled speaker when both
  labels are before you.
- Judge nothing else. Pronouns in fact_text, style, atomicity, duplication are
  NOT your concern. Support only.
- When genuinely uncertain whether an element is entailed, FAIL and name the
  element. A false FAIL costs little; a false PASS is fatal.

Write with the Write tool to {OUTFILE}, one line per card, no header, no
commentary:
index<TAB>PASS<TAB>-
index<TAB>FAIL<TAB><comma-separated missing/ungrounded elements>
Every input index must appear exactly once. Final message: counts only
("N cards, K FAIL").
```

Orchestrator neutralization note: the designer's Prompt A used a paraphrase
example echoing NPR-1394 wording; it was replaced with the neutral
plan/proposal example above before freezing so tracked files stay free of
transcript-derived text. No rule content changed.

## Prompt R — targeted repairer (S4; placeholders {LABELED}, {TURNS}, {FAILED}, {OUTFILE})

```
You are repairing fact cards that failed a strict support audit. Each card
asserts at least one element that is NOT present in its quoted evidence. Your
job is to make the evidence carry the fact, or remove the fact. Read EXACTLY
THREE files:
- {LABELED} — transcript, one turn per line, SPEAKER<TAB>text; Read line
  number = turn number.
- {TURNS} — same turns, no labels, for grep verification.
- {FAILED} — the failed cards: index, fact_text, primary_quote,
  antecedent_quote_or_-, primary_line, speaker, missing_elements.

For each card, output exactly one decision line. primary_line is FROZEN — it
is the asserting turn and may never change. Apply the first move that fully
grounds every missing element:

1. REANCHOR — keep fact_text unchanged; replace the antecedent_quote with a
   verbatim contiguous span (≤240 chars) from ONE line in
   [primary_line − 4, primary_line] that states the missing element. Re-cutting
   the existing antecedent to a longer contiguous span of its own line counts.
2. RECUT — extend or shift primary_quote to a longer verbatim contiguous span
   of the SAME primary line that includes the missing element. Never cut so
   that a negation, hedge, or quantity is lost.
3. REWRITE — minimally edit fact_text so every element is entailed by the
   (possibly re-anchored/re-cut) quotes: generalize the ungrounded descriptor
   to exactly what the quotes state, or delete it. The result must still be
   fully standalone (NO pronouns, NO deictics like "this authority"/"those
   reasons"/"that case", no relative time without an anchor, no dangling
   comparatives), still one atomic claim, with every hedge, negation,
   quantity, condition, reported-source framing, and attribution preserved
   exactly. You may not change what the claim says beyond removing the
   ungrounded element.
4. DROP — if no combination of moves 1–3 yields a card that is BOTH fully
   grounded in its two quotes AND fully standalone. Dropping is always
   acceptable; keeping an ungrounded card or a deictic card never is.

Verify EVERY quote you output with grep -Fq "<quote>" {TURNS} before writing;
if a quote fails, fix it to an exact copy or DROP the card.

Write with the Write tool to {OUTFILE}, one line per input card, tab-separated,
no header:
index<TAB>REPAIR<TAB>fact_text<TAB>primary_quote<TAB>antecedent_quote_or_-<TAB>primary_line
index<TAB>DROP<TAB>-<TAB>-<TAB>-<TAB>-
Every input index appears exactly once. Final message: counts only.
```

## Designer's failure-mode analysis (paraphrased, quote-free)

- All five r02 support fails are of the class "referent/descriptor asserted but
  absent from both spans"; S3 catches the whole class by construction because
  the missing token is absent from the auditor's entire evidence universe.
  Expected disposition on the dev doc: most repair via REANCHOR (a nearby turn
  states the missing referent inside the antecedent window), remainder DROP;
  coverage margin (93.2% vs 75% gate) absorbs sole-coverer drops.
- The 18 r03 standalone fails are not inherited: they came from C3's
  generation-time groundability constraint, which C4 deletes by running the C2
  generator verbatim. Prompt R's move ordering and hard "drop beats deictic"
  rule forbid reintroducing them; only REWRITE can touch fact_text, and only on
  audit-failed cards (~10%), bounding worst-case standalone erosion to ~88–89%
  even if every rewrite failed standalone.
- Main risks: (1) auditor–evaluator disagreement (gate demands 0 unsupported;
  audit replicates the evaluator's information condition, not its instance);
  (2) over-strict auditing mass-failing sound paraphrases and eroding coverage
  via drops; (3) blind spot: a spans-only auditor cannot see transcript-level
  disqualifiers (quote that is in-context a question/hypothetical/unmarked
  reported speech) — first thing to check if a replicate shows a surprise
  unsupported fact.
