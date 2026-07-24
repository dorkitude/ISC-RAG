# ISC-RAG

ISC-RAG studies whether ingest-time extraction of atomic, source-linked facts
improves retrieval-augmented generation over conventional text chunking for real
human dialogue.

The authoritative project roadmap is [PLAN.md](PLAN.md).

## Status

- **2026-07-23:** Project created.
- **2026-07-23:** Selected MediaSum as the first large-scale resource. Downloaded
  the pinned three-file release (4.45 GB) and verified every SHA-256 checksum; see
  [MediaSum provenance](resources/mediasum/README.md).
- **2026-07-23:** Added the official MediaSum repository as a Git submodule, pinned
  at commit `7c92589bb6e69e678d043a9cf8c8c7285cd88055`.
- **2026-07-23:** Added local Fireworks credential configuration.
- **2026-07-23:** Replaced the shell downloader with the Go/Cobra `isc-rag`
  command-line interface.
- **2026-07-23:** Adopted a proposed source-grounded, two-stage fact-pipeline
  design adapted from the Endgame implementation; see
  [fact-pipeline design](docs/fact-pipeline-design.md).
- **2026-07-23:** Chose PostgreSQL with pgvector as the fully open-source
  canonical database and vector-search layer; see
  [storage design](docs/storage-design.md).
- **2026-07-23:** Implemented PostgreSQL migrations and MediaSum JSONL ingestion;
  verified with 10 documents and 309 source turns in the local database.
- **2026-07-23:** Implemented all four retrieval representations, Fireworks-backed
  fact extraction with exact quote validation, pgvector embeddings/search, and a
  question/evidence evaluation runner.
- **2026-07-24:** Added and database-validated a 100-question canonical MediaSum
  test set with exact turn evidence and deterministic reference answers; see
  [canonical-question notes](examples/canonical-questions-v1.README.md).

## Initial study design

Compare fixed-token, turn-aware, and semantic text chunks with ISC fact cards.
Keep embedding model, retriever, generator, and answer-context budget constant.
Each fact card must retain its source document, speaker, turn range, and verbatim
source span.

The proposed fact-pipeline contract, provenance validation, enrichment pass, and
ablation boundaries are documented in [fact-pipeline design](docs/fact-pipeline-design.md).
The canonical-data, vector-projection, and index-synchronization boundaries are
documented in [storage design](docs/storage-design.md).

## Resources

| Resource | Status | Role |
| --- | --- | --- |
| MediaSum | Downloaded and verified | Primary large-scale corpus of NPR/CNN broadcast-interview transcripts |
| Historical FOMC transcripts | Candidate | In-domain replication / continuity with ISC |
| NewsInterview | Candidate | Small CC-BY evaluation corpus with rich interview annotations |

## License and data handling

MediaSum is restricted to research use by its authors. Its data must not be
redistributed outside the terms of the source release. The raw files are ignored
by Git; this repository records acquisition metadata and processing code only.

## Upstream source

The official MediaSum repository is pinned as the
[`sources/MediaSum`](sources/MediaSum) submodule. It contains the source project,
paper split, and examples, but **not** the full corpus. The complete raw dataset
is the separately downloaded, checksum-verified local resource at
`resources/mediasum/`.

## Bootstrap

From this project directory, bootstrap through the Go/Cobra CLI with:

```sh
./bootstrap.sh
```

`bootstrap.sh` only verifies that Go is installed and delegates to
`isc-rag bootstrap`. The command loads the local `.env`, requires the Fireworks key,
then verifies existing MediaSum files or downloads the pinned files atomically
into `resources/mediasum/` with SHA-256 validation. MediaSum remains restricted
to research use under the source release terms.

## Local credentials

Store experiment credentials only in `.env`, which is intentionally ignored by
Git. `.env_example` is the tracked template:

```sh
cp .env_example .env
```

Set `FIREWORKS_API_KEY` in the local `.env` before
running experiments. Do not put real keys in the README, scripts, or commits.

