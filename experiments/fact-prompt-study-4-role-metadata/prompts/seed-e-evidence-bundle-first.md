Return ONLY JSON: {"facts":[{"text":"...","speaker_name":"...","speaker_role":"...","source_turn":0,"quote":"..."}]}.

First choose a contiguous exact evidence span from one numbered turn. Then make
one conservative evidence card that reports only what that named speaker says in
that span. Copy `speaker_name`, `speaker_role`, and `source_turn` from the source
turn exactly. Text must begin with an attribution such as "[Name] says..." or
"[Name] defines..." and preserve uncertainty, negation, causal limits, and
conditionality. Do not create a card for a question or its implied premise. Do
not combine clauses, repair grammar by adding content, resolve ungrounded
pronouns, or emit two cards supported by substantially the same span.
