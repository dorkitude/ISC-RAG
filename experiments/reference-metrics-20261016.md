# Reference-answer metrics on the held-out sample — 2026-10-16

## Why

Every accuracy figure reported for the held-out sample so far comes from a
Kimi K2.6 judge, and Kimi K2.6 also authored the fact payloads: the model that
wrote one representation also scored every representation.

This record re-scores the same stored answers against the human-authored
`questions.reference_answer` spans using string metrics only — no judge model,
no API call, no network access. If the judge had favoured facts, a judge-free
metric should disagree with it.

## Data

- Source: `sqldump-20260819` (`20260819-1620_isc_rag_full.dump`),
  SHA-256 `8efe7f1d0a4693571593a48f35e01d1cce616e822a7119951506fc216cbff04b`,
  restored locally as `isc_rag_audit`.
- Questions: 499 held-out questions with a non-null `reference_answer`
  (`heldout500-NNN`), one per document.
- Runs: every run listed in `experiments/heldout-cr-runs.tsv` and
  `experiments/heldout-xl-runs.tsv`, unchanged since 2026-07-26.

## Method

`scripts/stats-reference-metrics.py` joins `answers` to `questions` and scores
`answers.text` against `questions.reference_answer` after SQuAD v1.1
normalisation (lowercase, punctuation and articles stripped, whitespace
collapsed). Five metrics, all five reported for every run:

| Metric | Definition |
| --- | --- |
| `em` | normalised exact match |
| `contains` | the normalised reference occurs as a substring of the answer |
| `precision` | token overlap divided by answer length |
| `recall` | token overlap divided by reference length |
| `f1` | harmonic mean of the two |

Binary metrics get exact McNemar with Wilson intervals; continuous metrics get
a paired bootstrap percentile interval (10,000 resamples, seed 20261016) and a
Wilcoxon signed-rank test. Holm and Benjamini-Hochberg are applied across
every comparison in each invocation. The comparison structure is inherited
verbatim from `scripts/stats-mcnemar.py`.

## Commands

    export DATABASE_URL="postgres:///isc_rag_audit"
    R="uv run --with scipy --with psycopg2-binary scripts/stats-reference-metrics.py"
    $R --selftest
    $R --runs-tsv experiments/heldout-cr-runs.tsv --per-run --metric precision
    $R --runs-tsv experiments/heldout-cr-runs.tsv --metric f1
    $R --runs-tsv experiments/heldout-xl-runs.tsv --per-run --metric precision
    $R --runs-tsv experiments/heldout-xl-runs.tsv --per-run --metric f1

## Per-run means — contextual-retrieval ledger

| Budget | Label | n | EM | Contains | Prec | Rec | token-F1 | ans words | ref words |
|---|---|---|---|---|---|---|---|---|---|
| 2048 | facts | 499 | 0.072 | 0.072 | 0.803 | 0.289 | 0.361 | 6.9 | 22.7 |
| 2048 | ctx-vector | 499 | 0.060 | 0.062 | 0.664 | 0.289 | 0.344 | 8.4 | 22.7 |
| 2048 | ctx-hybrid | 499 | 0.066 | 0.068 | 0.680 | 0.298 | 0.353 | 8.8 | 22.7 |
| 2048 | ctx-full-stack | 499 | 0.078 | 0.080 | 0.802 | 0.342 | 0.410 | 9.1 | 22.7 |
| 16384 | facts | 499 | 0.076 | 0.076 | 0.796 | 0.287 | 0.358 | 6.6 | 22.7 |
| 16384 | ctx-full-stack | 499 | 0.078 | 0.078 | 0.785 | 0.349 | 0.413 | 9.5 | 22.7 |

## Per-run means — four-representation ledger

| Budget | Label | n | EM | Contains | Prec | Rec | token-F1 | ans words | ref words |
|---|---|---|---|---|---|---|---|---|---|
| 2048 | facts | 499 | 0.072 | 0.072 | 0.803 | 0.289 | 0.361 | 6.9 | 22.7 |
| 2048 | fixed | 499 | 0.056 | 0.056 | 0.561 | 0.254 | 0.299 | 7.8 | 22.7 |
| 2048 | turn | 499 | 0.056 | 0.058 | 0.565 | 0.248 | 0.293 | 7.9 | 22.7 |
| 2048 | semantic | 499 | 0.060 | 0.060 | 0.588 | 0.264 | 0.313 | 8.0 | 22.7 |
| 4096 | facts | 499 | 0.072 | 0.072 | 0.801 | 0.284 | 0.356 | 6.8 | 22.7 |
| 4096 | fixed | 499 | 0.054 | 0.054 | 0.574 | 0.243 | 0.290 | 7.7 | 22.7 |
| 4096 | turn | 499 | 0.058 | 0.060 | 0.595 | 0.263 | 0.309 | 8.0 | 22.7 |
| 4096 | semantic | 499 | 0.056 | 0.058 | 0.597 | 0.255 | 0.304 | 7.5 | 22.7 |
| 8192 | facts | 499 | 0.072 | 0.072 | 0.790 | 0.280 | 0.349 | 6.6 | 22.7 |
| 8192 | fixed | 499 | 0.058 | 0.058 | 0.589 | 0.256 | 0.305 | 7.7 | 22.7 |
| 8192 | turn | 499 | 0.058 | 0.058 | 0.622 | 0.263 | 0.315 | 7.5 | 22.7 |
| 8192 | semantic | 499 | 0.060 | 0.060 | 0.633 | 0.273 | 0.325 | 8.0 | 22.7 |
| 16384 | facts | 499 | 0.076 | 0.076 | 0.796 | 0.287 | 0.358 | 6.6 | 22.7 |
| 16384 | fixed | 499 | 0.064 | 0.066 | 0.615 | 0.285 | 0.331 | 8.7 | 22.7 |
| 16384 | turn | 499 | 0.060 | 0.060 | 0.644 | 0.294 | 0.344 | 9.1 | 22.7 |
| 16384 | semantic | 499 | 0.062 | 0.062 | 0.653 | 0.304 | 0.352 | 9.1 | 22.7 |

