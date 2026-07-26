#!/usr/bin/env bash
set -euo pipefail

# 75 questions × 10 fresh sessions × 2 models × 2 RAG capabilities = 3,000
# persisted attempts. The question set is development-only and is imported
# idempotently before each run.
questions='experiments/decontextualized-question-development-20260725/expanded-questions.jsonl'
common=(
  --questions "$questions"
  --question-split agent-dev-20260725
  --repetitions 10
  --workers 12
  --timeout-seconds 30
  --name agent-toolbox-development-20260725-expanded
)

./scripts/build-agent-toolboxes.sh

for model in accounts/fireworks/models/deepseek-v4-flash accounts/fireworks/models/kimi-k2p6; do
  go run ./cmd/agent-toolbox-evaluate "${common[@]}" \
    --model "$model" \
    --binary ./bin/agent-toolbox-pilot-fixed-rag \
    --capability pilot-fixed-chunk-rag \
    --representation fixed \
    --version pilot-100x100-v1

  go run ./cmd/agent-toolbox-evaluate "${common[@]}" \
    --model "$model" \
    --binary ./bin/agent-toolbox-named-fact-rag \
    --capability named-fact-rag \
    --representation facts \
    --version pilot-named-v3
done
