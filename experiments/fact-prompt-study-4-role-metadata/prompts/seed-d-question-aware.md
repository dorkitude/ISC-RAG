Return ONLY JSON: {"facts":[{"text":"...","speaker_name":"...","speaker_role":"...","source_turn":0,"quote":"..."}]}.

This is an interview, not a factual article. The interviewer may ask loaded,
hypothetical, or compound questions. Never extract the premise of a question as
a fact. Extract claims only from declarative answer content or clearly marked
host narration, and state who made the claim. For every answer, preserve its
speech act and scope: retain "I think", "not likely", "could", "may", "if", and
time qualifiers. Use deterministic speaker name and role copied exactly from the
quoted turn. Make one proposition per card; the quote must be a contiguous exact
substring from that same turn. If a pronoun or ellipsis lacks a local explicit
antecedent, omit it rather than resolving it.
