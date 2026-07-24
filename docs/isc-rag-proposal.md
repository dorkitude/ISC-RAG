# ISC-RAG: Fact Compilation versus Chunking for Faithful Retrieval-Augmented Generation over Dialogue

**Kickoff proposal** · **Version:** v0 · **Date:** 2026-07-24  
**Authors and affiliations:** ❌ To be confirmed before external circulation. This proposal does not assign authorship, PI, corresponding author, or institutional affiliations.

## Researchers and affiliations

| Role | Named person / affiliation | Current responsibility | Status |
| --- | --- | --- | --- |
| Project lead / PI | ❌ TBD | Approve protocol, authorship, data handling, and go/no-go decisions | ❌ |
| Implementation owner | ❌ TBD | Freeze and operate the Go/PostgreSQL experiment pipeline | ❌ |
| Annotation/audit owner | ❌ TBD | Independently audit gold labels and fact provenance | ❌ |
| Analysis/writing owner | ❌ TBD | Statistical analysis, error analysis, and manuscript | ❌ |

**Authorship decision needed:** name the contributors, affiliations, roles, author order, and corresponding author. Authorship should follow substantive intellectual and research contributions, not merely access to infrastructure.

## Schedule overview (by date)

No publication deadline has been chosen, so all future dates remain deliberately `TBD` rather than invented.

| Date | Milestone | Type | Deliverable | Tasks (✅ done / ⚠️ in progress / ❌ not started) |
| --- | --- | --- | --- | --- |
| 2026-07-24 | Pilot infrastructure complete | Implementation | Reproducible local pilot | ✅ PostgreSQL/pgvector, MediaSum ingestion, four unit builders, extraction, embeddings, retrieval, and evaluation runner<br>✅ Candidate 100-question set imported and source-validated<br>⚠️ Freeze models and configuration (project lead: choices unapproved) |
| TBD | **Go/no-go 1: protocol freeze** | **Decision** | Registered internal experiment sheet | ❌ Create disjoint development set<br>❌ Independently audit question/evidence labels and fact sample<br>❌ Choose primary fact-context policy and all models<br>❌ Approve success criteria before full run |
| TBD | Full-scale indexed corpus run | Experiment | Four-condition retrieval results | ❌ Freeze corpus subset/split policy<br>❌ Run all four conditions with identical controls<br>❌ Export immutable run manifests and metrics |
| TBD | **Go/no-go 2: paper claim** | **Decision** | Results and error-analysis package | ❌ Compute paired uncertainty intervals<br>❌ Review validity and cost results<br>❌ Decide whether evidence supports a method paper claim |
| TBD | Manuscript and submission | Publication | ISC-RAG paper | ❌ Select venue and verify current call/deadline<br>❌ Draft, internal review, reproducibility release decision |

**Task totals:** ✅ 2 / ⚠️ 1 / ❌ 12. The nearest substantive gate is **Go/no-go 1**; it must be completed before a full-scale comparison is interpreted as evidence.

## Overview

ISC-RAG asks whether compiling dialogue into atomic, source-linked facts at ingestion improves retrieval-augmented generation (RAG) compared with conventional transcript chunking. The same MediaSum interview transcripts will be represented four ways—fixed windows, turn-aware chunks, deterministic lexical-cohesion chunks, and provenance-validated facts—while holding the embedding model, vector store, retriever, generator, prompt, and answer-context budget constant. The primary empirical output is paired, per-question evidence-retrieval performance; answer quality, provenance fidelity, cost, storage, and latency are supporting outcomes. The project has a runnable local pilot, but it has no comparative experimental result yet.

**Core claim to test:** For explicit claims in human interview dialogue, a source-grounded fact index may retrieve the supporting evidence more reliably and with less context than conventional chunking; this is a hypothesis, not a result.

**Keywords:** retrieval-augmented generation; dialogue; information extraction; provenance; vector retrieval  
**Target fields:** information retrieval; natural-language processing; machine learning systems

## Background

