# Contextual-retrieval chunk baseline (frozen-100)

**Status:** controlled development result on the frozen 100-document sample
and the `fact-dev-sample-100-expansion-20260725` questions; same boundaries
as the frontier-XL record (same-document data, automated Kimi screens).

## Question

Anthropic's contextual retrieval is the strongest published chunk-side
remedy: prepend LLM-generated situating context to each chunk before
embedding, add a BM25 leg, and rerank. Does the fact representation's
advantage survive against that full stack at matched budgets?

## Protocol

- **Contextualization:** each `fixed/fact-dev-sample-100-fixed-v1` chunk was
  situated by DeepSeek V4 Flash reading the full transcript (prompt frozen in
  `internal/cli/contextualize.go`), producing
  `fixed/fact-dev-sample-100-ctx-v1` (context + chunk embedded together).
- **BM25 approximation:** Postgres full-text (`ts_rank_cd` over
  `websearch_to_tsquery`) fused with pgvector by reciprocal rank fusion
  (k=60), `--retrieval-mode hybrid`.
- **Reranking:** Fireworks `qwen3-reranker-8b` over 50 hybrid candidates
  (`--rerank-model`, `--rerank-candidates 50`); the facts bonus arm used 100
  vector candidates.
- Otherwise identical to frontier-XL: top-k 300 base retrieval, greedy
  `cl100k_base` packing to 2,048 / 16,384 tokens, DeepSeek V4 Flash
  answering with structured citations, blind Kimi screens under the
  `judge_kimi_frontier_xl_v1` keys.

## Results (correct / corr+supp of 100; gold-packed coverage)

| Condition | @2,048 | @16,384 | Gold packed @2k |
| --- | --- | --- | ---: |
| Plain fixed chunks (frontier-XL reference) | 62 (59) | 75 (73) | 80 |
| Ctx embeddings only | 76 (74) | 78 (77) | 86 |
| Ctx + hybrid FTS | 79 (77) | 80 (79) | 88 |
| Ctx + hybrid + rerank (full recipe) | 88 (86) | 87 (84) | 99 |
| Facts + quote (reference) | 92 (86) | 86 (81) | 100 |
| Facts + rerank (bonus) | 92 (85) | — | 100 |

## Findings

1. The full contextual stack closes roughly 70% of the chunk gap at 2,048
   (62 → 88) at essentially unchanged read cost; facts retain the best cell
   (92) and tie on correct-and-supported (86).
2. Mechanism decomposition: the situating preamble lifts *reading* accuracy
   (~76% → ~89% given packed gold) because it names referents — i.e., it
   performs partial, unvalidated compilation inside the payload. The
   reranker lifts *coverage* (88 → 99 gold-packed). Chunk-side remedies work
   insofar as they imitate compilation.
3. Reranking does not change the facts condition (92 → 92): fact retrieval
   is already saturated (100% gold coverage), so query-time crutches add
   nothing.
4. The ctx conditions barely improve from 2,048 to 16,384 (88 → 87): once
   retrieval is precise, additional raw budget adds nothing — corroborating
   the read-payload account.
5. Trade-offs the paper should state: the full recipe needs a reranking
   model in the query path (a per-query cost paid forever), its preamble is
   unvalidated generation (the HaluMem failure mode), and ingest-cost
   accounting differs (contextualization ≈ 1.2M tokens for this sample —
   see the itemized accounting below — versus the heavier fact-extraction
   bill).

### Significance and full-stack cost accounting (added after review request)

McNemar, facts vs the full CR stack at 2,048: discordant pairs **8/4**, exact
p = 0.39 — *not significant*. Facts vs the weaker CR variants remain
significant (hybrid: 17/4, p = 0.0072; embeddings-only: 21/5, p = 0.0025).
The honest claim is therefore: **facts match the strongest available chunk
stack on accuracy within statistical noise, at materially lower query cost,
with validated provenance.**

Query-path token accounting for the full CR stack (measured from stored
units and provider-reported usage): the reranking leg scores the top-50
contextualized candidates at a mean 916 tokens each ≈ **45.8k reranker
tokens per query**, plus a 1,848-token answer read — ≈ 47.7k tokens per
query versus 2,233 for the fact path: **≈21× the query-path tokens**. Caveat:
the reranker is an 8B model with a lower per-token price than the answer
model, so the dollar multiple is smaller than the token multiple but remains
large; the fact path needs no reranker at all (its bonus rerank arm changed
nothing). Ingest accounting: contextualization consumed ≈1.2M tokens for
this sample (0.91M transcript + 0.23M chunk + outputs); fact extraction's
ingest cost is higher, which is precisely the amortization trade the ISC
cost model prices.

*Ingest-figure reconciliation (2026-10-16).* Findings item 5 originally said
≈2.3M; that line predates this section and was not updated when it was added
(`54eaa9a`). Contextualization usage was not logged, so the bill was
reconstructed from `sqldump-20260819`: `internal/cli/contextualize.go` sends
the full speaker-labelled transcript (capped at 120,000 characters) plus one
chunk per call, for 283 `fact-dev-sample-100-ctx-v1` chunks over 100
documents. Under `o200k_base` that is ≈0.85M transcript + 0.19M chunk +
0.03M system-prompt input tokens and ≈0.02M output tokens, ≈1.09M in total.
This is consistent with ≈1.2M under DeepSeek V4 Flash's own tokenizer and
cannot reach 2.3M, which is retired.

Caveats: single runs per cell (judge noise plausibly ±3–4), rerank
candidate counts differ between chunk (50) and fact (100) arms for payload
size reasons.

## Reproducibility

Run IDs: ctx-vector ed4e1cc1-5a33-4f54-a91b-8c50bc45942b @2048,
9694cde6-94a4-4cd7-88f7-e6e9903cd8a0 @16384; ctx-hybrid
77b616ce-fc29-46db-9339-ac313885f878 @2048,
e58993e6-1f2e-4e8c-af89-d6139259e281 @16384; ctx-hybrid-rerank
95b8f4e3-3718-4c70-8f1c-3de34e9aab93 @2048,
453ddf0c-0605-48b3-a6b0-1f60674f3dcf @16384; facts-rerank
b2e31edc-f490-4eb2-b834-e080fba87dd5 @2048. All judged 700/700 without
failures.
