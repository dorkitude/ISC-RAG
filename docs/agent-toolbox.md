# Agent toolbox

`agent-toolbox` is a family of small, capability-limited CLIs for Q&A agents in
the representation comparison. A binary is compiled for one retrieval condition
and exposes only one operation: `search`.

## Why separate binaries

Giving an evaluation agent a general experiment CLI creates an easy accidental
confound: it could select another representation, version, model, or corpus
split. Each binary instead compiles those choices into the executable:

| Binary | Fixed capability |
| --- | --- |
| `agent-toolbox-fixed-rag` | `fixed/v1`, MediaSum test split, Nomic v1.5 |
| `agent-toolbox-turn-rag` | `turn/v1`, MediaSum test split, Nomic v1.5 |
| `agent-toolbox-semantic-rag` | `semantic/v1`, MediaSum test split, Nomic v1.5 |
| `agent-toolbox-pilot-fixed-rag` | development-only `fixed/pilot-100x100-v1`, MediaSum test split, Nomic v1.5 |
| `agent-toolbox-named-fact-rag` | development-only `facts/pilot-named-v3`, MediaSum test split, Nomic v1.5 |
| `agent-toolbox-named-fact-quote-rag` | development-only `facts/pilot-named-v3-quote-v1`, MediaSum test split, Nomic v1.5 |
| `agent-toolbox-turnwise-opinion-fact-rag` | development-only `facts/turnwise-opinion-v1`, MediaSum test split, Nomic v1.5 |
| `agent-toolbox-fact-development-fixed-rag` | development-only `fixed/fact-dev-sample-100-fixed-v1`, MediaSum test split, Nomic v1.5 |

Each binary supports only:

```sh
agent-toolbox-fixed-rag search --query 'What profession does the guest’s father have?' --top-k 5
```

It writes JSON containing the immutable capability name, fixed configuration,
query, ranked chunk text, document ID, turn span, and similarity score. `top-k`
is bounded to 1--5; no flag can select another representation, version, split,
or embedding model.

## Build and use

```sh
./scripts/build-agent-toolboxes.sh
./bin/agent-toolbox-fixed-rag search --query 'What profession does the guest’s father have?'
```

The tools load the project's `.env` only to connect to PostgreSQL and embed the
query. They open PostgreSQL sessions with `default_transaction_read_only=on` and
perform no schema migration or write. This prevents writes from the executable,
but it is not a complete security boundary: give each external agent a separate,
least-privilege database/API credential and isolate its filesystem/network access
if adversarial containment is required.

The five non-`v1` capabilities are explicitly development-only and remain
separate from the three primary chunk-condition binaries. Their presence in the
build script does not promote any fact condition: use their frozen names only
for the documented diagnostics, and do not include their results in a primary
comparison until the fidelity and held-out-evaluation gates are passed.

## Kimi and Flash tool-use smoke check

**Status:** exploratory interface test, 2026-07-24; not a retrieval or answer
quality result. `agent-toolbox-smoke` presents a model with only the fixed-RAG
search contract, validates its JSON request, executes the compiled binary, then
asks for an answer and source citation from that returned JSON.

Two known in-index questions about the same `NPR-454` chunk were used. Under
the provider's no-reasoning option, both models produced a valid `top_k=1`
search request for both questions, retrieved `NPR-454`, answered with the
expected surface form, and cited the source range:

| Model | Tool requests | Retrieved expected chunk | Answer surface | Citation surface |
| --- | ---: | ---: | ---: | ---: |
| `accounts/fireworks/models/deepseek-v4-flash` | 2 / 2 | 2 / 2 | 2 / 2 | 2 / 2 |
| `accounts/fireworks/models/kimi-k2p6` | 2 / 2 | 2 / 2 | 2 / 2 | 2 / 2 |

Without no-reasoning, both models sometimes emitted visible planning text rather
than the required JSON: Flash produced one valid request in two attempts, while
Kimi produced zero. The harness therefore enables that option for these two
specific models. Raw tool responses are ignored under
`artifacts/agent-toolbox-smoke/`.

Reproduce after building the binaries:

```sh
go run ./cmd/agent-toolbox-smoke --model accounts/fireworks/models/deepseek-v4-flash
go run ./cmd/agent-toolbox-smoke --model accounts/fireworks/models/kimi-k2p6
```

This demonstrates only that the constrained interface can be used in a short,
favorable loop. A paper evaluation must use a frozen, broader question set and
measure evidence retrieval and answer quality separately.

The follow-on 50-loop check is recorded in
[agent-toolbox-smoke-20260724.md](agent-toolbox-smoke-20260724.md). Its sampled
candidate question IDs are now consumed for held-out evaluation use.

The separate, source-grounded development question study is recorded in
[agent-toolbox-development-20260724.md](agent-toolbox-development-20260724.md).
It uses only the ten currently indexed pilot documents and keeps all questions,
answers, evidence, and raw outputs ignored locally.
