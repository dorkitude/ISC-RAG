# Held-out 500 replication — results (I6b) — 2026-07-26

Full replication of the budget-matched and frontier-XL studies on the frozen
held-out sample of 500 documents (`examples/heldout-500-20260726.txt`),
executed under the pre-registered protocol in
`docs/heldout-500-protocol-20260726.md`. 499 questions (one per document; one
document excluded after repeatedly failing the exact-quote validation gate).
Answer model DeepSeek V4 Flash at temperature zero; blind Kimi K2.6 judging
(keys `judge_kimi_heldout500_v1` / `source_support_kimi_heldout500_v1`).
Run ledgers: `experiments/heldout-matched-runs.tsv`,
`experiments/heldout-xl-runs.tsv`. McNemar tables:
`experiments/heldout-500-mcnemar-20260726.md`.

## Headline

**The development-set effect replicates on held-out data.** Facts win all 32
cells. All 24 facts-vs-chunk McNemar comparisons survive Holm FWER correction
at p ≤ 4.9e-06 (most below 1e-10; strongest 1.4e-42). At n=499 the paired
tests have roughly 5× the discordant pairs of the development runs.

## Retrieval (evidence recall)

| Representation | recall@20 (matched) | recall@300 (XL) |
|---|---|---|
| facts | 0.932 (465/499) | 0.986 (492/499) |
| fixed | 0.902 (450/499) | 0.984 (491/499) |
| turn | 0.892 (445/499) | 0.982 (490/499) |
| semantic | 0.890 (444/499) | 0.974 (486/499) |

Same ordering as development; facts lead at practical depth, everyone is
near ceiling at top-300.

## Answer accuracy (Kimi-judged correct / 499)

Budget-matched (top-20):

| Budget | facts | fixed | turn | semantic |
|---|---|---|---|---|
| 256 | **347 (69.5%)** | 126 (25.3%) | 135 (27.1%) | 140 (28.1%) |
| 512 | **371 (74.3%)** | 225 (45.1%) | 231 (46.3%) | 233 (46.7%) |
| 1024 | **388 (77.8%)** | 290 (58.1%) | 296 (59.3%) | 294 (58.9%) |
| 2048 | **393 (78.8%)** | 308 (61.7%) | 313 (62.7%) | 323 (64.7%) |

Frontier-XL (top-300):

| Budget | facts | fixed | turn | semantic |
|---|---|---|---|---|
| 2048 | **425 (85.2%)** | 310 (62.1%) | 313 (62.7%) | 323 (64.7%) |
| 4096 | **421 (84.4%)** | 316 (63.3%) | 332 (66.5%) | 326 (65.3%) |
| 8192 | **415 (83.2%)** | 320 (64.1%) | 340 (68.1%) | 344 (68.9%) |
| 16384 | **417 (83.6%)** | 344 (68.9%) | 357 (71.5%) | 362 (72.5%) |

Facts at 256 tokens (69.5%) beat every chunk policy at 2,048 tokens
(61.7–64.7%). Facts at 2,048/top-300 (85.2%) beat every chunk cell at any
budget, including chunks at 16,384 (68.9–72.5%).

Source-support rates track correctness closely for all representations
(facts XL: 414/409/400/404 supported; chunk cells within ±11 of their
correct counts), so the correctness gap is not an artifact of unsupported
guessing.

## Coverage vs reading-accuracy decomposition (XL cells)

Gold packed = any packed source overlaps a gold turn range; read-acc = blind
correctness among gold-packed questions only.

| Representation | Gold packed @2k/4k/8k/16k | Read-acc given gold @2k/4k/8k/16k |
|---|---|---|
| facts | 97.6% / 98.4% / 98.6% / 98.6% | 86.9% / 85.3% / 83.9% / 84.3% |
| fixed | 74.1% / 80.4% / 86.4% / 90.6% | 81.1% / 77.6% / 73.1% / 75.9% |
| turn | 72.3% / 80.8% / 86.4% / 89.4% | 84.5% / 81.1% / 78.4% / 79.8% |
| semantic | 74.9% / 82.0% / 87.2% / 90.4% | 85.3% / 78.2% / 78.9% / 79.6% |

Identical structure to development: facts pack gold evidence essentially
always from 2k upward, so their curve is flat; chunk budget gains are mostly
coverage gains, and chunk reading accuracy *declines* as contexts grow
(fixed drops 81%→73% from 2k to 8k) — the context-rot effect, now at n≈370–450
per cell.

