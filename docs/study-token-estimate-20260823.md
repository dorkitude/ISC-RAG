# ISC-RAG study token estimate — 2026-08-23

## Headline

The best estimate for the experimental work behind the ISC-RAG results reported
in *RAG Deserves an Index* is **about 280 million model tokens**, with a
reasonable range of **260–310 million**.

This is an experiment-inference estimate. It does not include tokens spent in
Claude, Codex, or other assistants while designing the study, writing code, or
editing the paper. Those sessions do not have a complete common usage ledger.

## Accounting

| Component | Tokens | Basis |
|---|---:|---|
| RAG answer generation | 123.33M | Exact sum of provider-reported `answers.usage.total_tokens` across 28,967 recorded answer calls |
| Automated correctness and source-support judging | about 19.2M | 3.15M provider-reported; remaining 36,540 calls estimated from the persisted prompts and completions |
| Qwen contextual-retrieval reranking | about 75.6M | 1,647 executed query attempts at the measured approximately 45.9k reranker tokens per query, including the retained partial attempt |
| Main fact extraction | about 31.6M | The measured 20-document replay extrapolation of 26.3M for held-out 500, plus approximately 5.3M for the development 100 |
| Contextual-chunk generation | about 9.7M | Reconstructed from 1,768 persisted calls, full-transcript and chunk prompt sizes, and 50–100-token outputs |
| Corpus and query embeddings | about 10.6M | Tokenized-size estimate over 99,873 persisted embedded units and 79,783 recorded query embeddings |
| Rejection-rate replay | 1.05M | Exact provider-reported total in `experiments/rejection-rate-sample-20260803/observed-result.json` |
| Isolated Study #2 runs | 0.84M | Exact sum of persisted prompt and completion token columns |
| Question authoring and smaller prompt/extraction pilots | roughly 3–15M | Range estimate; these calls lack a unified provider-usage record |

The directly metered subtotal is **128.37M** tokens: answer generation,
metered judges, the rejection-rate replay, and Study #2. Reconstructable but
unmetered components bring the central estimate to approximately **280M**.

## Important limits

- Provider tokenizers differ, so character-reconstructed prompt totals are not
  exact billing totals.
- Several failed or superseded extraction attempts do not retain usage. The
  range allows for them, but the central estimate is not an invoice-level sum.
- The database contains development work as well as the frozen paper cells.
  This estimate answers “what did the study take,” not merely “what did the
  final table cells take.”
- Human reading, annotation, and writing effort are outside the token count.
