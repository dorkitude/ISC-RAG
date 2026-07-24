CREATE TABLE questions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  external_id TEXT NOT NULL UNIQUE,
  corpus TEXT NOT NULL DEFAULT 'mediasum',
  corpus_split TEXT NOT NULL DEFAULT 'dev',
  question_type TEXT NOT NULL,
  text TEXT NOT NULL,
  reference_answer TEXT,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE gold_evidence (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  turn_start INTEGER NOT NULL,
  turn_end INTEGER NOT NULL,
  UNIQUE (question_id, document_id, turn_start, turn_end)
);

CREATE TABLE experiment_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  kind TEXT NOT NULL,
  representation TEXT NOT NULL,
  representation_version TEXT NOT NULL,
  embedding_model TEXT NOT NULL,
  generator_model TEXT,
  parameters JSONB NOT NULL DEFAULT '{}'::jsonb,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ,
  status TEXT NOT NULL DEFAULT 'running'
);

CREATE TABLE retrievals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  experiment_run_id UUID NOT NULL REFERENCES experiment_runs(id) ON DELETE CASCADE,
  question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  source_unit_id UUID NOT NULL REFERENCES source_units(id) ON DELETE CASCADE,
  rank INTEGER NOT NULL,
  score DOUBLE PRECISION NOT NULL,
  UNIQUE (experiment_run_id, question_id, rank)
);
CREATE INDEX retrievals_question_idx ON retrievals (question_id);

CREATE TABLE answers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  experiment_run_id UUID NOT NULL REFERENCES experiment_runs(id) ON DELETE CASCADE,
  question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  text TEXT NOT NULL,
  citations JSONB NOT NULL DEFAULT '[]'::jsonb,
  usage JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (experiment_run_id, question_id)
);