## Measured read-path cost (mean prompt tokens per query)

| Cell | Mean prompt tokens | Accuracy |
|---|---|---|
| facts matched@256 | 374 | 69.5% |
| facts XL@2048 | 2,211 | 85.2% |
| best chunk cell anywhere (semantic XL@16384) | 16,283 | 72.5% |

Facts' best cell beats the best chunk cell by +12.7 points using 7.4× fewer
prompt tokens. The fact store exhausts at ~8.1k tokens (top-300 packs
everything), so facts@16384 reads only ~8.5k tokens — the representation
cannot even spend the largest budget, and doesn't need to.

## Honest deltas vs development

- Absolute accuracy shrinks for every representation: facts XL 92/87/88/86%
  (dev, n=100) → 85.2/84.4/83.2/83.6% (held-out, n=499); chunks shrink
  comparably (e.g. fixed@16384 75% → 68.9%).
- Effect direction and structure fully replicate; XL risk differences are
  +0.110 to +0.230 held-out vs +0.110 to +0.330 dev. Budget-matched
  differences at small budgets are *larger* held-out (+0.415 to +0.443 at
  256).
- Significance strengthens enormously with n=499: all 24 comparisons pass
  Holm; the weakest is p=4.9e-06 (dev's weakest was p=0.061 uncorrected).
- One harness fix mid-study (invalid-citation strip, see protocol correction
  log) — a robustness fix required to complete `semantic@512`; no completed
  run was affected and no scoring rule changed.

## Extension E1: contextual-retrieval baseline (held-out)

Pre-registered before execution (protocol Extension E1). All 1,485 held-out
fixed chunks were situated by DeepSeek V4 Flash with the frozen prompt,
embedded, and evaluated in four cells against the existing held-out facts
runs. Run ledger appended below; stats in
`experiments/heldout-500-mcnemar-20260726.md`.

**Correction (execution, not scoring):** the first E1 attempt ran retrieval
before the contextualized chunks had embeddings — all three affected runs
returned recall 0.000 and were deleted before any judging; the pipeline was
re-run with the embed step in place. No judged result was affected.

| Cell | Correct/499 | Evidence recall@300 | Mean prompt tokens |
|---|---|---|---|
| facts @2048 (comparator) | 425 (85.2%) | 0.986 | 2,211 |
| ctx-vector @2048 | 369 (73.9%) | 0.994 | 1,823 |
| ctx-hybrid @2048 | 377 (75.6%) | 0.994 | 1,825 |
| ctx-full-stack @2048 | 439 (88.0%) | 0.994 | 1,824 (+ ~45.9k rerank) |
| facts @16384 (comparator) | 417 (83.6%) | 0.986 | 8,479 |
| ctx-full-stack @16384 | 436 (87.4%) | 0.994 | 16,179 (+ ~45.9k rerank) |

McNemar (family of 4, Holm/BH):

- facts beat ctx-vector (+0.112, p=2.96e-06 Holm-corrected 1.18e-05) and
  ctx-hybrid (+0.096, Holm p=1.41e-04) decisively.
- facts vs ctx-full-stack is a **statistical tie at both budgets**
  (2048: 45/59 discordant, p=0.202; 16384: 42/61, p=0.076; neither survives
  any correction), with the full stack nominally ahead by 2.8–3.8 points —
  the development-set tie (8/4, p=0.39) replicates at n=499, now with the
  nominal direction favoring the full stack.
- Contextualization lifts retrieval above plain chunks (0.994 vs 0.984
  recall@300), consistent with development.

**Cost accounting.** The held-out contextualized chunks average ≈917
estimated tokens (dev: 916), so the reranking leg scores 50 candidates at
≈45.9k tokens per query, every query, before the reader sees anything. Full
query-path totals: ctx-full-stack ≈47.7k tokens/query versus facts ≈2.2k —
a ≈21× multiple for parity accuracy, unchanged from development. The fact
path uses no reranker. Ingest-side, contextualization is itself an
unvalidated partial compilation pass (a preamble per chunk with no
provenance gate); the fact pipeline's compile cost buys validated,
quote-linked units instead.

**Interpretation.** On held-out data the strongest chunk-side stack reaches
fact-level accuracy only by adding a per-query reranking stage that spends
~21× the tokens of the entire fact path, and by shipping ~16k-token contexts
at the top budget where the decomposition shows chunk reading accuracy
degrading. Facts deliver the same statistical accuracy tier with validated
provenance, ~2.2k tokens per query, and no query-time reranking dependency.
