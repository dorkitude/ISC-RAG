#!/usr/bin/env zsh
# Generate one additional, transcript-only Claude Fable question per frozen
# document. Results enter Postgres only after exact quote/turn validation.
set -u

project_dir=${0:A:h:h}
cd "$project_dir"
set -a
source .env
set +a

split=${1:-fact-dev-sample-100-expansion-20260725}
sample_file=${2:-examples/fact-development-sample-100-20260725.txt}
question_types=(pronoun_coreference attribution decontextualization opinion_belief_judgment direct_control)
index=0

while IFS= read -r doc; do
  [[ -z "$doc" || "$doc" == \#* ]] && continue
  index=$((index + 1))
  qtype=${question_types[$(((index - 1) % ${#question_types[@]} + 1))]}
  external_id=$(printf 'fact-dev100x-%03d' "$index")
  present=$(psql "$DATABASE_URL" -At -c "SELECT count(*) FROM questions WHERE external_id='$external_id'")
  [[ "$present" != "0" ]] && continue

  prior=$(psql "$DATABASE_URL" -At -c "SELECT q.text FROM questions q JOIN gold_evidence ge ON ge.question_id=q.id JOIN documents d ON d.id=ge.document_id WHERE q.corpus_split='fact-dev-sample-100-20260725' AND d.external_id='$doc' LIMIT 1")
  turns=$(psql "$DATABASE_URL" -At -F $'\t' -c "SELECT t.ordinal,COALESCE(t.speaker_label,''),replace(replace(t.text,E'\\n',' '),E'\\r',' ') FROM turns t JOIN documents d ON d.id=t.document_id WHERE d.external_id='$doc' ORDER BY t.ordinal")
  prompt=$(printf '%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s' \
    'You are independently authoring ONE source-grounded evaluation question for a RAG study. You may use only the transcript below. Do not use facts, chunks, retrieval outputs, or outside knowledge.' \
    "DOCUMENT_ID: $doc" "REQUIRED_TYPE: $qtype" \
    'The question must be natural, answerable by explicit text, and test the requested type. Pick a different claim/topic from the prior question. Preserve exact attribution. Return ONLY one JSON object with keys: text, reference_answer, document_external_id, turn_start, turn_end, supporting_quote. The supporting_quote must be an exact contiguous quote in one cited turn; cited turns must contain all needed support. Do not add markdown.' \
    "PRIOR_QUESTION_TO_AVOID: $prior" \
    'TRANSCRIPT_TURNS (ordinal<TAB>speaker<TAB>text):' "$turns")
  response=$(claude -p "$prompt" --model fable --output-format json 2>/dev/null)
  payload=$(printf '%s' "$response" | jq -r '.result // empty')
  if ! printf '%s' "$payload" | jq -e '.text and .reference_answer and .document_external_id and (.turn_start|type=="number") and (.turn_end|type=="number") and .supporting_quote' >/dev/null 2>&1; then
    print -u2 "generation failed for $doc"
    continue
  fi
  payload_b64=$(printf '%s' "$payload" | base64 | tr -d '\n')
  inserted=$(psql "$DATABASE_URL" -At -c "WITH x AS (SELECT convert_from(decode('$payload_b64','base64'),'UTF8')::jsonb AS j), v AS (SELECT j,d.id AS document_id,(j->>'turn_start')::int AS turn_start,(j->>'turn_end')::int AS turn_end FROM x JOIN documents d ON d.external_id=j->>'document_external_id' WHERE d.external_id='$doc'), checked AS (SELECT v.* FROM v WHERE turn_start<=turn_end AND (SELECT count(*) FROM turns t WHERE t.document_id=v.document_id AND t.ordinal BETWEEN v.turn_start AND v.turn_end)=(v.turn_end-v.turn_start+1) AND EXISTS (SELECT 1 FROM turns t WHERE t.document_id=v.document_id AND t.ordinal BETWEEN v.turn_start AND v.turn_end AND position(v.j->>'supporting_quote' IN t.text)>0)), upsert AS (INSERT INTO questions(external_id,corpus_split,question_type,text,reference_answer,metadata) SELECT '$external_id','$split','$qtype',j->>'text',j->>'reference_answer',jsonb_build_object('generator_model','claude-fable-5','protocol','fresh transcript-only second question per frozen document','supporting_quote',j->>'supporting_quote') FROM checked ON CONFLICT(external_id) DO UPDATE SET text=EXCLUDED.text,reference_answer=EXCLUDED.reference_answer,metadata=EXCLUDED.metadata RETURNING id) INSERT INTO gold_evidence(question_id,document_id,turn_start,turn_end) SELECT u.id,c.document_id,c.turn_start,c.turn_end FROM upsert u CROSS JOIN checked c ON CONFLICT DO NOTHING RETURNING question_id")
  if [[ -z "$inserted" ]]; then
    print -u2 "validation failed for $doc"
  else
    print "validated $external_id $doc"
  fi
done < "$sample_file"
