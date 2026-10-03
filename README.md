# ISC-RAG

ISC-RAG studies whether ingest-time extraction of atomic, source-linked facts
improves retrieval-augmented generation over conventional text chunking for real
human dialogue.

## Papers

- [*RAG Deserves an Index: Why Ingest-Time Compilation Beats Query-Time Interpretation*](https://arxiv.org/abs/2608.20845) ([repository PDF](PAPER.PDF), [arXiv PDF](https://arxiv.org/pdf/2608.20845), [DOI](https://doi.org/10.48550/arXiv.2608.20845)) — the position paper reporting this repository's document-disjoint retrieval and QA comparison alongside the ISC maintenance result.
- [*Cost Scales with Change, Not Corpus Size: Incrementally Maintaining an Evolving Semantic Substrate*](https://arxiv.org/abs/2608.16621) ([PDF](https://arxiv.org/pdf/2608.16621), [DOI](https://doi.org/10.48550/arXiv.2608.16621)) — the companion incremental-maintenance study.

**Venue status:** *RAG Deserves an Index* is CIDR 2027 submission #150. It is
public on arXiv but has not yet been accepted; CIDR author notification is due
October 6, 2026.

The complete chronological record is in [CHANGELOG.md](CHANGELOG.md). The
reconstructed [study token estimate](docs/study-token-estimate-20260823.md)
separates directly metered usage from estimated components.

## Main result

The frozen held-out study used 500 document-disjoint MediaSum transcripts and
499 source-grounded questions. At the 2,048-token answer budget, compiled facts
reached 85.2% answer accuracy with roughly 2,200 prompt tokens per question,
compared with 72.5% using roughly 16,300 tokens for the strongest ordinary
chunk configuration. An Anthropic-style contextual-retrieval stack reached
comparable accuracy, but used about 21 times as many query-path tokens and did
not provide exact-quote-validated provenance.

These are controlled results over one dialogue corpus. The detailed conditions,
corrections, failed attempts, and statistical tests are documented in:

- [held-out protocol](docs/heldout-500-protocol-20260726.md)
- [held-out results](docs/heldout-500-results-20260726.md)
- [contextual-retrieval baseline](docs/contextual-retrieval-baseline-20260726.md)
- [paired McNemar tables](experiments/heldout-500-mcnemar-20260726.md)
- [raw-result audit](experiments/raw-data-audit-20260820.md)

## Architecture

The canonical store is PostgreSQL with pgvector. All retrieval conditions use
the same Nomic embedding model and answer/evaluation machinery while varying
the retrieval representation:

| Representation | Retrieval payload |
|---|---|
| Fixed chunks | Overlapping transcript windows |
| Turn-aware chunks | Consecutive dialogue turns grouped to a size target |
| Semantic chunks | Consecutive turns grouped by lexical cohesion |
| Compiled facts | Self-contained fact text plus exact supporting quotation and source coordinates |
| Contextual retrieval | Fixed chunks prepended with generated document context, optionally combined with BM25 and reranking |

The [fact-pipeline design](docs/fact-pipeline-design.md) defines extraction,
validation, provenance, and serving contracts. The [storage design](docs/storage-design.md)
defines canonical records, embeddings, and index synchronization.

## Repository layout

```text
cmd/                 Go command-line entry points
internal/            ingestion, extraction, embedding, retrieval, and evaluation
internal/storage/    PostgreSQL and pgvector migrations
experiments/         frozen run ledgers, aggregate statistics, protocols, and tools
examples/            deterministic sample membership and question fixtures
docs/                design records, experiment reports, and limitations
resources/mediasum/  acquisition manifest and provenance; no corpus files
results/              public, source-text-free result SQL
sources/MediaSum/     pinned official MediaSum repository submodule
```

## Bootstrap

Requirements: Go, Docker, and Fireworks access. On an exe.dev VM with the
`fireworks` integration attached, nothing else is needed: every Fireworks call
goes to `https://fireworks.int.exe.xyz/inference/v1`, which injects the key at
the network edge. Elsewhere, set
`FIREWORKS_API_BASE=https://api.fireworks.ai/inference/v1` and
`FIREWORKS_API_KEY` in the environment or the ignored `.env`.

```sh
git clone --recurse-submodules https://github.com/dorkitude/ISC-RAG.git
cd ISC-RAG
cp .env_example .env
# Off exe.dev only: add FIREWORKS_API_BASE and FIREWORKS_API_KEY to .env.
docker compose up -d postgres
./bootstrap.sh --accept-mediasum-research-use
go run ./cmd/isc-rag migrate
```

The bootstrap command requires explicit acknowledgement of the MediaSum
research-use notice, verifies Go, loads the ignored local `.env`, and verifies
or downloads the pinned files atomically with SHA-256 validation.

## Data handling

The code and original project documentation are released under the
[MIT License](LICENSE). MediaSum's authors request research-only use and
citation. The upstream repository does not provide an explicit license, so the
4.45 GB corpus is not relicensed or copied into this repository. Instead, the
pinned downloader obtains it from the upstream mirror and verifies every file.
See [MEDIASUM-NOTICE.md](MEDIASUM-NOTICE.md) before downloading or using it.

Source-derived research fixtures and calibration materials are retained so the
study can be inspected and rerun. They remain subject to the MediaSum
research-use request and the rights in the underlying NPR/CNN transcripts; the
repository's MIT license does not override those terms. Full database backups
are not distributed because they contain hundreds of complete transcripts.

The public SQL contains experiment configuration, counts, token totals, and
aggregate verdicts only. It intentionally excludes documents, turns, facts,
quotes, questions, reference answers, generated answers, citations, and model
prompts/completions. See [results/README.md](results/README.md).

## Reproduction boundaries

The frozen public ledgers identify the runs behind the reported comparisons.
Reproducing source-bearing rows requires independently obtaining MediaSum under
its research-use terms through
`./bootstrap.sh --accept-mediasum-research-use` and rerunning the pipeline.
Provider model availability
and pricing can change; model identifiers and the dates used by the study remain
recorded in the experiment documentation.

The authoritative project roadmap is [PLAN.md](PLAN.md).
