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
- **2026-07-24:** Added a source-verified literature landscape and a kickoff
  proposal; both distinguish the runnable pilot from unrun empirical claims.
- **2026-07-24:** Added an isolated DeepSeek V4 fact-rephrasing iteration CLI and
  completed a ten-iteration, fresh-evaluator prompt study on one MediaSum
  document. The account exposed V4 Pro rather than V4 Flash; the study is
  diagnostic only and did not promote a primary extraction configuration. See
  [the iteration report](docs/deepseek-fact-extraction-iterations.md).
- **2026-07-24:** Compared the read-only the production system call-transcript implementation
  with ISC-RAG and recorded evidence-backed, ranked transfer hypotheses. The
  memo distinguishes observed code behavior from proposed experiments; see
  [the production system fact-extraction hypotheses](docs/study2-fact-extraction-hypotheses.md).
- **2026-07-24:** Added a production-mechanics review of the production system. Its main
  conclusion is that provenance gating, speaker-aware context, domain metadata,
  source recovery, and evidence-first serving—not a magic standalone-fact
  prompt—likely explain production usefulness; see
  [why the production system works in production](docs/study2-production-why-it-works.md).
- **2026-07-24:** Added deterministic MediaSum speaker metadata. Corpus labels
  such as `ALEX COHEN, host` are parsed into normalized name and role fields
  while retaining the original label; no role is inferred from dialogue text.
- **2026-07-24:** Stream-audited all 463,596 local MediaSum records (13,919,244
  turns) through the ingestion schema and deterministic speaker parser. Every
  required source field passed; titles and programme names are often absent but
  remain optional document metadata.
- **2026-07-24:** Ran the isolated `fw-smoke-v1` Fireworks embedding smoke test
  over 50 MediaSum test documents: 174 fixed, 167 turn-aware, and 213 semantic
  units received 768-dimensional `nomic-ai/nomic-embed-text-v1.5` vectors.
  PostgreSQL persistence, model-specific HNSW creation, and a live Fireworks
  query retrieval all succeeded. Fact units were intentionally excluded pending
  a frozen fact configuration.
- **2026-07-24:** Completed `pilot-100x100-v1`: 100 test documents, all four
  retrieval representations, 10,000 deterministic source-phrase probes, and
  Fireworks Nomic embeddings. Turn-aware/fixed/semantic recall@5 was
  33.40%/33.21%/30.78%; the deliberately low-coverage fact baseline was 11.69%
  but stored 19× less retrieval text and returned an exact quote with every
  fact. This is a lexical plumbing pilot, not a paper result; see the local
  [results page](reports/pilot-100x100-v1.html).
- **2026-07-24:** Completed the controlled higher-coverage fact rerun
  `pilot-coverage-v2`: reducing fact-input units to 2,500 bytes increased the
  validated fact index from 629 to 2,353 facts and recall@5 from 11.69% to
  15.51% on the unchanged 10,000-probe set. Every stored fact retained an exact
  quote; the diagnostic lexical-pilot caveat remains.
- **2026-07-24:** Re-extracted and embedded the 100-document fact pilot as
  `pilot-named-v3`. Extraction now receives deterministic source speaker labels,
  forbids generic speaker terms in fact text, and mechanically resolves any
  residual generic attribution to the corpus-supplied name. The result is 2,405
  named fact vectors with zero remaining `the speaker` phrases.
- **2026-07-24:** Switched exploratory DeepSeek rephrasing defaults to Fireworks
  DeepSeek V4 Flash and added a persisted, source-auditable two-pass trial
  runner. It holds the isolated first-pass claims fixed and records each claim,
  quote, turn span, and Flash rewrite under ignored `artifacts/`.
- **2026-07-24:** Completed a five-document, 95-claim two-pass Flash check. The
  conservative contextual-rewrite prompt completed all claims without provider
  or JSON failures and outperformed its role-reference prompt variant on the
  mechanical screens. It remains an unpromoted engineering candidate pending an
  independent source-blind audit; see [the Flash study](docs/two-pass-flash-study.md).
- **2026-07-24:** Ran a separate four-run manual DeepSeek prompt diagnostic. Its
  result is that prompt length is not the leading failure hypothesis; see
  [manual prompt tests](docs/deepseek-manual-prompt-tests.md).
