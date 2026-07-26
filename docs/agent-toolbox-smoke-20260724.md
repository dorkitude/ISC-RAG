# Agent-toolbox 50-loop smoke check

**Status:** exploratory integration test, 2026-07-24. This is not a retrieval,
answer-quality, or model comparison result.

## Protocol

The fixed-chunk capability binary was the only executable exposed to each model.
For each question, the model had to return a JSON request for
`agent-toolbox-fixed-rag search`; the harness validated that request, executed
the binary, supplied the resulting JSON back to the model, and required a JSON
answer/citation. The model could request only 1--5 results and could not select
the representation, version, split, or embedding model.

The run used the same deterministic sample of 25 candidate canonical questions
for each model (seed `20260724`), for 50 model-question loops total:

| Model | Valid restricted calls | Provider/JSON failures | Expected chunk hit | Correct among hits (manual) | Citation among hits | Abstained when not hit |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| DeepSeek V4 Flash | 25 / 25 | 0 | 3 / 25 | 3 / 3 | 3 / 3 | 22 / 22 |
| Kimi K2.6 | 25 / 25 | 0 | 3 / 25 | 3 / 3 | 3 / 3 | 22 / 22 |

Both used Fireworks' no-reasoning option. Without it, model-visible planning
text sometimes violated the strict JSON contract. The raw outputs are ignored
under `artifacts/agent-toolbox-smoke/`.

## Interpretation limit

At the time of this run, the fixed `v1` index contained only 10 of the 10,000
MediaSum test documents. Only three documents represented in the 25-question
sample were indexed, so the remaining 22 questions could not test retrieval or
answer accuracy. Both agents appropriately returned insufficient-evidence
responses in those cases rather than answering from the wrong chunk.

The three retrieved answers were visibly supported by the retrieved chunk. The
simple exact-string screen undercounts them because reference answers contain
minor wording differences (for example, an omitted qualifier); it must not be
reported as a paper metric.

## Held-out discipline

The following candidate question IDs were exposed to models in this interface
test and are **consumed**: do not use them in a future untouched held-out paper
evaluation. Move them to development/integration use, and create replacement
held-out labels before freezing the primary question set.

```text
canonical-npr-759-03 canonical-npr-426-03 canonical-npr-771-01
canonical-npr-419-04 canonical-npr-1046-04 canonical-npr-2824-01
canonical-npr-640-04 canonical-npr-668-04 canonical-npr-2862-01
canonical-npr-1026-01 canonical-npr-759-01 canonical-npr-3378-01
canonical-npr-3499-04 canonical-npr-3528-01 canonical-npr-1683-02
canonical-npr-668-01 canonical-npr-3528-03 canonical-npr-419-03
canonical-npr-1683-04 canonical-npr-3378-02 canonical-npr-3499-03
canonical-npr-813-02 canonical-npr-771-04 canonical-npr-640-02
canonical-npr-3528-04
```

To reproduce the run after building the binary:

```sh
go run ./cmd/agent-toolbox-smoke \
  --model accounts/fireworks/models/deepseek-v4-flash \
  --questions examples/canonical-questions-v1.jsonl --limit 25 --seed 20260724

go run ./cmd/agent-toolbox-smoke \
  --model accounts/fireworks/models/kimi-k2p6 \
  --questions examples/canonical-questions-v1.jsonl --limit 25 --seed 20260724
```
