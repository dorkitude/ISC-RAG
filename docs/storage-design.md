# ISC-RAG storage design

## Decision

Use **PostgreSQL** as the system of record and **Turbopuffer** as the derived
vector-search layer. Do not use MongoDB or SQLite as the primary experiment store.

PostgreSQL stores the canonical corpus, fact provenance, pipeline versions,
annotations, experiment configurations, and results. Turbopuffer stores only
retrieval projections that can be regenerated from PostgreSQL.

## Why PostgreSQL

The experiment needs durable relationships and reproducible joins:

- a fact belongs to a document and one or more source turns/spans;
- an extraction run produces facts under a particular model, prompt, and config;
- a query has gold evidence, retrieval results, generated answers, and judgments;
- an index record must be traced to an exact canonical fact or chunk version.

These are relational invariants, not merely flexible documents. Use normalized
tables for the stable entities and `jsonb` only for changing provider responses
and auxiliary metadata.

SQLite is suitable for a throwaway local prototype but not as the primary store:
the intended pipeline has asynchronous ingestion, repeatable runs, and eventually
multiple workers. MongoDB would make the provenance and experiment-result joins
less explicit while offering no needed advantage over PostgreSQL `jsonb`.

## Canonical tables

| Table | Purpose |
| --- | --- |
| `documents` | Canonical MediaSum/interview record and immutable source metadata |
| `turns` | Ordered speaker turns with canonical text and character offsets |
| `source_units` | Versioned chunks or transcript units used by pipeline stages |
| `facts` | Atomic fact text, topic/triple enrichment, state, and content hash |
| `fact_evidence` | Verified quote, source offsets, and turn range for every fact |
| `pipeline_runs` | Extractor, prompt hash, parameters, status, timestamps, and costs |
| `index_outbox` | Durable create/update/delete events for Turbopuffer projections |
| `corpus_splits` | Frozen train/dev/test or scale-study assignments |
| `questions` / `gold_evidence` | Independently authored evaluation set and supporting spans |
| `experiment_runs` / `retrievals` / `answers` / `judgments` | Reproducible RAG evaluation artifacts |

## Turbopuffer projections

Use separate, versioned namespaces for each retrieval representation, for example:

```text
isc-rag-facts-v1
isc-rag-fixed-chunks-v1
isc-rag-turn-chunks-v1
isc-rag-semantic-chunks-v1
```

Each vector ID is the corresponding PostgreSQL UUID. Metadata includes at least
`document_id`, `source_unit_id`, `corpus_split`, `representation_version`,
`embedding_model`, and `pipeline_run_id`. Fact vectors also carry `fact_id`.

Never make a Turbopuffer namespace the only record of a fact, quote, or result.

## Write path

1. Persist a completed extraction result and its verified evidence in one
   PostgreSQL transaction.
2. In that same transaction, add an `index_outbox` event for each new, changed, or
   superseded retrieval unit.
3. A worker sends the event to the appropriate Turbopuffer namespace and marks the
   event complete only after success.
4. Re-extraction supersedes prior facts transactionally, then emits matching
   delete/upsert events. A reconciliation job can rebuild any namespace from the
   canonical database.

This avoids treating a successful vector write as proof that the canonical source
and provenance were persisted, or vice versa.

## Configuration

Add `DATABASE_URL` to local `.env` only when PostgreSQL provisioning begins. Keep
database migrations in the Go project and do not add `pgvector`; Turbopuffer is
the sole vector-retrieval system for this experiment.