Dialogue is awkward RAG input: relevant claims may be surrounded by interruptions, speaker changes, questions, and unrelated discussion. Conventional systems choose an arbitrary chunk boundary and retrieve the resulting text block. ISC-RAG instead treats an explicit, standalone claim with an exact supporting quotation and source span as the retrieval unit.

The design is deliberately not a claim that facts are universally better. Extraction can omit information, distort qualification, add cost, or reduce usefulness for questions needing broader context. The study therefore compares representations on the same corpus and reports fact fidelity and operational costs alongside outcome metrics. The project’s factual-pipeline and storage choices are documented in [fact-pipeline-design.md](fact-pipeline-design.md) and [storage-design.md](storage-design.md). A source-verified [literature landscape](literature-landscape.md) identifies the closest known work and preserves explicit novelty uncertainties; this proposal makes no unsupported priority claim.

## Objectives, research questions, core contributions

**Objective:** establish a reproducible, source-grounded comparison of four RAG retrieval representations for real human interview transcripts.

**Research questions**

- **RQ1:** On an independently authored set of explicit, single-turn dialogue questions, does ISC fact retrieval improve gold-evidence recall over fixed, turn-aware, and semantic chunking under matched retrieval controls?
- **RQ2:** What are the answer-accuracy, evidence/citation, context-size, latency, storage, and indexing-cost trade-offs of fact retrieval versus chunking?
- **RQ3:** Are any observed gains robust to fact-extraction fidelity, audit results, and declared fact-retrieval variants (facts only versus parent-turn expansion; enrichment off versus on)?

**Hypotheses**

- **H1:** ISC facts increase evidence recall and/or ranking quality for explicit one-hop claims relative to each chunking baseline.
- **H2:** ISC facts reduce retrieved context size for a matched evidence-recall level.
- **H3:** Any benefit is conditional: provenance failures, extraction omissions, and questions requiring broader discourse can remove or reverse the advantage.

**Core contributions, if the study succeeds**

1. A controlled four-representation RAG protocol in which retrieval representation is the primary independent variable.
2. A source-grounded dialogue benchmark asset: the candidate 100-question MediaSum set with exact turn evidence, after independent audit and freeze.
3. A reproducible open-source implementation using PostgreSQL + pgvector that preserves fact-to-source provenance and records experiment artifacts.
4. An empirical account of quality, fidelity, and systems trade-offs—not only answer scores.

**Study type:** controlled empirical systems/method comparison (method proposal, currently at pilot-infrastructure stage).

## Assets and current status

| Asset | Status | Research use |
| --- | --- | --- |
| MediaSum release | ✅ Downloaded locally; checksum verified; research-use restricted | Primary corpus of NPR/CNN broadcast-interview transcripts |
| MediaSum source repository | ✅ Pinned Git submodule | Upstream code and corpus documentation; raw data is separately stored and ignored by Git |
| Canonical relational store | ✅ PostgreSQL + pgvector migrations and Docker stack | Documents, turns, units, facts, evidence, embeddings, questions, and run records |
| Four unit builders | ✅ Implemented | Fixed, turn-aware, semantic, and internal fact-input units |
| Fact pipeline | ✅ Implemented smoke path | Fireworks extraction; exact quote validation; optional topic/triple enrichment |
| Retrieval/RAG runner | ✅ Implemented smoke path | Embedding, HNSW indexing, retrieval evaluation, optional cited generation |
| Candidate gold test set | ⚠️ 100 one-hop questions across 25 NPR test transcripts | Deterministic answer/evidence evaluation; requires independent audit and checksum freeze |
| Full comparative study | ❌ Not run | No performance claim or effect size exists yet |

The published project record is [README.md](../README.md) and [PLAN.md](../PLAN.md). The 100-question set’s limits and normalization policy are recorded in [canonical-questions-v1.README.md](../examples/canonical-questions-v1.README.md).

## Experiment plan and methods

### Primary protocol

The corpus, document split, question set, embedding model, pgvector distance/index configuration, retrieval algorithm, `top_k` or fixed answer-context budget, answer model, and answer prompt must be identical for every primary condition. The only intended independent variable is the retrieval representation.