## Local database

Start the project-local PostgreSQL + pgvector service, add the template `DATABASE_URL` to
your local `.env`, then apply migrations and ingest a small working sample:

```sh
docker compose up -d postgres
go run ./cmd/isc-rag migrate
go run ./cmd/isc-rag ingest mediasum --split test --limit 100
```

For a host-native installation, use `./bootstrap.sh --install-deps`. It supports
Homebrew on macOS and apt, dnf, or pacman on Linux; the Docker Compose stack is
the portable default.

The ingestion command reads MediaSum JSONL, upserts canonical documents, and
stores every ordered speaker turn with source-character offsets.

## Runnable experiment pipeline

The default `v1` setup uses 4,000-byte units, 400-byte fixed-window overlap, and
a 2,000-byte minimum for the deterministic lexical-cohesion semantic baseline.
These byte offsets match the stored canonical UTF-8 transcript offsets. Build the
three chunking conditions plus the internal 10,000-byte fact-input units:

```sh
go run ./cmd/isc-rag build --representation all --version v1 --split test --limit 100
```

Extract facts only after `fact-input/v1` is built. The extractor accepts only
JSON facts with a short exact quotation; a fact is discarded unless that quotation
is found in the canonical transcript. The `facts/v1` source units are the ISC-RAG
fourth condition.

```sh
go run ./cmd/isc-rag extract facts --version v1 --input-version v1 --split test --limit 100
go run ./cmd/isc-rag enrich --version v1
go run ./cmd/isc-rag embed --representation fixed --version v1 --split test
go run ./cmd/isc-rag embed --representation turn --version v1 --split test
go run ./cmd/isc-rag embed --representation semantic --version v1 --split test
go run ./cmd/isc-rag embed --representation facts --version v1 --split test
go run ./cmd/isc-rag index --embedding-model nomic-ai/nomic-embed-text-v1.5
go run ./cmd/isc-rag retrieve --representation facts --version v1 --split test --query 'What did the guest say about climate policy?'
```

`enrich` is optional metadata enrichment: it adds topics and subject–relation–target
triples but never changes fact text or evidence, so it is not part of the primary
fact-only condition unless explicitly enabled as an ablation.

Model names are command flags, so the experiment records a chosen model instead
of hiding a provider choice. The current account-visible default is Fireworks
GPT-OSS 120B for extraction/generation; Nomic Embed Text v1.5 remains the
embedding default. Confirm model availability on the account before a large run.
API calls can incur provider costs.

### Evaluation labels and run

Create independently authored questions as JSONL; the supplied MediaSum summary is
not used as an answer or evidence label. Each evidence item identifies a canonical
document external ID and inclusive turn range:

```json
{"external_id":"pilot-001","corpus_split":"dev","question_type":"single_fact","text":"What did the guest say?","reference_answer":"...","evidence":[{"document_external_id":"<MediaSum id>","turn_start":4,"turn_end":4}]}
```

Import labels and run a retrieval trial. The runner persists each ranked source
unit and reports gold-evidence recall@k; `--generate` additionally stores answers
that must cite `[S#]` retrieved source labels.

```sh
go run ./cmd/isc-rag evaluate questions pilot-questions.jsonl
go run ./cmd/isc-rag evaluate run --representation facts --version v1 --question-split dev --corpus-split test --top-k 5 --generate
```

Run the same fixed question set and flags once for `fixed`, `turn`, `semantic`,
and `facts`. Do not tune one representation on the held-out evaluation set.
`examples/smoke-questions.jsonl` is only a wiring test, not a study label set.
The actual canonical pilot set is
`examples/canonical-questions-v1.jsonl`: 100 one-hop questions across 25 MediaSum
test documents, each with one source turn and a verbatim deterministic answer.
Its selection and normalization policy are in
[canonical-question notes](examples/canonical-questions-v1.README.md).

## Project log

Update this README and [PLAN.md](PLAN.md) when a resource, experimental decision,
dataset freeze, or result changes.
