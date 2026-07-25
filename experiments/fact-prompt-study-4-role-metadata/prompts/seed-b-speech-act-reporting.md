Return ONLY JSON: {"facts":[{"text":"...","speaker_name":"...","speaker_role":"...","source_turn":0,"quote":"..."}]}.

Build evidence cards for explicit interview speech. Treat the named speaker and
role on each turn as authoritative metadata. Prefer conservative reporting verbs
that preserve stance: "X defines...", "X says...", "X argues...", "X says that
X thinks...", or "X describes...". A question is not evidence for its premise;
make a card from it only when the question itself is the speech act of interest,
which it is not in this study. Do not silently convert opinions, forecasts,
possibilities, or proposed remedies into settled facts. Each card has exactly one
atomic proposition, an exact contiguous quote, and the exact metadata copied from
the quote's source turn. Omit ambiguous pronoun-dependent material rather than
guessing an antecedent.
