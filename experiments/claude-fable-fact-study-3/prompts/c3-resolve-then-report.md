# Config C3 — "Resolve-then-Report" (structural peak-jump)

Proposed by a fresh study_designer subagent as a design structurally distinct
from C1 (compact single-quote baseline) and C2 (contract + short backward
antecedent window). Two structural moves: (1) coreference resolution is split
into its own document-global Pass A that emits a provenance-carrying Referent
Ledger; (2) every fact is framed as an attribution report so that the mandatory
first evidence span (the asserting turn) both supports the claim and drives
deterministic speaker attribution, which lets later anchor spans reach anywhere
in the document without ever stealing attribution.

Executed by ONE fresh Claude Fable extractor subagent that performs both passes
internally and self-validates every span with `grep -Fq`. Deterministic
validation (grep, by the mechanical validator) then re-checks all spans.

Evidence schema for C3 (up to three spans; span 1 mandatory):
`index, fact_text, span1_quote, span2_quote_or_-, span3_quote_or_-, span1_line, speaker`

Deterministic validation rules for C3:
- span1_quote must be a verbatim substring of the turn line span1_line, ≤240 chars.
- span2_quote / span3_quote, when present, must each be a verbatim substring of
  SOME turn line anywhere in the document (global reach), ≤240 chars.
- speaker is attached from span1_line's canonical label.
- exact-duplicate fact_text dropped (first kept).
Support is judged by the evaluator against the UNION of all cited spans.

---

## Extractor prompt (— {LABELED}/{TURNS} substituted per document —)

You are compiling provenance-exact, standalone fact cards from a broadcast-
dialogue transcript, in two passes. Read ONLY these two files:
- labeled transcript (`SPEAKER<TAB>text`, Read line number = turn number): {LABELED}
- same transcript, one turn per line, no labels (line N == turn N): {TURNS}

### Pass A — build a Referent Ledger for the whole document
Scan the ENTIRE transcript. For every referring expression (a pronoun, "this
report", "the program", "that decision", "there", a bare last name, etc.) whose
identity is fixed elsewhere by an explicit verbatim mention, record one ledger
row with:
- canonical: the fully resolved noun phrase, containing NO content word that is
  not grounded in an anchor quote below;
- anchor_quote(s): 1–2 exact contiguous substrings (each from one turn line,
  ≤240 chars) that state the canonical identity (e.g. a name in an early turn,
  a role in another). Anchors must be locally unique multi-word spans, never a
  bare pronoun.
Include a referent ONLY when its binding is unambiguous; when in doubt, omit it
(protect support over recall). Verify each anchor_quote with
`grep -Fq "<anchor>" {TURNS}`; drop rows whose anchor is not found.

### Pass B — attribution-framed, coverage-first harvest
1. Walk every turn top to bottom and enumerate EVERY explicit declarative
   assertion (world facts/events, stated positions, predictions, estimates,
   explanations, first-person experiences). Be exhaustive; over-harvest — the
   orchestrator prunes duplicates deterministically. Skip greetings, thanks,
   host housekeeping, bare questions, filler.
2. Write each fact as `<exact source speaker label or name> <attribution verb
   matching the source: says/said/asserts/predicts/estimates/recalls/reports/
   alleges> <claim>` when the claim is a stance/prediction/estimate/first-person
   report or otherwise needs an asserting subject. A world fact a reporter
   asserts in their own voice may be stated directly.
3. Decontextualize the claim using ONLY the Referent Ledger: replace each
   internal pronoun/deictic (it/they/that/there/then/this X) with the matching
   row's canonical. Do NOT invent context. Preserve every negation, hedge,
   quantity, condition, temporal scope, and reported-source framing exactly.
   Exactly one predicate-level claim per card (split compounds).
4. Evidence spans, in THIS ORDER:
   - span1 (mandatory) = the asserting turn's verbatim quote — supports the
     claim and anchors attribution;
   - span2, span3 (optional) = the anchor_quote(s) of any ledger rows you used
     to resolve a referent, giving verbatim grounding for the substitution.
   Every element of the fact text must be entailed by the union of its spans
   read with their speaker labels. If a needed referent is not groundable, keep
   the short source-literal form or drop the card — prefer dropping to guessing.
5. Verify EVERY span with `grep -Fq "<span>" {TURNS}` before output; fix to an
   exact copy or drop. span1 must come from the asserting turn line; attach the
   speaker from that line's SPEAKER label (never invent it).

### Output
Write with the Write tool to {OUTFILE}, one card per line, tab-separated, no
header, columns:
`index<TAB>fact_text<TAB>span1_quote<TAB>span2_quote_or_-<TAB>span3_quote_or_-<TAB>span1_line<TAB>speaker`
(index from 1; unused span columns are a literal `-`.) Return at most 60 cards.
Also write the Pass A ledger you used to {LEDGERFILE} as JSON for audit.

Final message: report card count and confirm every span passed `grep -Fq`.
