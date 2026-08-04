-- Persist extraction candidates rejected by the validation gate. The gate
-- previously dropped failures silently, leaving no record of its rejection
-- rate; every gate decision is now auditable per pipeline run.
CREATE TABLE fact_rejections (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  pipeline_run_id UUID NOT NULL REFERENCES pipeline_runs(id) ON DELETE CASCADE,
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  text TEXT NOT NULL,
  quote TEXT NOT NULL,
  reason TEXT NOT NULL CHECK (reason IN ('empty_or_overlong_quote', 'bad_target_range', 'quote_not_in_target', 'generic_speaker', 'duplicate')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX fact_rejections_pipeline_run_id_idx ON fact_rejections (pipeline_run_id);
CREATE INDEX fact_rejections_document_id_idx ON fact_rejections (document_id);
