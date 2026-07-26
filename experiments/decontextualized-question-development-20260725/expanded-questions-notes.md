# Expanded agent-development questions

`expanded-questions.jsonl` contains 75 hand-authored, source-grounded
questions for agent-tool development only: exactly 25 each of
`pronoun_resolution`, `attribution`, and `decontextualization`. Every record
uses `corpus_split: "agent-dev-20260725"` so a runner can import/select it
without colliding with the earlier retrieval-only development questions.

The set uses five CNN MediaSum test documents that are present in both local
comparison projections: `fixed/pilot-100x100-v1` and
`facts/pilot-named-v3`:

| Document | Questions per stratum |
| --- | ---: |
| `CNN-100102` | 5 |
| `CNN-100178` | 5 |
| `CNN-100299` | 5 |
| `CNN-100343` | 5 |
| `CNN-100741` | 5 |

The prior 15-question probe uses only `CNN-10004`, `CNN-100082`, and
`CNN-100095`; this file neither reuses those questions nor those documents.
Each item has an exact canonical document and turn span, and the evidence was
checked against canonical `turns` in Postgres on 2026-07-25.

This is inspected, authored **development data**, not a held-out paper set.
It is suitable for repeated harness reliability checks and prompt iteration,
but must not be reported as a confirmatory experiment or merged into a future
held-out benchmark.