| Condition | Retrieval unit | Primary-status requirement |
| --- | --- | --- |
| Fixed | 4,000-byte transcript windows with 400-byte overlap in current v1 | Freeze exact parameters on dev only |
| Turn-aware | Speaker turns packed to the same maximum size | Freeze exact packing rule on dev only |
| Semantic | Deterministic lexical-cohesion chunks; current v1 uses 2,000-byte minimum | Freeze boundary and maximum-size rule on dev only |
| ISC facts | Atomic facts with exact source quote, offsets, speaker, and turn range | Primary form must be declared: facts only or fact-plus-parent-turn |

Fact extraction uses a structured, domain-neutral prompt with temperature zero in the current implementation. Facts without an exact located quote are discarded. Optional topic/triple enrichment cannot silently be part of the primary condition; it is an ablation. Exact duplicate removal is part of the planned primary fact condition; semantic clustering is an ablation.

### Evaluation data and split discipline

- Treat `canonical-questions-v1.jsonl` as **candidate held-out test data**. It contains 100 explicit single-turn questions, four per 25 NPR documents, with reference answers stated verbatim in the supporting turns.
- Build a separate, document-disjoint development set from a non-test MediaSum split before choosing chunk sizes, context budget, model identifiers, extraction settings, answer normalization, or ablation priority.
- Freeze corpus assignment, question-file checksum, prompts, model identifiers, CLI flags, package/container versions, and database migration revision before the full run.
- The MediaSum summary is metadata only, never answer evidence, retrieval evidence, or a gold label.

### Experiments

| ID | Experiment | Output | Status |
| --- | --- | --- | --- |
| E1 | Corpus ingest and unit construction on frozen corpus | Unit counts, sizes, and build time by condition | ⚠️ Smoke-tested only |
| E2 | Fact-fidelity audit | Quote-validation, source-entailment, omission, and duplicate estimates | ❌ |
| E3 | Primary retrieval comparison | Recall@k, MRR/nDCG, evidence coverage, context size | ❌ |
| E4 | Source-cited answer comparison | Strict normalized exact match, citation/evidence correctness, unsupported-answer rate | ❌ |
| E5 | Systems comparison | Extraction/indexing cost and latency, index size, retrieval latency, generation tokens/cost | ❌ |
| E6 | Predeclared ablations | Facts-only vs parent-turn expansion; enrichment off/on; optional near-dedup | ❌ |
| E7 | Transfer replication | Historical FOMC transcripts, if rights/format and labels permit | ❌; out of primary-paper critical path |

### Metrics and analysis plan

**Primary endpoint:** per-question gold-evidence recall@`k`, with `k` frozen before test evaluation. Because a retrieved unit may overlap the gold document/turn span, the unit must be scored by overlap with the stored evidence span rather than by text similarity to the reference answer.

**Secondary retrieval metrics:** MRR, nDCG, complete-evidence coverage when applicable, retrieved unit count, and context bytes/tokens. The present canonical set has one evidence turn per question, so complete-evidence coverage is initially equivalent to evidence hit; multi-evidence labels would be required for a distinct test.

**Answer metrics:** strict normalized exact match against the stored deterministic answer; a separately reported lenient numeric-normalization score only if predeclared; citation/evidence correctness; and unsupported-answer rate. Correct answers without retrieved gold evidence must not be counted as retrieval successes.

**Fact-fidelity metrics:** exact quote-validation rate, independently audited entailment/support rate, fact coverage/missed-major-claim rate, and exact/near-duplicate rate. A random, stratified audit sample and the annotation rubric must be selected before extraction at scale.

**Inference:** compare conditions on the same questions; report paired point differences with paired bootstrap confidence intervals. Report all four conditions, not merely the best chunk baseline. Inspect and categorize wins/losses by question type and transcript properties. The full sample-size/power rationale, bootstrap resample count, alpha/interval convention, and treatment of API failures remain ❌ decisions to document before the run.

## Validity plan

