# ISC-RAG storage design

## Decision

Use **PostgreSQL with pgvector** as the fully open-source system of record and
vector-search layer. Do not use MongoDB, SQLite, or an external vector database.

PostgreSQL stores the canonical corpus, provenance, pipeline versions,
annotations, experiment configurations, results, and embeddings in one durable
database. This makes every retrieved vector traceable to an exact source unit and
its evidence.

## Why PostgreSQL + pgvector

The experiment needs durable relationships: a fact belongs to a document and
source turns, an extraction run has a model/prompt/configuration, and an answer
must join to retrievals and gold evidence. These are relational invariants.

Use normalized tables for stable entities and `jsonb` only for changing provider
responses and auxiliary metadata. pgvector removes a second operational data
store while keeping similarity search alongside those canonical relations.

SQLite remains useful for a throwaway single-process prototype, but not as the
primary store once batch ingestion and repeated experiment runs begin. MongoDB
would make provenance and experiment joins less explicit without a compensating
benefit.

## Canonical tables

| Table | Purpose |
| --- | --- |
| `documents` | Canonical interview record and immutable source metadata |
| `turns` | Ordered speaker turns with canonical text and character offsets |
| `source_units` | Versioned chunks or transcript units |
| `facts` | Atomic fact text, enrichment, state, and content hash |
| `fact_evidence` | Verified quote and source-span provenance |
| `pipeline_runs` | Model, prompt, configuration, cost, and status |
| `embeddings` | pgvector embedding for a source unit and declared model/dimensions |
| `corpus_splits` | Frozen train/dev/test or scale-study assignments |
| `questions`, `gold_evidence` | Evaluation questions and supporting spans |
| `experiment_runs`, `retrievals`, `answers`, `judgments` | Reproducible RAG artifacts |

## Retrieval representations

Store each representation in `source_units` with an explicit version, for example
`facts-v1`, `fixed-chunks-v1`, `turn-chunks-v1`, or `semantic-chunks-v1`.
Embeddings point at their source-unit UUID and record the embedding model and
dimension. The database is therefore both the source of truth and the index that
can be rebuilt through migrations and representation builders.

Create the pgvector HNSW/IVFFlat similarity index only after fixing the embedding
model and dimension; all vectors in a given index must be dimension-compatible.

## Local installation

The portable default is the project Docker Compose stack, based on the official
`pgvector/pgvector:pg17` image. Run `docker compose up -d postgres`, then apply
migrations with `go run ./cmd/isc-rag migrate`.

For host-native installation, `./bootstrap.sh --install-deps` invokes the correct
package-manager path for Homebrew on macOS and apt, dnf, or pacman on Linux. Host
installation is opt-in because it changes system packages.
