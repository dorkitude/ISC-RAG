# MediaSum resource

## Purpose

MediaSum is the first large-scale resource for ISC-RAG. It contains actual NPR
radio and CNN television interview transcripts, including utterance order and
speaker strings. It is a strong setting for testing whether source-linked facts
retrieve more efficiently and faithfully than text chunks.

## Source and version

- Dataset: [`nbroad/mediasum`](https://huggingface.co/datasets/nbroad/mediasum)
- Pinned revision: `bcbcef7de3b4e702e526352d825a4ff06de2becb`
- Official source repository: `../../sources/MediaSum/` (Git submodule at
  `7c92589bb6e69e678d043a9cf8c8c7285cd88055`)
- Retrieval date: 2026-07-23
- Files: `train.json`, `validation.json`, `test.json`
- Expected sizes: 4,260,826,970 B; 94,509,662 B; 96,548,679 B, respectively
- Expected SHA-256 values are recorded in `MANIFEST.sha256`.

## Scope

The source release describes 463.6K NPR/CNN media-interview transcript instances.
CNN entries can be topic-level broadcast segments rather than complete interviews.
Do not use provided summaries as answer evidence; they are auxiliary metadata.
Keep original program, date, URL, speaker, and utterance order during processing.

The official Git repository supplies the project materials and canonical data split,
but not the complete corpus. This local resource is the full, verified dataset
release associated with that pinned source.

## Reuse restriction

The source authors require research-only use. This resource is downloaded for
local research and must not be republished or redistributed. Consult the source
dataset card before any broader use.

## Acquisition

The Go/Cobra CLI downloads the pinned Hugging Face revision to a temporary file,
validates its checksum, then atomically places it in this directory. Run:

```sh
./bootstrap.sh --accept-mediasum-research-use
```

The flag confirms that the user has reviewed
[`MEDIASUM-NOTICE.md`](../../MEDIASUM-NOTICE.md) and will use the resource for
research purposes.

### Completed acquisition

- Completed: 2026-07-23
- Downloaded size: 4,451,885,311 B (4.45 GB)
- Verification result: `train.json`, `validation.json`, and `test.json` each match
  their pinned SHA-256 value in `MANIFEST.sha256`.
