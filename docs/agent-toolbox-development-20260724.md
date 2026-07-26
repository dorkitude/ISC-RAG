# Fixed-chunk Q&A-agent development check

**Status:** exploratory development/integration study, 2026-07-24. This is not
a representation comparison or a paper result.

## Purpose and boundary

We needed to verify that Q&A agents can use the capability-limited fixed-chunk
toolbox for source-grounded questions, and get a first latency estimate. The
agent could make exactly one request to `agent-toolbox-fixed-rag search`; the
binary fixes representation (`fixed/v1`), MediaSum's test split, and the
embedding model. The model then answered only from that JSON response.

The local development set contains 25 newly written questions from the ten
documents actually present in the fixed `v1` test index (`NPR-12`, `NPR-134`,
`NPR-142`, `NPR-143`, `NPR-234`, `NPR-345`, `NPR-363`, `NPR-419`, `NPR-426`,
and `NPR-454`). Source-derived questions, reference answers, evidence turn
ranges, raw tool outputs, and manual labels are deliberately ignored under
`artifacts/agent-toolbox-development-20260724/`. They must not become a
held-out paper set.

For the bounded first run, the harness used the first ten questions (nine
distinctly scoreable; one was later found ambiguous) with no-reasoning mode and
a 15-second end-to-end deadline per question. Flash and Kimi each received the
same questions. Question-file SHA-256:
`5d5926519e9479b6a1e51a624ef02ca98f5f6198a2354cc90ce0166609d03de3`.

## Results

| Model | Valid constrained calls | Evidence-containing retrieval | Strict alias screen | Manual correctness | Citation coverage | Median / mean E2E latency |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| DeepSeek V4 Flash | 9 / 10 | 7 / 10 | 5 / 10 | 7 / 9 completed | 7 / 9 completed | 6.27 s / 6.61 s |
| Kimi K2.6 | 10 / 10 | 9 / 10 | 6 / 10 | 8 / 9 scoreable | 9 / 10 | 1.42 s / 1.53 s |

Flash had one provider request reach the 15-second deadline before a tool call;
this is counted as an interface failure, not an incorrect answer. Kimi made a
valid constrained request on all ten. The strict alias screen deliberately
under-counts semantically correct short forms (for example, `5` versus
`five`), hence the separate manual correctness column. Citation coverage means
the model's cited retrieved turn range contains the preannotated evidence turn;
it is not a claim that citation text has been independently adjudicated.

One question was excluded from manual comparative scoring after the fact: it
asked for two of three source-stated labor-market pressures. Flash selected the
prewritten pair and Kimi selected a different source-supported pair. This is a
question-authoring defect, not a model win. The remaining shared miss was a
three-name Giants defensive-line question: neither model retrieved the chunk
containing that evidence; Kimi supplied only two names and Flash abstained.

## Interpretation

The constrained-tool interface is usable by both models, and the questions are
substantively answerable when retrieval lands on the evidence. The bottleneck
in this tiny fixed-index condition is retrieval coverage, rather than the
answer formatting layer. Kimi was much faster in this small run; the one Flash
timeout and tiny sample make this an operational observation rather than a
model-performance conclusion.

Do not reuse these questions in a held-out evaluation. Repair the ambiguous
question, add a separate document-disjoint development set, and pre-register
manual grading before comparing chunk representations with facts.

## Reproduce

```sh
./scripts/build-agent-toolboxes.sh
shasum -a 256 artifacts/agent-toolbox-development-20260724/questions.jsonl

go run ./experiments/agent-toolbox-development-20260724 \
  --limit 10 --model accounts/fireworks/models/deepseek-v4-flash \
  --output artifacts/agent-toolbox-development-20260724/flash-results-v3.json

go run ./experiments/agent-toolbox-development-20260724 \
  --limit 10 --model accounts/fireworks/models/kimi-k2p6 \
  --output artifacts/agent-toolbox-development-20260724/kimi-results-v2.json
```

The `experiments/agent-toolbox-development-20260724` runner is intentionally
scoped to this development study and does not change an agent-toolbox capability
or a retrieval representation.
