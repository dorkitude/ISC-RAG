Return JSON `{ "facts": [{"speaker":"...","evidence_quote":"...","context_quote":"...","retrieval_label":"..."}] }`.
Do not paraphrase a fact. Select one exact, self-contained atomic source span
from a declarative source unit; optionally select an exact contextual span from
the same turn.  `retrieval_label` is search-only and must not add information.
Do not select questions, greetings, or host framing as a claim.  A quoted
writer's words remain that writer's words even when read by the host.
