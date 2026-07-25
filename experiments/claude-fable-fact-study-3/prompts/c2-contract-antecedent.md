# Config C2 — explicit contract + antecedent evidence window

Structural changes versus C1: (a) a detailed extraction contract with failure
taxonomy and worked examples; (b) an antecedent-quote slot so a fact may cite a
second verbatim span from up to 4 turns earlier that resolves a referent —
support is judged against both quotes, which removes the standalone-versus-
single-quote conflict identified in the prior DeepSeek studies; (c) a
coverage-first sweep instruction.

The extraction subagent receives exactly this prompt, with {LABELED} replaced by
the document's labeled transcript path.

---

You are compiling provenance-exact fact cards from a broadcast-dialogue
transcript.

Read EXACTLY ONE file and nothing else: {LABELED}
Each line is one speaker turn: `SPEAKER<TAB>text`. The line number shown by the
Read tool is the turn number.

Sweep the whole transcript in order and produce one fact card for every
explicit substantive claim: world facts and events, a speaker's stated
position, prediction, estimate, explanation, or first-person experience. Skip
greetings, thanks, host housekeeping, bare questions, and filler.

Rules for the fact text:
1. Standalone: a reader with the transcript hidden must fully understand it.
   Resolve every pronoun and deictic ("he", "this report", "there", "the
   program") to its explicit named referent from the transcript. Convert
   relative time to the explicit anchor stated in the transcript, or drop the
   time phrase if none is stated. Never leave a dangling comparative.
2. Atomic: exactly one predicate-level claim per card. Split compound
   statements into separate cards.
3. Faithful qualifiers: preserve every negation, hedge ("may", "I think",
   "expected to"), quantity ("about 10 percent"), condition, temporal scope,
   and reported source ("officials said") exactly; never strengthen or weaken.
4. Attribution: if the claim is a stance, prediction, estimate, or first-person
   report, name the speaker in the fact text ("<Speaker> says/estimates/
   believes ..."). A world fact a reporter asserts in their own voice may be
   stated directly without naming the reporter.
5. Explicit only: no inference, no outside knowledge, no meta-claims about the
   conversation, nothing from questions unless the answer asserts it.

Evidence rules:
- primary_quote: verbatim, character-for-character copy from ONE turn line,
  under 240 characters, containing the core assertion including its operator
  ("I think", "not", "about") — never cut the quote so a qualifier is lost.
- antecedent_quote: when the fact text resolves a pronoun/referent or needs
  context stated earlier, copy the verbatim span that states that antecedent
  from the SAME line or one of the 4 preceding lines. Otherwise write `-`.
- Every element of your fact text must be entailed by these one or two quotes
  read together with their speaker labels. If you cannot quote the evidence,
  do not emit the card.

Worked examples (patterns, not from this transcript):
- BAD: "He said the program was failing." (unresolved "he", "the program")
  GOOD: "John Smith said the Head Start program was failing." with
  primary_quote "the program was failing" and antecedent_quote naming Head
  Start from an earlier line.
- BAD: "Sanctions were imposed and oil exports fell." (two claims)
  GOOD: two cards, one per claim.
- BAD: "The economy will collapse." when the source says "I think the economy
  may collapse" (dropped hedge and attribution)
  GOOD: "Jane Doe says she thinks the economy may collapse."

Output contract: your final message must be ONLY fact lines, one card per line,
tab-separated, with no header and no commentary:

fact_text<TAB>primary_quote<TAB>antecedent_quote_or_-<TAB>primary_turn_line_number

Return at most 60 cards.
