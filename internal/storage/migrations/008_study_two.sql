-- Study #2 is intentionally isolated from primary facts and Study #1.  Raw
-- completions remain local database artifacts and are never committed.
CREATE TABLE study_two_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  study_version TEXT NOT NULL,
  candidate TEXT NOT NULL,
  replicate INTEGER NOT NULL CHECK (replicate > 0),
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  model TEXT NOT NULL,
  effective_prompt TEXT NOT NULL,
  prompt_hash TEXT NOT NULL,
  config JSONB NOT NULL,
  raw_output TEXT NOT NULL,
  raw_output_hash TEXT NOT NULL,
  prompt_tokens INTEGER NOT NULL DEFAULT 0,
  completion_tokens INTEGER NOT NULL DEFAULT 0,
  latency_ms INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(study_version,candidate,replicate,document_id)
);
CREATE TABLE study_two_facts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id UUID NOT NULL REFERENCES study_two_runs(id) ON DELETE CASCADE,
  ordinal INTEGER NOT NULL,
  text TEXT NOT NULL,
  content_hash TEXT NOT NULL,
  turn_start INTEGER NOT NULL,
  turn_end INTEGER NOT NULL,
  speaker_label TEXT,
  UNIQUE(run_id,ordinal)
);
CREATE TABLE study_two_evidence (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  fact_id UUID NOT NULL REFERENCES study_two_facts(id) ON DELETE CASCADE,
  ordinal INTEGER NOT NULL,
  quote TEXT NOT NULL,
  quote_char_start INTEGER NOT NULL,
  quote_char_end INTEGER NOT NULL,
  turn_start INTEGER NOT NULL,
  turn_end INTEGER NOT NULL,
  speaker_label TEXT,
  validation_method TEXT NOT NULL DEFAULT 'exact',
  UNIQUE(fact_id,ordinal)
);
CREATE INDEX study_two_facts_run_idx ON study_two_facts(run_id,ordinal);
