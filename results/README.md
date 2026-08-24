# Public result SQL

`public-results.sql` is a source-text-free export reconstructed from the frozen
local PostgreSQL study database on 2026-08-23.

It includes:

- retrieval/answer run configuration, status, row counts, and token totals;
- aggregate correctness and source-support verdict counts by run and judge key;
- extraction-pipeline run status and fact/unit counts;
- aggregate agent-evaluation, Study #2, and rephrasing-study run counts.

It excludes documents, turns, transcript text, facts, supporting quotations,
questions, reference answers, generated answers, citations, prompts, raw model
completions, embeddings, and human-calibration worksheets.

Restore into PostgreSQL with:

```sh
psql -d your_database -f results/public-results.sql
```

The dump creates only the `public_release` schema and does not depend on
pgvector.

SHA-256:
`9e0f49819a3f6c32da03277a3af051090dfdea4fcf47f6e38885fd9a8ee3a95e`.
