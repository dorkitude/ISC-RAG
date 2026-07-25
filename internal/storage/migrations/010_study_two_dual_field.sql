-- Dual-field RAG study: retrieval text is explicitly non-authoritative.
-- Answer generation must dereference the linked exact capsule/evidence only.
CREATE TABLE study_two_dual_field (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  capsule_fact_id UUID NOT NULL REFERENCES study_two_facts(id) ON DELETE CASCADE,
  retrieval_source_unit_id UUID UNIQUE REFERENCES source_units(id) ON DELETE CASCADE,
  retrieval_expansion TEXT NOT NULL,
  expansion_hash TEXT NOT NULL,
  expansion_model TEXT NOT NULL,
  expansion_prompt TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(capsule_fact_id)
);
CREATE INDEX study_two_dual_field_capsule_idx ON study_two_dual_field(capsule_fact_id);
