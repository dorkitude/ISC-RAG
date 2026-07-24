# ISC-RAG plan

## 1. Research objective

**Question:** Does ISC-style ingest-time extraction of atomic, source-linked
facts improve retrieval-augmented generation over conventional chunking for real
human interview transcripts?

**Working title:** *ISC-RAG: Fact Compilation versus Chunking for Faithful
Retrieval-Augmented Generation over Dialogue*.

The study is a method proposal (Type B) intended to reach P3 evidence: a
quantitative, controlled comparison with retrieval, answer-quality, and system
metrics. It extends ISC from fact compilation for compact representation to the
choice of RAG retrieval unit.

## 2. Corpus

### Primary corpus: MediaSum

MediaSum contains real NPR and CNN broadcast-interview transcripts. The local,
pinned release is checksum-verified at `resources/mediasum/`; raw data is
research-only and ignored by Git.

Preserve the original document ID, programme, date, source URL, speaker strings,
and turn order. CNN records may be topic-level broadcast segments rather than
complete interviews, so source grouping must be documented before evaluation.
Provided summaries are auxiliary metadata only: they are never retrieval evidence
or gold answer labels.

### Corpus strategy

1. Create a reproducible, stratified working subset for pipeline development.
2. Freeze a document-level train/dev/test split before tuning retrieval settings.
3. Scale the frozen pipeline to a larger MediaSum subset or full corpus.
4. Later add historical FOMC transcripts as an in-domain ISC continuity and
   transfer evaluation.

## 3. Retrieval conditions

Compare four representations of exactly the same source corpus:

| Condition | Retrieval unit |
| --- | --- |
| Fixed chunks | Fixed token-size windows with declared overlap |
| Turn-aware chunks | Whole speaker turns, split only when required by a maximum size |
| Semantic chunks | Deterministic lexical-cohesion boundaries under a declared size range |
| ISC facts | Atomic, source-grounded facts extracted at ingest time |

Keep the embedding model, vector-store configuration, retrieval algorithm,
generator, prompt, and answer-context token budget constant across the primary
conditions. This makes the independent variable the retrieval representation,
not a bundle of system changes.

## 4. Fact pipeline

The fact condition adapts useful structural ideas from the Endgame implementation,
but is domain-neutral and does not include sales ontology or workflow logic.

### Fact contract

Each persisted fact contains:

```text
fact_id, document_id, fact_text, supporting_quote,
source_span_start/end, speaker_id/label, turn_start/end,
source_url, programme, date, optional topic, optional triple,
extractor_version, model, prompt_hash
```

Facts are explicit, standalone, and atomic. Speaker identity belongs in
provenance metadata rather than being unnecessarily repeated in the fact text.
Any fact without validated provenance is discarded.

### Pipeline stages

1. Canonicalize documents and speaker turns; split only overlong transcripts into
   approximately 10,000-character units with minimal preceding context.
2. Extract raw facts and concise exact supporting quotes with a domain-neutral,
   structured-output prompt.
3. Locate each quote in the canonical transcript. Attach source offsets, speaker,
   turn range, and metadata; drop unlocatable facts.
4. Remove exact duplicate facts within a document.
5. In a separate batch pass, add optional 1–3 word topics and
   subject–relation–target triples without altering fact text.
6. Embed and publish fact projections for retrieval.

The extractor must not speculate, invent information, summarize conversational
dynamics, or emit meta-claims. The enrichment model must not rewrite raw facts.

### Fact ablations

These are not hidden preprocessing steps; report them separately:

- semantic clustering for near-duplicate fact reduction;
- topic-prefixed versus fact-only embeddings;
- facts-only retrieval versus fact-to-parent-turn expansion;
- document-level strategic/important-fact selection.

The first offline research implementation does not need Temporal-style
asynchronous orchestration.

## 5. Evaluation

Build an independently authored question set with gold supporting document/turn
spans. Stratify questions into:

1. single-fact lookup;
2. dialogue/coreference and attribution resolution; and
3. multi-fact synthesis.

### Retrieval metrics

- Gold evidence recall@k
- nDCG
- Complete-evidence coverage
- Retrieved context-token count

### Answer metrics

- Answer correctness
- Evidence/citation precision and recall
- Unsupported-claim rate

### Fact-pipeline metrics

- Quote-validation rate
- Source-entailment rate
- Information-loss / missed-major-claim rate
- Exact and semantic deduplication rate

### System metrics

- Extraction/indexing cost and latency
- Vector index size
- Retrieval latency
- Generation cost and token use

Use fixed seeds or repeated answer-generation runs, paired bootstrap confidence
intervals, and per-corpus reporting. A fact improvement is only credible when its
extraction fidelity and provenance coverage are also reported.

## 6. Storage architecture

Use PostgreSQL with pgvector as the canonical database and vector-search layer.
Do not use MongoDB, SQLite, or an external vector store as the primary experiment
store.

