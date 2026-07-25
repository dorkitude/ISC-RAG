CREATE TABLE study_two_batch_checkpoints (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  study_version TEXT NOT NULL,
  candidate TEXT NOT NULL,
  replicate INTEGER NOT NULL,
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  batch_ordinal INTEGER NOT NULL,
  config JSONB NOT NULL,
  candidate_json JSONB NOT NULL,
  repair_json JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(study_version,candidate,replicate,document_id,batch_ordinal)
);
