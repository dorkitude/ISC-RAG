# ISC-RAG

ISC-RAG studies whether ingest-time extraction of atomic, source-linked facts
improves retrieval-augmented generation over conventional text chunking for real
human dialogue.

## Status

- **2026-07-23:** Project created.
- **2026-07-23:** Selected MediaSum as the first large-scale resource. Downloaded
  the pinned three-file release (4.45 GB) and verified every SHA-256 checksum; see
  [MediaSum provenance](resources/mediasum/README.md).
- **2026-07-23:** Added the official MediaSum repository as a Git submodule, pinned
  at commit `7c92589bb6e69e678d043a9cf8c8c7285cd88055`.
- **2026-07-23:** Added local Fireworks and Turbopuffer credential configuration.
- **2026-07-23:** Replaced the shell downloader with the Go/Cobra `isc-rag`
  command-line interface.

## Initial study design

Compare fixed-token, turn-aware, and semantic text chunks with ISC fact cards.
Keep embedding model, retriever, generator, and answer-context budget constant.
Each fact card must retain its source document, speaker, turn range, and verbatim
source span.

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
`isc-rag bootstrap`. The command loads the local `.env`, requires both API keys,
then verifies existing MediaSum files or downloads the pinned files atomically
into `resources/mediasum/` with SHA-256 validation. MediaSum remains restricted
to research use under the source release terms.

## Local credentials

Store experiment credentials only in `.env`, which is intentionally ignored by
Git. `.env_example` is the tracked template:

```sh
cp .env_example .env
```

Set `FIREWORKS_API_KEY` and `TURBOPUFFER_API_KEY` in the local `.env` before
running experiments. Do not put real keys in the README, scripts, or commits.

## Project log

Update this README when a resource, experimental decision, dataset freeze, or
result changes.
