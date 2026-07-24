CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE embeddings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  source_unit_id UUID NOT NULL REFERENCES source_units(id) ON DELETE CASCADE,
  model TEXT NOT NULL,
  dimensions INTEGER NOT NULL,
  vector vector NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (source_unit_id, model)
);

CREATE INDEX embeddings_source_unit_idx ON embeddings (source_unit_id);