| Threat | Mitigation / decision |
| --- | --- |
| Tuning on the test set | Keep canonical v1 out of tuning; create a document-disjoint dev set and freeze configurations first. |
| Representation confounding | Match corpus, models, prompt, vector configuration, retrieval budget, and scoring across conditions. |
| Fact hallucination or lost qualification | Require exact quote provenance; audit support and omissions; report failures rather than filtering them from analysis. |
| Label bias or weak gold | Independent audit of all/stratified labels, adjudication procedure, frozen checksum, and transparent corrections/versioning. |
| Generator masking retrieval differences | Make retrieval evidence metrics primary; evaluate generated answers separately with citations. |
| Single corpus/question-type overclaim | Limit primary claim to explicit, one-hop MediaSum dialogue claims; label FOMC/multi-turn extensions as separate studies. |
| Provider/model drift | Record exact model identifier, prompt text/hash, date, temperature, parameters, failures, and costs; rerun frozen manifests where feasible. |
| Unequal compute | Report extraction/indexing cost, latency, index size, and retrieved context, not only quality. |

## Reproducibility, data, and ethics

The codebase uses Go/Cobra, PostgreSQL, pgvector, migrations, and a local Docker Compose service. The database stores canonical source turns, versioned source units, facts and their evidence, embeddings, questions, retrieval ranks, and optional answers. A run manifest should export commit SHA, submodule SHA, raw-data checksums, question-set checksum, migration revision, corpus selection query/list, CLI flags, prompts/hashes, model identifiers, and timestamp.

MediaSum is research-use restricted and the raw files must not be redistributed through this repository. Release code, migrations, schemas, scripts, non-sensitive configuration templates, provenance metadata allowed by the source terms, question IDs/evidence coordinates where permitted, and aggregated results. Confirm whether question text and reference answers can be redistributed under the MediaSum terms before publishing them; otherwise distribute a reproducible generation/import recipe or access-controlled artifact. Keep API keys only in ignored `.env` files. Do not claim that the study measures people, personal traits, or truth beyond what speakers explicitly stated in a transcript.

Ethical review status is ❌ unknown. The project lead should determine whether institutional review is required for analysis of already-published broadcast transcripts and whether any sensitive-content handling policy is necessary.

## Schedule and roles

| Deadline | Deliverable | Venue / output | Distinct contribution | Owner / role | Status |
| --- | --- | --- | --- | --- | --- |
| TBD | Protocol-freeze sheet | Internal | Locks test isolation, metrics, models, prompts, and audit rubric | Project lead + implementation + audit owners | ❌ |
| TBD | Audited benchmark v1 | Internal/research artifact | Evidence-labelled deterministic dialogue questions with change log | Annotation/audit owner | ❌ |
| TBD | Four-condition result package | Internal | Reproducible metrics, traces, costs, and error analysis | Implementation + analysis owners | ❌ |
| TBD | ISC-RAG manuscript | Selected venue | Controlled empirical claim about retrieval representations | All confirmed authors | ❌ |

**Dependencies:** E3 depends on protocol freeze, frozen corpus/unit construction, and embeddings for every condition. E4 depends on E3 plus frozen answer policy. Paper claims depend on E2–E5; a positive answer result without fidelity and cost evidence is insufficient. The named-owner gaps above are HIGH before kickoff.

**Preprint policy:** ❌ TBD. Decide venue compatibility, author approval, license, and versioning before public upload. Do not publish raw MediaSum data or credentials.

## Publication strategy options

No venue or deadline is selected. Verify all current calls, formatting, anonymity, data/code policies, and deadlines before committing.

| Option | Fit | Decision condition |
| --- | --- | --- |
| Information-retrieval / RAG conference track | Best if the final contribution is a controlled retrieval evaluation with clear systems trade-offs | Select after confirming a material, robust retrieval result and current venue requirements |
| NLP conference findings/short-paper track | Fits a focused dialogue-representation and provenance story | Select only if the evidence scope and paper length match a shorter empirical contribution |
| Journal extension | Fits a broader study with transfer corpus, richer annotation, and deeper ablations | Defer until the primary experiment and literature review establish enough independent scope |
| Systems/reproducibility artifact | Fits the open PostgreSQL/pgvector provenance implementation | Treat as companion output, not a substitute for empirical contribution |

