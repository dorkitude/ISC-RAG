# Config C1 — direct baseline (compact contract, no antecedent mechanism)

The extraction subagent receives exactly this prompt, with {LABELED} replaced by
the document's labeled transcript path.

---

You are extracting explicit, atomic facts from a broadcast-dialogue transcript.

Read EXACTLY ONE file and nothing else: {LABELED}
Each line is one speaker turn: `SPEAKER<TAB>text`. The line number shown by the
Read tool is the turn number.

Extract only claims explicitly stated in the transcript. Do not speculate,
summarize conversational dynamics, infer unstated context, or write
meta-claims. Every fact must be a standalone atomic claim understandable on its
own. Each fact needs one verbatim supporting quote copied exactly,
character-for-character, from a single turn line, under 240 characters.

Output contract: your final message must be ONLY fact lines, one fact per line,
tab-separated, with no header and no commentary:

fact_text<TAB>primary_quote<TAB>-<TAB>turn_line_number

Return at most 60 facts.
