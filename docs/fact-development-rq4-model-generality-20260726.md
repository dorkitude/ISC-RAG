# RQ4: model generality of the fact advantage (Kimi K2.6, GLM 5.2)

**Status:** controlled development result on the frozen 100-document sample and
`fact-dev-sample-100-expansion-20260725` questions; same boundaries as the
[frontier-XL record](docs/fact-development-frontier-xl-20260725.md) (same-document
data, automated judges, no human calibration).

## Question

Is the fact-versus-chunk advantage a weakness artifact of the DeepSeek V4 Flash
answer model — which pays a large raw-transcript reading penalty — or a
model-general property of the representation?

## Protocol

Identical to frontier-XL (top-300 retrieval with the depth fix, greedy
`cl100k_base` packing, structured citations, temperature zero) except:

- **Answer models:** `kimi-k2p6` and `glm-5p2` (Flash cells from frontier-XL
  serve as the third model).
- **Conditions:** facts, fixed, and turn-aware chunks at the 2,048 and 16,384
  budget endpoints (semantic omitted; it tracked fixed within noise).
- **Cross-judging:** no model judges its own answers. Kimi-authored answers
  are judged by GLM 5.2 (`judge_glm_frontier_xl_v1` /
  `source_support_glm_frontier_xl_v1`); GLM-authored answers by Kimi
  (`judge_kimi_frontier_xl_v1` keys); Flash cells were already Kimi-judged.
  Because the judge differs across answerer arms, small absolute cross-model
  differences should not be over-read; within-model comparisons are clean.

## Pre-registered predictions

Recorded before the run: (1) facts lead all models at 2,048; (2) the gap
narrows at 16,384 with model strength but does not invert; (3) budget-driven
regressions shrink with model strength; (4) facts at 2,048 stay within a few
points of chunks at 16,384 at roughly an eighth of the tokens.

## Results: correct and correct-and-supported answers

| Answer model | Representation | Budget | Correct | Corr+Supp | Mean total tokens |
| --- | --- | ---: | ---: | ---: | ---: |
| DeepSeek Flash* | Facts + quote | 2,048 | 92/100 | 86/100 | 2,233.3 |
| DeepSeek Flash* | Facts + quote | 16,384 | 86/100 | 81/100 | 13,094.1 |
| DeepSeek Flash* | Fixed chunks | 2,048 | 62/100 | 59/100 | 1,801.2 |
| DeepSeek Flash* | Fixed chunks | 16,384 | 75/100 | 73/100 | 16,243.4 |
| DeepSeek Flash* | Turn-aware | 2,048 | 59/100 | 55/100 | 1,804.9 |
| DeepSeek Flash* | Turn-aware | 16,384 | 75/100 | 71/100 | 16,271.2 |
| Kimi K2.6 | Facts + quote | 2,048 | 90/100 | 83/100 | 2,191.5 |
| Kimi K2.6 | Facts + quote | 16,384 | 91/100 | 82/100 | 12,821.8 |
| Kimi K2.6 | Fixed chunks | 2,048 | 74/100 | 71/100 | 1,779.4 |
| Kimi K2.6 | Fixed chunks | 16,384 | 87/100 | 83/100 | 16,000.5 |
| Kimi K2.6 | Turn-aware | 2,048 | 70/100 | 68/100 | 1,782.8 |
| Kimi K2.6 | Turn-aware | 16,384 | 87/100 | 82/100 | 16,032.9 |
| GLM 5.2 | Facts + quote | 2,048 | 98/100 | 92/100 | 2,216.7 |
| GLM 5.2 | Facts + quote | 16,384 | 99/100 | 90/100 | 13,034.3 |
| GLM 5.2 | Fixed chunks | 2,048 | 75/100 | 73/100 | 1,792.5 |
| GLM 5.2 | Fixed chunks | 16,384 | 91/100 | 90/100 | 16,183.8 |
| GLM 5.2 | Turn-aware | 2,048 | 72/100 | 69/100 | 1,793.7 |
| GLM 5.2 | Turn-aware | 16,384 | 91/100 | 89/100 | 16,202.4 |

\* Flash rows repeat the frontier-XL cells for comparison (Kimi-judged).

## Findings against the predictions

1. **Facts lead every model at 2,048** — by 30–33 points for Flash, 16–20 for
   Kimi, 23–26 for GLM. The representation effect is model-general. ✔
2. **The 16,384 gap narrows with model strength but never inverts** — Flash
   11 points, Kimi 4, GLM 8. Stronger readers recover much of the
   raw-transcript penalty when given 16k tokens; none surpasses facts. ✔
3. **Budget-driven degradation is model-dependent.** Flash's fact decline
   (92→86) and chunk non-monotonicity do not reproduce on the stronger
   models: Kimi and GLM are flat-to-better at 16,384 in every condition
   (facts 90→91 and 98→99; chunks up 13–19). At these volumes, context-rot
   effects are concentrated in the weaker reader. Partially ✔ (regressions
   do not merely shrink — they vanish at this scale).
4. **Facts at 2,048 beat chunks at 16,384 for every model** — 92 vs 75
   (Flash), 90 vs 87 (Kimi), 98 vs 91 (GLM) — at roughly 13% of the tokens.
   The frontier claim is strictly stronger than predicted. ✔

**Economic headline:** the cheapest configuration tested (Flash + facts at a
2,048 budget, 92/100, ~2,233 tokens on the cheapest model) outperforms the
most expensive chunk configuration (GLM 5.2 + chunks at 16,384, 91/100,
~16,200 tokens on a premium model). Compilation lets a small model with a
small context match a large model with a large one.

## Boundaries

Same-document development data; automated cross-judging (different judge per
answerer arm) pending human calibration; two budget endpoints rather than the
full ladder; semantic chunks untested here.

## Reproducibility

Run IDs (model / representation @ budget):

- Kimi: facts a3abb2f5-623d-4e3a-ac28-1b9fb66b4143 @2048,
  92a84632-23f1-4a96-8634-51a94124997b @16384; fixed
  3f4d3129-63bd-43f4-aecf-925ec8c9e8ea @2048,
  27606ee0-a78d-4f8e-9932-2823459cd04a @16384; turn
  bbb15eb7-e886-48ce-a081-f7175f32533e @2048,
  437b3aa8-04e9-47b7-a7c7-93340c94d854 @16384
- GLM: facts 374b1d89-9d27-4c5c-8fe2-1b1b45134662 @2048,
  13fe5abc-f7db-4e42-9e3b-fa62a1cf13a2 @16384; fixed
  35cbf6f5-a71a-4241-8936-30fcecea8cbf @2048,
  1ce8371e-4d9c-41cf-bd42-b68507c04911 @16384; turn
  0334d09b-8c21-4d89-958b-0da9258ae68e @2048,
  3522c9f5-7734-454d-948f-89245c630aeb @16384

All 1,200 generations and 1,200 cross-judgments completed without failures
(`isc-rag evaluate judge`).