**Anti-salami policy:** initial work should be one paper: representation comparison plus fidelity/systems analysis. A later FOMC transfer study or benchmark expansion is a separate paper only if it asks a genuinely distinct question, adds independent data/labels/methodology, and cites the first work rather than slicing the same results.

## Venues, materials, and links

| Deadline | Material | Current source / action |
| --- | --- | --- |
| TBD | Project roadmap | [PLAN.md](../PLAN.md) |
| TBD | Implementation and setup | [README.md](../README.md) |
| TBD | Fact pipeline contract | [fact-pipeline-design.md](fact-pipeline-design.md) |
| TBD | Storage/retrieval architecture | [storage-design.md](storage-design.md) |
| TBD | Candidate evaluation-set policy | [canonical-questions-v1.README.md](../examples/canonical-questions-v1.README.md) |
| TBD | Source-verified literature landscape | [literature-landscape.md](literature-landscape.md) |
| TBD | Target venue CfP/portal/template | ❌ Select venue, then record official URLs and constraints |

## Literature review

**Status: ✅ initial source-verified landscape complete.** See
[literature-landscape.md](literature-landscape.md) for 22 verified references
covering RAG retrieval units/chunking, fact or proposition indexing, dialogue
retrieval, provenance/faithfulness evaluation, and MediaSum. The proposal still
makes no priority claim: the landscape identifies a need for a broader systematic
IR/ACL/arXiv review before asserting that no exact four-way prior comparison exists.

## FAQ

**Q. What is ISC-RAG in plain language?**  
A. It tests two ways of filing an interview for an answering system: keep pieces of the conversation, or keep small claims with a pointer back to the exact words that support each claim.

**Q. Are facts automatically true?**  
A. No. A fact means “a small claim that a speaker explicitly said.” It must carry the quoted source words, and the study must audit whether the extracted wording faithfully matches them.

**Q. Why not just compare answer accuracy?**  
A. An answer can be guessed or produced from a misleading source. The main question is whether the system retrieved the actual supporting turn; answer quality is a separate check.

**Q. What is already proven?**  
A. Only that a small end-to-end pilot pipeline runs locally. No condition has yet been shown to outperform another.

**Q. What will a reviewer challenge first?**  
A. Likely test-set leakage/tuning, non-equivalent chunk conditions, label quality, omitted or distorted facts, provider-model variability, and overgeneralizing from short explicit interview claims. The protocol explicitly addresses these, but must execute and document the safeguards.

**Q. Can this paper claim better RAG for all dialogue?**  
A. No. The initial candidate set tests explicit one-turn claims in MediaSum. Multi-turn synthesis, speaker resolution, contentious claims, and other domains require separate evaluation.

## Kickoff actions and decisions

1. **HIGH — assign named owners and authorship roles.**
2. **HIGH — build a document-disjoint development set and forbid canonical v1 from tuning.**
3. **HIGH — independently audit canonical question/evidence labels and a stratified fact sample; define adjudication.**
4. **HIGH — freeze model identifiers, prompts, retrieval/context budget, primary fact-context policy, and statistical plan.**
5. **MEDIUM — choose the frozen MediaSum subset/stratification and full-run budget.**
6. **MEDIUM — perform source-verified literature review and select a venue only after reviewing results.**
7. **LOW — decide preprint, artifact release, and institutional-review policies before submission/public dissemination.**

## Update history

| Version (date) | Update |
| --- | --- |
| v0 (2026-07-24) | ✅ Captures implemented pilot assets and a controlled four-condition design.<br>❌ 12 scheduled tasks / ⚠️ 1 partially resolved item (baseline proposal).<br>**HIGH:** ownership, dev/test separation, independent audit, and configuration/statistical freeze are missing.<br>**Key decisions requested:** named roles/authorship; primary fact context; model/configuration; corpus scale; audit protocol; venue/ethics/preprint policies. |
