-- Isolated storage for prompt-iteration studies.  These rows deliberately do
-- not share the primary facts table, whose document-level replacement semantics
-- are appropriate for production extraction but unsafe for research variants.
CREATE TABLE fact_rephrasing_study_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  namespace TEXT NOT NULL,
  version TEXT NOT NULL,
  iteration INTEGER NOT NULL CHECK (iteration > 0),
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  model TEXT NOT NULL,
  prompt_hash TEXT NOT NULL,
  config JSONB NOT NULL,
  raw_output TEXT NOT NULL,
  raw_output_hash TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'complete',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(namespace, version, iteration, document_id)
);

CREATE TABLE fact_rephrasing_study_facts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id UUID NOT NULL REFERENCES fact_rephrasing_study_runs(id) ON DELETE CASCADE,
  ordinal INTEGER NOT NULL,
  text TEXT NOT NULL,
  quote TEXT NOT NULL,
  quote_char_start INTEGER NOT NULL,
  quote_char_end INTEGER NOT NULL,
  turn_start INTEGER NOT NULL,
  turn_end INTEGER NOT NULL,
  speaker_label TEXT,
  content_hash TEXT NOT NULL,
  UNIQUE(run_id, ordinal)
);
CREATE INDEX fact_rephrasing_study_facts_run_idx ON fact_rephrasing_study_facts(run_id, ordinal);