### PostgreSQL source-of-truth tables

| Table | Purpose |
| --- | --- |
| `documents` | Canonical transcript record and immutable source metadata |
| `turns` | Ordered speaker turns and canonical offsets |
| `source_units` | Versioned chunks or transcript units |
| `facts` | Atomic fact text, enrichment, state, and content hash |
| `fact_evidence` | Verified quote and source-span provenance |
| `pipeline_runs` | Model, prompt, configuration, cost, and status |
| `corpus_splits` | Frozen split and scale-study assignments |
| `questions`, `gold_evidence` | Evaluation questions and supporting spans |
| `experiment_runs`, `retrievals`, `answers`, `judgments` | Reproducible experiment artifacts |

Use normalized columns for stable entities and `jsonb` only for flexible provider
responses and auxiliary metadata.

### pgvector representations

Maintain independent versioned namespaces, such as:

```text
facts / source_units rows with representation_version values:
facts-v1, fixed-chunks-v1, turn-chunks-v1, semantic-chunks-v1
```

Embeddings are stored in PostgreSQL, keyed by source-unit UUID and embedding model.
Store document/source-unit IDs, corpus split, representation version, embedding
model, and pipeline run in relational columns. Never separate a fact from its
canonical provenance. The CLI creates one partial HNSW pgvector index per
embedding model and dimension after embeddings exist; it refuses mixed-dimension
model rows.

### Synchronization rule

Persist a canonical extraction result, its evidence, and embedding in one
PostgreSQL transaction. Re-extraction supersedes prior facts transactionally.
The database is the only index to reconcile.

## 7. Tooling and local setup

The project uses Go and Cobra for its CLI.

```sh
./bootstrap.sh
```

`bootstrap.sh` checks that Go is installed and invokes `isc-rag bootstrap`. The
CLI loads `.env`, requires `FIREWORKS_API_KEY`, then
verifies or downloads the pinned MediaSum release.

Local secrets go only in `.env`, which is ignored by Git. `.env_example` is the
tracked template. Add `DATABASE_URL` when PostgreSQL provisioning starts.

The official MediaSum source repository is pinned as `sources/MediaSum`; the raw
corpus is separately downloaded at `resources/mediasum/`.

## 8. Immediate implementation sequence

1. Provision PostgreSQL locally and add migrations for documents, turns, pipeline
   runs, facts, fact evidence, and pgvector embeddings.
2. Implement MediaSum canonical ingestion and frozen subset/split creation.
3. Implement the four representation builders, starting with fixed and turn-aware
   chunks as reproducibility baselines. **Implemented:** `build` has fixed,
   turn-aware, lexical-cohesion semantic, and fact-input builders.
4. Implement the two-stage fact pipeline and quote-provenance validator.
   **Implemented:** `extract facts` performs raw extraction and exact quote
   validation; `enrich` adds optional topic/triple metadata as a declared ablation.
5. Add pgvector embedding writes and representation-specific similarity indexes.
   **Implemented:** `embed`, `index`, and `retrieve` use the database only.
6. Create the question/evidence annotation format and annotate a pilot set.
   **Format and importer implemented; human pilot annotation remains.**
7. Run a pilot, inspect failure modes, freeze configuration, then execute the
   full comparison and ablations.

## 9. Open decisions

- Freeze Fireworks models for extraction, enrichment, embedding, and answer generation.
- MediaSum working-subset size and stratification policy.
- CNN grouping rule and maximum transcript/unit size.
- Neutral topic/tag policy and triple schema.
- Quote-length limit and conservative fuzzy-match threshold.
- Primary fact context: facts only or source-turn expansion.
- PostgreSQL provisioning method and migration library.
- Question-set size, annotation protocol, and judging model/human review plan.

## 10. Project state

- MediaSum is downloaded locally, pinned, and SHA-256 verified.
- The official MediaSum repository is a pinned Git submodule.
- The Go/Cobra bootstrap command is implemented, validates the Fireworks key,
  supports opt-in PostgreSQL/pgvector host installation, and ensures the resource
  is present.
- PostgreSQL migrations and MediaSum JSONL ingestion are implemented and verified
  locally with a 10-document / 309-turn sample.
- All primary representation builders are implemented and verified on a local
  ten-document / 309-turn MediaSum smoke sample (25 fixed, 24 turn-aware, 31
  semantic, and 13 fact-input units).
- The raw fact extractor performs exact quote-provenance validation before facts
  become retrieval units. Credentialed Fireworks extraction, enrichment,
  embedding, retrieval, and answer-generation smoke runs have succeeded; model
  selection still needs to be frozen before the pilot.
- pgvector embedding writes, model-specific HNSW indexing, single-query search,
  question/evidence import, retrieval evaluation, and optional answer persistence
  are implemented. Independently authored pilot questions and full experiment
  execution remain.
