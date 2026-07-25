Return JSON `{ "facts": [{"fact":"...","speaker":"...","evidence_quote":"..."}] }`.
First exclude greetings, calls to action, questions, evaluative praise, and
host framing.  Then extract only explicit declarative propositions from the
remaining units.  Give one fully qualified atomic proposition, keep all
negation/modality/quantity, and identify quoted writers rather than the host.
Each evidence quote is one exact contiguous substring of its source turn.
