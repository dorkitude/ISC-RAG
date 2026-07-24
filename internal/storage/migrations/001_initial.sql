CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE documents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  external_id TEXT NOT NULL UNIQUE,
  corpus TEXT NOT NULL,
  corpus_split TEXT NOT NULL,
  programme TEXT,
  published_on DATE,
  source_url TEXT,
  title TEXT,
  summary TEXT,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE turns (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  ordinal INTEGER NOT NULL,
  speaker_label TEXT,
  text TEXT NOT NULL,
  char_start INTEGER NOT NULL,
  char_end INTEGER NOT NULL,
  UNIQUE (document_id, ordinal)
);
CREATE INDEX turns_document_ordinal_idx ON turns (document_id, ordinal);

CREATE TABLE pipeline_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  kind TEXT NOT NULL,
  version TEXT NOT NULL,
  model TEXT,
  prompt_hash TEXT,
  parameters JSONB NOT NULL DEFAULT '{}'::jsonb,
  status TEXT NOT NULL DEFAULT 'running',
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ
);

CREATE TABLE source_units (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  representation TEXT NOT NULL,
  version TEXT NOT NULL,
  turn_start INTEGER NOT NULL,
  turn_end INTEGER NOT NULL,
  char_start INTEGER NOT NULL,
  char_end INTEGER NOT NULL,
  text TEXT NOT NULL,
  content_hash TEXT NOT NULL,
  pipeline_run_id UUID REFERENCES pipeline_runs(id),
  UNIQUE (document_id, representation, version, content_hash)
);

CREATE TABLE facts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  text TEXT NOT NULL,
  topic TEXT,
  subject_entity TEXT,
  relation_type TEXT,
  target_entity TEXT,
  content_hash TEXT NOT NULL,
  state TEXT NOT NULL DEFAULT 'active',
  pipeline_run_id UUID REFERENCES pipeline_runs(id),
  superseded_by UUID REFERENCES facts(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (document_id, content_hash)
);

CREATE TABLE fact_evidence (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  fact_id UUID NOT NULL REFERENCES facts(id) ON DELETE CASCADE,
  quote TEXT NOT NULL,
  char_start INTEGER NOT NULL,
  char_end INTEGER NOT NULL,
  turn_start INTEGER NOT NULL,
  turn_end INTEGER NOT NULL,
  validation_method TEXT NOT NULL,
  UNIQUE (fact_id, char_start, char_end)
);

CREATE TABLE index_outbox (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  namespace TEXT NOT NULL,
  operation TEXT NOT NULL CHECK (operation IN ('upsert', 'delete')),
  entity_type TEXT NOT NULL,
  entity_id UUID NOT NULL,
  payload JSONB NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'complete', 'failed')),
  attempts INTEGER NOT NULL DEFAULT 0,
  last_error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ
);
CREATE INDEX index_outbox_pending_idx ON index_outbox (status, created_at);
