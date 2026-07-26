# Fresh-session agent-evaluation ledger

**Status:** development infrastructure. It is designed for the expanded
pronoun/attribution/decontextualization probe and is not itself a paper result.

`agent-toolbox-evaluate` runs one model behind exactly one compiled,
capability-limited search binary. Every question/replicate is a new API session:
there is no conversation history, cache, or previous answer presented to the
model. A session has two model calls: constrained JSON search selection, then a
fresh answer call receiving only the returned tool JSON. The host validates the
request and executes the binary; the model never has database or shell access.

## PostgreSQL records

Migration `012_agent_evaluation.sql` adds three audit tables:

| Table | One row per | Contents |
| --- | --- | --- |
| `agent_evaluation_runs` | model × capability evaluation run | Frozen system prompts, representation/version, embedding model, repetition count, parameters, lifecycle. |
| `agent_evaluation_attempts` | question × fresh replicate | Session UUID, raw selection and answer completions, parsed tool request/output, answer/citation, automatic screens, errors, and timing. |
| `agent_evaluation_retrievals` | retrieved rank | Attempt-local source-unit ID, rank, and similarity score. |

Question and canonical-evidence data remain in the existing `questions` and
`gold_evidence` tables. The runner imports its JSONL idempotently before a run.
It stores four automatic diagnostics, none of which substitutes for human answer
adjudication: valid tool JSON; retrieved-gold evidence coverage; cited-gold
turn coverage; and a deliberately conservative literal reference-answer match.

## Expanded development matrix

The planned command runs both Q&A models against both locally matched,
immutable capabilities:

| Capability | Representation/version | What the model receives |
| --- | --- | --- |
| `pilot-fixed-chunk-rag` | `fixed/pilot-100x100-v1` | Fixed chunk text only. |
| `named-fact-rag` | `facts/pilot-named-v3` | Rephrased fact plus its verified quote and speaker label. |

With 75 source-annotated development questions, ten independent repetitions,
two models, and two capabilities, this is **3,000 fresh sessions** (6,000 model
completions and 3,000 tool-side query embeddings). The question set is not
held out, so use these runs to debug stability, retrieval, and answer behavior;
do not promote their percentage to the paper.

## Run and inspect

```sh
./scripts/run-agent-development-20260725.sh
```

Each capability/model invocation creates a distinct run UUID and writes directly
to Postgres. Inspect the aggregate screens with:

```sql
SELECT r.model, r.capability, r.status,
       count(a.id) AS attempts,
       count(*) FILTER (WHERE a.tool_request_valid) AS valid_requests,
       count(*) FILTER (WHERE a.evidence_hit) AS evidence_hits,
       count(*) FILTER (WHERE a.citation_correct) AS citation_hits,
       count(*) FILTER (WHERE a.reference_answer_match) AS literal_reference_matches,
       round(avg(a.end_to_end_latency_ms)::numeric, 1) AS mean_end_to_end_ms
FROM agent_evaluation_runs r
JOIN agent_evaluation_attempts a ON a.run_id = r.id
WHERE r.question_split = 'agent-dev-20260725'
GROUP BY r.model, r.capability, r.status
ORDER BY r.model, r.capability;
```

To inspect instability by question, add `questions.external_id` and aggregate
over `agent_evaluation_attempts`; the `replicate` column is the fresh-session
index. Raw completions are intentionally in Postgres, rather than ignored local
artifacts, so a later human review can trace every aggregate back to tool output
and source-unit ranks.

### Invalidated configuration attempt

The first two fixed-chunk runs on 2026-07-25 are retained but marked `failed`:
they invoked the original `agent-toolbox-fixed-rag`, which is correctly frozen
to `fixed/v1`, while their ledger metadata incorrectly claimed the CNN pilot
version. Since that index only contains the unrelated ten-document NPR pilot,
all 1,500 fixed sessions had zero possible gold coverage. The dedicated
`agent-toolbox-pilot-fixed-rag` binary now freezes the actual
`fixed/pilot-100x100-v1` capability, and the two fixed runs must be repeated.
Never include invalidated run IDs in an aggregate.

## Expanded development run: 2026-07-25

The valid matrix is complete: 75 questions × 10 fresh sessions × two models ×
two capabilities. Aggregates below are direct PostgreSQL counts over the four
valid runs in `agent-dev-20260725`; `evidence` means the agent's chosen search
retrieved a source unit covering a gold turn, and `citation` means its cited
turn range covered gold evidence. They are **not** answer-correctness scores.

| Model | Capability | Sessions | Provider timeouts | Evidence | Citation | Mean completed E2E |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| Flash | `pilot-fixed-chunk-rag` | 750 | 1 | 320 (42.7%) | 201 (26.8%) | 2.59 s |
| Flash | `named-fact-rag` | 750 | 22 | 172 (22.9%) | 25 (3.3%) | 3.69 s |
| Kimi | `pilot-fixed-chunk-rag` | 750 | 0 | 314 (41.9%) | 305 (40.7%) | 1.44 s |
| Kimi | `named-fact-rag` | 750 | 0 | 186 (24.8%) | 127 (16.9%) | 1.67 s |

The broad development set does **not** reproduce the earlier 15-question
retrieval-only fact advantage. Fixed chunks led strongly on the new
decontextualization items (Flash 160 vs. 12 evidence hits; Kimi 164 vs. 37),
while attribution was close and facts had a small Kimi pronoun-retrieval edge
(50 vs. 39). This is useful negative development evidence: the first tiny
probe was not representative, and the fact pipeline/configuration needs work
before any primary hypothesis claim. The apparent fact deficit could reflect
fact coverage, the agent's one-search policy, fact retrieval granularity, or
the authoring distribution; it is not yet a controlled paper comparison.

The valid run IDs are retained in `agent_evaluation_runs`; the two initial
fixed-v1 runs and the one Docker-interrupted partial rerun are `failed` with
their invalidation reasons recorded in `parameters`. The raw traces make each
of these statements auditable and make future manual answer grading possible.

### Fact-plus-quote embedding ablation

`facts/pilot-named-v3-quote-v1` is an immutable copy of the named-fact source
units. Its displayed result is still the clean fact plus separate verified
quote, but its embedding input is `Fact: {fact}` followed by `Supporting quote:
{quote}`. All 2,405 facts were embedded with that input and the 1,500-session
agent matrix was rerun.

| Model | Fact only: evidence / citation | Fact + quote: evidence / citation | Fixed chunks: evidence / citation |
| --- | ---: | ---: | ---: |
| Flash | 172 / 25 | 190 / 30 | 320 / 201 |
| Kimi | 186 / 127 | 192 / 146 | 314 / 305 |

The quote improves both models modestly and eliminates the prior Flash fact
run's provider timeouts, but it does not close the retrieval gap. Retain it as
the stronger of the two current fact embedding ablations; do not promote it as
a winning experimental condition.
