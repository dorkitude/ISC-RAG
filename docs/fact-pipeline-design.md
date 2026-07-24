# ISC-RAG fact-pipeline design

## Purpose

ISC-RAG will compare conventional transcript chunks with a source-grounded fact
index. This design adapts the useful mechanics of the Endgame implementation to
real human interview transcripts; it does **not** use sales-specific categories,
claims, or downstream workflow behavior.

## Fact contract

Every persisted fact must contain:

```text
fact_id
document_id
fact_text                 # one standalone, atomic, explicit claim
supporting_quote           # exact source quotation
source_span_start/end      # character offsets in the canonical transcript
speaker_id / speaker_label
turn_start / turn_end
source_url, programme, date
topic                      # optional short retrieval aid
triple {subject, relation, target}  # optional structured metadata
extractor_version, model, prompt_hash
```

`fact_text` must not embed the speaker's name when speaker metadata can express
it. It must nevertheless be understandable without the originating dialogue
turn. A fact without verified source provenance is discarded.

## Proposed two-stage pipeline

1. **Prepare source units.** Canonicalize MediaSum records while retaining
   document identity, speaker order, programme, date, and source URL. Split only
   overlong transcripts into approximately 10,000-character units. When splitting,
   carry minimal preceding dialogue context and retain the exact canonical offsets.
2. **Extract raw facts.** Ask the extraction model for explicit, standalone,
   one-claim facts and supporting quotations. The extraction prompt must be domain
   neutral: it should focus on what interview participants explicitly state, not a
   predefined sales ontology.
3. **Validate provenance.** Locate every quote exactly in the canonical transcript;
   allow a documented, conservative normalization/fuzzy-match fallback. Attach
   source offsets, speaker, turn range, and source metadata. Drop unlocatable facts.
4. **Deduplicate.** Remove exact duplicate raw facts within each source document.
   Semantic clustering is an optional ablation, not an unreported preprocessing
   step, because it can itself affect retrieval recall and index size.
5. **Enrich facts.** In a separate batched pass, add an optional short topic and a
   subject–relation–target triple while preserving raw fact text verbatim. Retry
   only missing enrichment outputs. The enrichment model must not revise facts.
6. **Embed and index.** Embed the same fact text (and, if adopted, a declared
   topic prefix) into PostgreSQL with pgvector. Store all provenance and extraction
   metadata in the canonical relational records.

## Prompt constraints

The extractor should state that it:

- emits only explicit, source-supported information;
- preserves attribution in metadata;
- creates a fully qualified, standalone, atomic fact;
- provides a concise exact supporting quote; and
- does not speculate, summarize conversational dynamics, or produce meta-claims.

The enrichment pass should output a 1–3 word topic and an optional compact triple
with an ALL_CAPS relation. Its topic/tag vocabulary must be domain-neutral or
empirically derived from the corpus; the Endgame sales tags are out of scope.

## Experimental safeguards

- The primary comparison uses the same embedding model, corpus split, retriever,
  generator, and answer-context token budget for all four representations.
- The fact condition reports extraction coverage, quote-validation rate,
  source-entailment rate, and information-loss rate alongside RAG outcomes.
- Provided MediaSum summaries are never retrieval evidence or answer gold labels.
- Document-summary fact selection is excluded from the primary retrieval condition;
  it may be evaluated later as a clearly named ablation.
- Exact deduplication belongs to the main fact condition. Semantic clustering,
  topic-prefix embeddings, fact-to-parent expansion, and importance selection are
  separate ablations.

## Implementation decisions still open

- Fireworks models for raw extraction, enrichment, embedding, and answer generation.
- The neutral topic/tag policy, if any.
- Maximum quote length and fuzzy-match threshold.
- Whether primary fact retrieval returns facts only or expands each selected fact to
  its parent source turns.
- Batch size and concurrency for the research pipeline; Temporal-style asynchronous
  orchestration is not needed for the initial offline experiment.
