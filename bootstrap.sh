#!/usr/bin/env bash
# Delegates ISC-RAG bootstrapping to the project's Go CLI.
set -euo pipefail

if ! command -v go >/dev/null 2>&1; then
  printf '%s\n' 'bootstrap: Go is required. Install Go, then rerun this script.' >&2
  exit 1
fi

project_dir="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
cd "$project_dir"
exec go run ./cmd/isc-rag bootstrap "$@"