- **2026-07-24:** Started isolated, document-disjoint production-inspired Study #2
  with a fixed source-only three-document claim ledger, distinct context/proof
  concepts, and fresh blind audits. No candidate is yet promoted; see the
  [Study #2 report](docs/fact-study-2.md).
- **2026-07-24:** Ran three independent five-seed prompt studies, each followed
  by two controlled hill-climb refinements. None promoted a standalone-fact
  condition: role metadata improved attribution but not third-party
  coreference; dialogue units exposed the need to distinguish an utterance
  speaker from a quoted/asserted source; and dual-field prompts preserved exact
  evidence but did not make retrieval descriptions trustworthy facts. See
  [Study #4](docs/fact-prompt-study-4-role-metadata.md),
  [Study #5](docs/fact-prompt-study-5-dialogue-units.md), and
  [Study #6](docs/fact-prompt-study-6-dual-field.md).

## Initial study design

Compare fixed-token, turn-aware, and semantic text chunks with ISC fact cards.
Keep embedding model, retriever, generator, and answer-context budget constant.
Each fact card must retain its source document, speaker, turn range, and verbatim
source span.

The proposed fact-pipeline contract, provenance validation, enrichment pass, and
ablation boundaries are documented in [fact-pipeline design](docs/fact-pipeline-design.md).
The canonical-data, vector-projection, and index-synchronization boundaries are
documented in [storage design](docs/storage-design.md).
The [literature landscape](docs/literature-landscape.md) provides 25 verified
references and the current related-work positioning; the
[kickoff proposal](docs/isc-rag-proposal.md) records the evidence gates and open
project decisions.

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

Parse the corpus-provided speaker labels for already-ingested turns with:

```sh
go run ./cmd/isc-rag annotate speaker-metadata
```

This extracts only explicit label suffixes (for example, `host` or
`correspondent`) into structured fields; it does not infer roles from the
conversation.

Audit the full raw corpus without writing database records:

```sh
go run ./cmd/isc-rag audit mediasum --split all
```

The audit streams every JSONL record through the ingestion schema and the same
deterministic speaker parser. It checks the ID and aligned utterance/speaker
arrays required for canonical turns, non-empty source text and labels, UTF-8,
and reports optional document-metadata completeness separately.

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

### Isolated DeepSeek rephrasing study

`deepseek-study` is separate from the primary `extract facts` pipeline: it writes
raw provider output and exact-provenance facts to isolated study tables and local,
ignored artifacts. It defaults to
`accounts/fireworks/models/deepseek-v4-flash`; use `--model` to run a controlled
comparison with another available deployment. The full ten-iteration protocol and
results are in [the study report](docs/deepseek-fact-extraction-iterations.md).

`two-pass-test` is the narrower coreference experiment: it takes the frozen
`pilot-named-v3` isolated claims as pass one, supplies each claim's exact quote
and two labeled turns of local context, then asks Flash for a standalone rewrite.
It writes a complete, ignored JSON artifact (default
`artifacts/two-pass-rephrasing/flash-results.json`) so a sample can always be
traced back to its source rather than being selected from terminal output.

The [manual prompt diagnostic](docs/deepseek-manual-prompt-tests.md) found that
the earlier instruction files were only roughly 229--361 tokens beside an
approximately 8.3k-token labelled transcript input. The
[the production system comparison](docs/study2-fact-extraction-hypotheses.md) therefore ranks
speaker-aware context construction, a richer extraction contract, and staged
validation above prompt shortening as development hypotheses.

The active successor is a structural C4 replication in the isolated
production-inspired Study #2: a full-transcript antecedent-window generator is
followed by exact provenance validation, a transcript-blind quote-span support
audit, narrowly constrained repair, mechanical revalidation, and a fresh blind
re-audit. Its frozen protocol and eventual evaluations are recorded in
[the Study #2 report](docs/fact-study-2.md). It is not part of
the primary extraction pipeline unless it clears the three-document gate and
unchanged-config repeats.

The next isolated Study #2 condition uses deterministic interviewer-question /
answer source units. When a question supplies the scope needed to understand an
answer, the resulting fact is an attributed dialogue act with both exact Q and
A evidence stored; unpaired declarative turns remain their own source units.

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
