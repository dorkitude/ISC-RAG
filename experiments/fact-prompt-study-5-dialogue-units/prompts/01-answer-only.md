Return a JSON object with `facts`, an array of at most 24 objects.  For every
object return `fact`, `speaker`, and `evidence_quote`.  Extract only explicit
atomic claims from quoted letters or declarative narration.  Use the answer or
statement turn as evidence; do not convert questions or host introductions into
claims about the host.  Resolve a name only if it occurs in the evidence turn.
`evidence_quote` must be an exact contiguous substring of that turn.
