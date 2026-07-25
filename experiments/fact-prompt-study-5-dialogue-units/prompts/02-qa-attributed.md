Return JSON `{ "facts": [{"fact":"...","speaker":"...","evidence_quote":"...","context_quote":"..."}] }`.
Treat each quoted letter as an answer/report unit.  Render each output as an
attributed speech act: `[named writer] said/reported/believed ...`; preserve
uncertainty and scope.  A host turn that introduces a letter is context only,
not the speaker of its contents.  Emit one atomic claim per card.  Both quotes
must be exact substrings of their own source turns; omit `context_quote` when
unneeded.