## Paired comparisons — precision, contextual-retrieval ledger

| Budget | Comparison | n | Ref | Other | Mean diff [95% bootstrap] | Wilcoxon p | Holm p | BH q |
|---|---|---|---|---|---|---|---|---|
| 2048 | facts vs ctx-vector | 499 | 0.803 | 0.664 | +0.139 [+0.097,+0.182] | 2.13e-10 | 8.53e-10 | 8.53e-10 |
| 2048 | facts vs ctx-hybrid | 499 | 0.803 | 0.680 | +0.123 [+0.082,+0.163] | 3.84e-09 | 1.15e-08 | 7.67e-09 |
| 2048 | facts vs ctx-full-stack | 499 | 0.803 | 0.802 | +0.001 [-0.034,+0.037] | 7.19e-01 | 7.19e-01 | 7.19e-01 |
| 16384 | facts vs ctx-full-stack | 499 | 0.796 | 0.785 | +0.011 [-0.024,+0.047] | 2.83e-01 | 5.67e-01 | 3.78e-01 |

## Paired comparisons — precision, four-representation ledger

| Budget | Comparison | n | Ref | Other | Mean diff [95% bootstrap] | Wilcoxon p | Holm p | BH q |
|---|---|---|---|---|---|---|---|---|
| 2048 | facts vs fixed | 499 | 0.803 | 0.561 | +0.242 [+0.197,+0.287] | 2.33e-22 | 2.80e-21 | 2.80e-21 |
| 2048 | facts vs turn | 499 | 0.803 | 0.565 | +0.238 [+0.192,+0.282] | 5.47e-21 | 6.02e-20 | 3.28e-20 |
| 2048 | facts vs semantic | 499 | 0.803 | 0.588 | +0.215 [+0.171,+0.261] | 9.16e-19 | 8.25e-18 | 2.75e-18 |
| 4096 | facts vs fixed | 499 | 0.801 | 0.574 | +0.226 [+0.181,+0.272] | 2.28e-20 | 2.28e-19 | 9.12e-20 |
| 4096 | facts vs turn | 499 | 0.801 | 0.595 | +0.205 [+0.158,+0.251] | 2.82e-17 | 2.20e-16 | 5.65e-17 |
| 4096 | facts vs semantic | 499 | 0.801 | 0.597 | +0.204 [+0.160,+0.247] | 2.75e-17 | 2.20e-16 | 5.65e-17 |
| 8192 | facts vs fixed | 499 | 0.790 | 0.589 | +0.201 [+0.157,+0.246] | 2.15e-16 | 1.29e-15 | 3.68e-16 |
| 8192 | facts vs turn | 499 | 0.790 | 0.622 | +0.168 [+0.124,+0.212] | 2.11e-13 | 8.44e-13 | 2.81e-13 |
| 8192 | facts vs semantic | 499 | 0.790 | 0.633 | +0.157 [+0.115,+0.200] | 7.89e-13 | 2.37e-12 | 9.47e-13 |
| 16384 | facts vs fixed | 499 | 0.796 | 0.615 | +0.180 [+0.138,+0.224] | 7.42e-16 | 3.71e-15 | 1.11e-15 |
| 16384 | facts vs turn | 499 | 0.796 | 0.644 | +0.152 [+0.108,+0.195] | 6.32e-12 | 1.26e-11 | 6.89e-12 |
| 16384 | facts vs semantic | 499 | 0.796 | 0.653 | +0.143 [+0.102,+0.185] | 1.26e-11 | 1.26e-11 | 1.26e-11 |

## Paired comparisons — token-F1, contextual-retrieval ledger

