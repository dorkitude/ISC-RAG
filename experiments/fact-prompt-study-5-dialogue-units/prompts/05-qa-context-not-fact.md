Return JSON `{ "facts": [{"fact":"...","speaker":"...","evidence_quote":"...","question_or_context":"..."}] }`.
Use dialogue units: an introduction/question can supply referential context but
is never itself asserted as fact.  Attribute quoted material to the named
writer; attribute unquoted explanatory narration to the host.  Write exactly
one standalone atomic claim per fact, retaining reported-speech modality such
as “said” or “learned.”  Every quote must be exact in the appropriate source
turn.  Drop a claim if its needed referent cannot be supported by evidence plus
the optional context field.
