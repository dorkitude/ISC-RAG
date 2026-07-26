#!/usr/bin/env bash
set -euo pipefail

mkdir -p bin
go build -o bin/agent-toolbox-fixed-rag ./cmd/agent-toolbox-fixed-rag
go build -o bin/agent-toolbox-pilot-fixed-rag ./cmd/agent-toolbox-pilot-fixed-rag
go build -o bin/agent-toolbox-turn-rag ./cmd/agent-toolbox-turn-rag
go build -o bin/agent-toolbox-semantic-rag ./cmd/agent-toolbox-semantic-rag
go build -o bin/agent-toolbox-named-fact-rag ./cmd/agent-toolbox-named-fact-rag
go build -o bin/agent-toolbox-named-fact-quote-rag ./cmd/agent-toolbox-named-fact-quote-rag
go build -o bin/agent-toolbox-turnwise-opinion-fact-rag ./cmd/agent-toolbox-turnwise-opinion-fact-rag
go build -o bin/agent-toolbox-fact-development-fixed-rag ./cmd/agent-toolbox-fact-development-fixed-rag

printf '%s\n' 'Built capability-limited RAG toolboxes in bin/'