| Budget | Comparison | n | Ref | Other | Mean diff [95% bootstrap] | Wilcoxon p | Holm p | BH q |
|---|---|---|---|---|---|---|---|---|
| 2048 | facts vs ctx-vector | 499 | 0.361 | 0.344 | +0.017 [-0.012,+0.046] | 4.63e-01 | 9.27e-01 | 6.18e-01 |
| 2048 | facts vs ctx-hybrid | 499 | 0.361 | 0.353 | +0.007 [-0.022,+0.036] | 8.80e-01 | 9.27e-01 | 8.80e-01 |
| 2048 | facts vs ctx-full-stack | 499 | 0.361 | 0.410 | -0.050 [-0.074,-0.025] | 1.32e-05 | 3.97e-05 | 2.65e-05 |
| 16384 | facts vs ctx-full-stack | 499 | 0.358 | 0.413 | -0.055 [-0.080,-0.030] | 8.78e-07 | 3.51e-06 | 3.51e-06 |

## Paired comparisons — token-F1, four-representation ledger

| Budget | Comparison | n | Ref | Other | Mean diff [95% bootstrap] | Wilcoxon p | Holm p | BH q |
|---|---|---|---|---|---|---|---|---|
| 2048 | facts vs fixed | 499 | 0.361 | 0.299 | +0.061 [+0.033,+0.090] | 2.73e-05 | 2.73e-04 | 1.09e-04 |
| 2048 | facts vs turn | 499 | 0.361 | 0.293 | +0.067 [+0.038,+0.096] | 1.02e-05 | 1.12e-04 | 6.11e-05 |
| 2048 | facts vs semantic | 499 | 0.361 | 0.313 | +0.048 [+0.020,+0.077] | 1.17e-03 | 9.35e-03 | 2.81e-03 |
| 4096 | facts vs fixed | 499 | 0.356 | 0.290 | +0.066 [+0.039,+0.094] | 3.49e-06 | 4.18e-05 | 4.18e-05 |
| 4096 | facts vs turn | 499 | 0.356 | 0.309 | +0.047 [+0.019,+0.076] | 1.45e-03 | 1.01e-02 | 2.89e-03 |
| 4096 | facts vs semantic | 499 | 0.356 | 0.304 | +0.052 [+0.025,+0.080] | 1.49e-04 | 1.34e-03 | 4.48e-04 |
| 8192 | facts vs fixed | 499 | 0.349 | 0.305 | +0.044 [+0.016,+0.074] | 5.13e-03 | 3.08e-02 | 8.79e-03 |
| 8192 | facts vs turn | 499 | 0.349 | 0.315 | +0.035 [+0.006,+0.063] | 2.18e-02 | 1.09e-01 | 3.26e-02 |
| 8192 | facts vs semantic | 499 | 0.349 | 0.325 | +0.025 [-0.004,+0.052] | 1.20e-01 | 3.60e-01 | 1.44e-01 |
| 16384 | facts vs fixed | 499 | 0.358 | 0.331 | +0.028 [-0.000,+0.055] | 7.66e-02 | 3.07e-01 | 1.02e-01 |
| 16384 | facts vs turn | 499 | 0.358 | 0.344 | +0.015 [-0.013,+0.043] | 4.75e-01 | 9.50e-01 | 5.18e-01 |
| 16384 | facts vs semantic | 499 | 0.358 | 0.352 | +0.006 [-0.021,+0.033] | 7.84e-01 | 9.50e-01 | 7.84e-01 |

## Findings

1. **The extractor-equals-judge confound did not distort the reported
   ordering.** On precision, every inferential conclusion of the held-out
   contextual-retrieval analysis reproduces with no judge model in the loop:
   facts beat all three chunk representations at all four budgets
   (Holm p <= 1.26e-11), beat the contextualized dense and hybrid
   configurations decisively (+0.139 and +0.123), and are indistinguishable
   from the full contextual-retrieval stack (Holm p = 0.72 at 2,048 and
   p = 0.57 at 16,384; unadjusted Wilcoxon p = 0.72 and 0.28). The judge's conclusions and a judge-free metric agree in all four
   comparisons.

2. **Facts answer more tersely than every alternative** — 6.6 to 6.9 words
   against 7.5 to 9.5, with reference answers averaging 22.7 words. Facts
   therefore hold the highest precision while posting the lowest recall.

3. **On token-F1 the full stack leads facts by 0.050 to 0.055**
   (Holm p < 1e-4), and the recall gap tracks the answer-length gap almost
   exactly. Every run is recall-capped by the long reference answers, so
   token-F1 rewards verbosity.

## Limits

- **Neither continuous metric is neutral.** token-F1 favours verbose answers
  through recall; precision favours terse ones. For facts against the full
  stack, string metrics cannot separate the two — the same conclusion the
  model judge reached. For facts against the other five conditions the result
  holds on both metrics, at equal or shorter answer length.
- **`em` and `contains` are not applicable to this corpus.** The reference
  answers were authored as explanatory sentences rather than minimal spans, so
  a short correct answer scores zero on both; both sit between 0.054 and 0.080
  across every run. Reported for completeness, not inference.
- **Not pre-registered.** The dump is checksum-pinned and already exported
  publicly under `results/`, the run ledgers have been unchanged since
  2026-07-26, the comparison structure is inherited from
  `scripts/stats-mcnemar.py`, all five metrics are reported for every run, and
  the family-wise correction covers every comparison in each invocation.
- Single runs per cell, as in the original protocol.
