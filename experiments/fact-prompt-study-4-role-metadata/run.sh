#!/usr/bin/env bash
# Run one fixed Study #4 seed/refinement against the locked development document.
# Source-derived transcript, completions, and API response metadata stay ignored
# under artifacts/.  Nothing here writes to the experiment database or main pipeline.
set -euo pipefail

project_dir="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$project_dir"
seed="${1:?usage: run.sh <seed-name>}"
prompt="experiments/fact-prompt-study-4-role-metadata/prompts/${seed}.md"
out="artifacts/fact-prompt-study-4-role-metadata/${seed}"
external_id="NPR-1040"
model="${FACT_STUDY_MODEL:-accounts/fireworks/models/deepseek-v4-pro}"

test -f "$prompt"
mkdir -p "$out"
set -a; source "$project_dir/.env"; set +a
: "${DATABASE_URL:?DATABASE_URL is required in the local .env}"
: "${FIREWORKS_API_KEY:?FIREWORKS_API_KEY is required in the local .env}"

# Keep role labels deterministic: they are only the migrated speaker_name and
# speaker_role fields, never inferred from the dialogue by this runner.
psql "$DATABASE_URL" -At -F $'\t' -c "
 SELECT format('[TURN %s | speaker_name=%s | speaker_role=%s | speaker_label=%s]',
               ordinal, coalesce(speaker_name,''), coalesce(speaker_role,''), speaker_label)
        || E'\\n' || text
 FROM turns t JOIN documents d ON d.id=t.document_id
 WHERE d.corpus='mediasum' AND d.external_id='${external_id}'
 ORDER BY ordinal" >"$out/transcript.txt"

jq -n \
  --arg model "$model" \
  --rawfile system "$prompt" \
  --rawfile transcript "$out/transcript.txt" \
  '{model:$model,temperature:0,max_tokens:6000,reasoning_effort:"none",response_format:{type:"json_object"},messages:[{role:"system",content:$system},{role:"user",content:("<TRANSCRIPT>\\n"+$transcript+"\\n</TRANSCRIPT>\\nReturn at most 25 cards.")}]}' \
  >"$out/request.json"

started="$(date +%s)"
curl --fail-with-body --silent --show-error \
  -H "Authorization: Bearer $FIREWORKS_API_KEY" -H 'Content-Type: application/json' \
  -d @"$out/request.json" 'https://api.fireworks.ai/inference/v1/chat/completions' \
  >"$out/response.json"
ended="$(date +%s)"
jq -r '.choices[0].message.content' "$out/response.json" >"$out/completion.json"
jq -n --arg seed "$seed" --arg model "$model" --argjson seconds "$((ended-started))" \
  --slurpfile response "$out/response.json" \
  '{seed:$seed,model:$model,latency_seconds:$seconds,usage:$response[0].usage,request_sha256:("'"$(shasum -a 256 "$out/request.json" | awk '{print $1}')"'"),response_sha256:("'"$(shasum -a 256 "$out/response.json" | awk '{print $1}')"'")}' \
  >"$out/metadata.json"
echo "wrote $out"
