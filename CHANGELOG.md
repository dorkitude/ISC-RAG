# ISC-RAG changelog

This is the chronological research and engineering record for ISC-RAG. It was
reconciled on 2026-08-23 against all 69 commits in the original private history,
the local PostgreSQL run ledger, the frozen experiment run-ID files, and the
paper-facing audit records. Dates below use the date recorded by Git or the
database; the database stores UTC, so a few late-Pacific runs appear on the next
calendar date.

The pre-publication Git history and private database backups are retained in a
verified private recovery bundle. The public repository retains its development
history and research fixtures; the corpus itself remains an upstream download.

## 2026-08-23 — papers, accounting, and public release

- Linked the two public papers from the repository:
  [arXiv:2608.20845](https://arxiv.org/abs/2608.20845) and
  [arXiv:2608.16621](https://arxiv.org/abs/2608.16621).
- Archived the exact arXiv v1 of *RAG Deserves an Index* as `PAPER.PDF`.
- Reconstructed total experimental model usage at approximately 280 million
  tokens, with a 260–310 million uncertainty range. The estimate distinguishes
  128.37M directly metered tokens from reranking, missing judge usage,
  contextualization, embeddings, and other reconstructed components.
- Moved the historical day log out of the README and reconciled it with the Git
  and database histories in this file.
- Replaced full private database releases with a source-text-free public SQL
  export containing run configuration and aggregate results only.
- Added a MediaSum research-use and citation notice. The official repository is
  pinned as a submodule and the 4.45 GB upstream release remains available
  through the checksum-verifying bootstrap path rather than being copied into
  Git.
- Kept source-derived research fixtures for reproducibility under the MediaSum
  research-use notice, while retaining full database snapshots privately.
- Produced and verified a private recovery bundle of the complete 69-commit
  history and tags before changing repository visibility.

## 2026-08-20 — independent audit and calibration pre-registration

- Pre-registered the I5c human-calibration acceptance criteria, including the
  sampling rule, allowed labels, agreement statistics, and decision gates.
- Independently restored the August 19 database snapshot and regenerated every
  paper-facing McNemar table. All 56 rows matched the committed tables exactly,
  and all 56 run IDs in the five frozen ledgers existed in the snapshot.
- Verified all 90 recorded machine calibration verdicts against `answers.usage`.
- Added explicit audit limitations: deleted unjudged runs cannot be recovered
  from the snapshot; development/held-out membership depends on the frozen
  sample files; and the validation-gate rejection rate comes from a replay, not
  the frozen database.
- Reconciled fact yield: 13,064 development facts, 69,746 held-out facts, and
  2,385 facts from other prompt-study documents, for 85,195 current fact rows.
- Corrected the 20-document replay yield comparison to a 3.4% difference from
  the corpus-wide mean (2.9% for the frozen run on those documents) and removed
  an unsupported implication that fact yield scales token cost directly.
- Reconstructed deterministic ingest cost from the saved replay:
  963,048 prompt and 90,187 completion tokens over 20 held-out documents.

## 2026-08-04 — qualitative gallery

- Added and curated a ten-document held-out compiled-claim gallery to illustrate
  the relationship between a fact payload and its exact supporting quotation.
- The gallery remains a research-use-only derived artifact under the MediaSum
  notice; it is not covered by the repository's MIT license.

## 2026-08-03 — extraction-gate measurement and recoverability

- Added a read-only replay harness for the frozen Kimi K2.6 extractor and ran it
  on the first 20 held-out documents in production order.
- Observed 2,724 candidates: 2,693 validated, 29 rejected (28 quotations not
  locatable in the target turn and one empty/overlong quotation), and two
  duplicates. The measured rejection rate was 1.06%.
- Measured 629 extractor calls, 963,048 prompt tokens, and 90,187 completion
  tokens: about 52.7k tokens and USD 0.064 per document under the recorded
  prices, extrapolating to roughly 26.3M tokens and USD 32 for 500 documents.
- Added migration 013 so future extraction runs persist every rejected
  candidate and reason, plus per-run validation-gate tallies.
- Established full PostgreSQL disaster-recovery snapshots. The two original
  842 MB snapshots were kept private and later removed from GitHub before the
  repository became public.

## 2026-07-27 — held-out contextual retrieval (E1)

- Completed the pre-registered Anthropic-style contextual-retrieval extension:
  generated situating preambles, vector and hybrid retrieval, and Qwen3
  reranking on the held-out 500-document sample.
- Facts decisively beat contextual vector and hybrid retrieval. The reranked
  full stack was statistically tied with facts: 88.0% versus 85.2% at 2,048
  answer tokens (`p=0.20`), while consuming approximately 45.9k reranker tokens
  per query versus about 2.2k total query-path tokens for facts.
- Disclosed and corrected an execution slip in which three retrieval runs were
  executed before contextual embeddings existed; they produced zero recall,
  were deleted unjudged, and were rerun only after embedding.
- Database ledger: eight retrieval/answer runs began on July 27 UTC, seven
  completed, one remained as a disclosed partial reranked attempt, and 3,842
  answer rows were recorded. The partial 16,384-token reranked attempt contains
  349 answers; the completed replacement contains all 499.
- Integrated the held-out replication and contextual-retrieval evidence into
  manuscript v5.

## 2026-07-26 — frozen development frontier and held-out replication

### Development frontier

- Ran the first answer-level cost pilot: facts at a 512-character retrieval cap
  reached 64/100 blind-judge accuracy versus 45/100 for fixed chunks at 2,048
  characters, using 50.4% fewer answer-model tokens.
- Completed the RQ4 model-generality sweep with DeepSeek V4 Flash, Kimi K2.6,
  and GLM 5.2 under cross-judging. Facts led each model at matched budgets;
  Flash with facts at 2,048 tokens reached 92/100, edging GLM with chunks at
  16,384 tokens (91/100) at roughly 13% of the tokens.
- Added contextual-retrieval machinery and ran the 100-document development
  baseline. The full stack reached 88/100 at 2,048 tokens versus 92/100 for
  facts; exact McNemar gave discordant counts 8/4 and `p=0.388`.
- Hardened statistical reporting with exact paired McNemar tests, Wilson
  intervals, odds ratios, risk differences, Holm family-wise correction, and
  Benjamini–Hochberg false-discovery correction.
- Built the blinded 80-item I5a worksheet, ten-item I5b practice set, and I5c
  Cohen-kappa tooling under a sampling/backfill rule committed before rater
  exposure. Source-bearing generated worksheets, protocols, keys, ratings, and
  analysis code remain available for research reproduction under the MediaSum
  notice.

### Held-out replication

- Froze a deterministic 500-document sample disjoint from the development 100
  and parameterized the transcript-only question-authoring script.
- Recorded two extraction false starts: the CLI default selected
  `gpt-oss-120b`, then a second run selected DeepSeek V4 Flash. Both were purged
  with zero downstream use. The protocol authority became
  `pipeline_runs.model`, and the final extraction used Kimi K2.6.
- Built all four representations, yielding 69,746 provenance-validated fact
  units and 499 database-validated, transcript-only Fable-authored questions.
- Completed 32 frozen held-out comparison cells. Facts won all 32; every one of
  the 24 facts-versus-chunk paired comparisons survived Holm correction at
  `p <= 4.9e-06`.
- At the headline 2,048-token cell, held-out fact accuracy was 85.2%, compared
  with 72.5% for the strongest ordinary chunk configuration. The decline from
  the 92% development result was recorded rather than hidden.
- Fixed a deterministic invalid-citation crash by stripping out-of-range and
  duplicate source IDs while retaining them in the answer audit payload. No
  scoring rule changed.
- Pre-registered held-out contextual retrieval (E1) before creating any E1
  artifact.

### Database reconciliation

- 149 retrieval/answer runs began July 26 UTC: 142 complete, six failed, one
  partial/running; they produced 25,125 stored answers.
- The date includes the budget frontier, top-k/depth runs, model-generality
  cells, development contextual retrieval, held-out budget-matched cells, and
  the beginning of the held-out XL frontier.
- Failed/partial records were retained in the database unless a protocol note
  explicitly records deletion before judging.

## 2026-07-25 — scaled development study

### Agent and retrieval diagnostics

- Added separately compiled, capability-limited read-only search tools and a
  fresh-session PostgreSQL ledger recording prompts, raw tool selection,
  retrievals, answers, citations, automated screens, and latency.
- Ran 14 agent-evaluation configurations containing 6,429 total attempts:
  11 completed runs and three failed/superseded fixed-index runs. The corrected
  four-condition expanded matrix contributed 3,000 completed attempts; earlier
  and source-grounded diagnostic runs remain in the ledger.
- The broad 75-question matrix reversed the earlier tiny retrieval probe:
  pilot fixed chunks beat the then-current named-fact index on coverage. A
  fact-plus-exact-quote embedding projection improved facts modestly but did not
  close the gap.
- Froze a deterministic 100-document development sample and commissioned one
  transcript-grounded question per document. Exact evidence coordinates were
  validated in PostgreSQL.

### Fact condition and question replication

- Promoted `turnwise-opinion-v1` as the bounded next fact condition and built a
  matched fixed-chunk control from exactly the same documents.
- Completed 13,064 Kimi K2.6 fact-plus-quote units for the development 100.
- Ran 240 initial fresh answer sessions across facts/fixed and Flash/Kimi. The
  fixed control led that small agent diagnostic, while direct retrieval over all
  100 questions favored facts 95/100 versus 87/100 at top-5.
- Authored a second transcript-only question set. Facts again led direct
  retrieval: 82/100 versus 64/100 at top-1 and 95/100 versus 88/100 at top-5.

### Accuracy, cost, and depth

- Completed the 4-representation x 4-budget answer frontier at
  256/512/1,024/2,048 tokens. Facts led every control at every budget; at 512
  tokens, facts reached 81/100 versus 45–46/100 for chunks at similar cost.
- Ran top-k pilots at `k=1/3/5/10/20/30`, then scaled them to all 100 questions.
  From `k=3` onward facts led while using much less answer context; at the
  reported point facts reached 68/100 at 366 mean tokens versus 53–60/100 at
  1,892–2,348 tokens for chunks.
- Completed the 1,600-answer frontier-XL study at 2,048–16,384 tokens. A spot
  check exposed pgvector's default `ef_search=40` as an accidental depth cap;
  the retrieval code was fixed and affected fact runs were regenerated.
- Facts at 2,048 tokens reached 92/100, beating every chunk condition through
  16,384 tokens. Fact accuracy declined at larger contexts, establishing a
  small-context optimum rather than a monotonic-context claim.
- Added paired statistics, coverage-versus-reading decomposition, manuscript
  drafts v1–v4, author affiliation/ORCID, and related-work positioning against
  Dense X, Contextual Retrieval, Token Tax, Mem0, and ingest-time compilation.
- Added the MIT license, pinned the official MediaSum repository as a submodule,
  removed internal-system references, and ignored local secrets/build outputs.

## 2026-07-24 — pilots and prompt studies

- Added a canonical 100-question MediaSum set with deterministic reference
  answers and exact evidence coordinates, plus a source-verified literature
  landscape and kickoff proposal.
- Added a ten-iteration DeepSeek fact-rephrasing study. The provider exposed V4
  Pro rather than the requested Flash model, so the result was retained only as
  a diagnostic.
- Audited all 463,596 MediaSum records and 13,919,244 turns through the parser;
  every required field passed. Speaker labels were normalized using corpus
  metadata without inferring roles from dialogue.
- Completed the `fw-smoke-v1` embedding path on 50 documents: 174 fixed, 167
  turn-aware, and 213 semantic units received 768-dimensional Nomic v1.5
  vectors; HNSW creation and live retrieval succeeded.
- Ran `pilot-100x100-v1` over 100 documents and 10,000 deterministic lexical
  probes. Recall@5 was 33.40% turn-aware, 33.21% fixed, 30.78% semantic, and
  11.69% for the intentionally low-coverage fact baseline.
- Reduced fact-input size to 2,500 bytes for `pilot-coverage-v2`, increasing the
  validated index from 629 to 2,353 facts and recall@5 to 15.51%.
- Rebuilt as `pilot-named-v3`: 2,405 named fact vectors with no remaining
  generic `the speaker` attribution.
- Added a two-pass Flash rephrasing runner and completed a five-document,
  95-claim check; it remained exploratory pending source-blind audit.
- Added capability-limited agent tools and smoke-tested Flash and Kimi: both
  completed the two easy in-index questions, followed by 50 abstention/interface
  trials and a ten-question source-grounded development check (Flash 7/9 and
  Kimi 8/9 scoreable answers).
- Began isolated Study #2 and three five-seed prompt families. Role metadata,
  dialogue units, and dual fields each exposed different attribution,
  coreference, or evidence-versus-description limitations; none was promoted.
- Database record: eight retrieval runs, one answer, ten `v1` rephrasing runs,
  four manual diagnostics, seven dual-field runs, and 83 Study #2 generations
  consuming 838,591 recorded tokens.

## 2026-07-23 — repository and pipeline foundation

- Created the project and selected MediaSum as the first large-scale corpus.
- Downloaded and checksum-verified the pinned three-file, 4.45 GB release and
  pinned the official source repository at commit
  `7c92589bb6e69e678d043a9cf8c8c7285cd88055`.
- Replaced the shell downloader with the Go/Cobra `isc-rag` CLI and added local,
  ignored Fireworks credential loading.
- Implemented PostgreSQL ingestion and verified it with ten documents and 309
  turns.
- Replaced the proposed external vector service with PostgreSQL + pgvector as
  the canonical open-source store and retrieval layer.
- Implemented fixed, turn-aware, semantic, and fact representations; Nomic
  embeddings; HNSW search; Fireworks-backed fact extraction with exact quote
  validation; and question/evidence evaluation.
- Adopted the source-grounded two-stage fact design and documented the storage,
  provenance, and index-synchronization contracts.

## Database ledger snapshot used for reconciliation

Snapshot state on 2026-08-23:

| Ledger | Rows | Notes |
|---|---:|---|
| `experiment_runs` | 174 | 166 complete, six failed, two partial/running across July 24–27 UTC |
| `answers` | 28,968 | 28,967 carry provider token usage; one smoke answer predates usage persistence |
| `pipeline_runs` | 30 | Includes successful runs, empty preflight failures, superseded model starts, and final extractions |
| `facts` | 85,195 | 13,064 development; 69,746 held-out; 2,385 other prompt-study documents |
| `source_units` | 101,532 | All representations and experiment versions |
| `embeddings` | 99,873 | Persisted source-unit vectors; query embeddings are not stored as rows |
| `agent_evaluation_runs` | 14 | 11 complete and three failed/superseded |
| `agent_evaluation_attempts` | 6,429 | Includes the corrected 3,000-attempt matrix and earlier diagnostics |
| `study_two_runs` | 83 | 838,591 directly recorded prompt and completion tokens |
| `fact_rephrasing_study_runs` | 21 | Ten iterative, four manual-diagnostic, and seven dual-field runs |

The public result SQL contains only source-text-free run metadata and aggregate
counts derived from this snapshot. The full database remains private because it
contains MediaSum transcripts, extracted quotations, and model artifacts.
