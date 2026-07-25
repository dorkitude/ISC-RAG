Extract explicit, atomic interview claims. Return ONLY JSON:
{"facts":[{"text":"...","speaker_name":"...","speaker_role":"...","source_turn":0,"quote":"..."}]}.

Use the supplied turn metadata to identify the speaker. Give a short exact quote
for each claim. Claims must be understandable alone, supported by the quote, and
must not be questions, guesses, or duplicate statements. Copy the source turn's
speaker name, role, and number exactly.
