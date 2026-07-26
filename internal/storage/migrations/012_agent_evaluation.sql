-- Persist fresh-session, capability-limited Q&A trials separately from the
-- one-shot retrieval evaluator. Raw completions are retained for audit.
CREATE TABLE agent_evaluation_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  question_split TEXT NOT NULL,
  model TEXT NOT NULL,
  capability TEXT NOT NULL,
  representation TEXT NOT NULL,
  representation_version TEXT NOT NULL,
  embedding_model TEXT NOT NULL,
  repetitions INTEGER NOT NULL CHECK (repetitions > 0),
  tool_selection_system TEXT NOT NULL,
  answer_system TEXT NOT NULL,
  parameters JSONB NOT NULL DEFAULT '{}'::jsonb,
  status TEXT NOT NULL DEFAULT 'running' CHECK (status IN ('running','complete','failed')),
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ
);
CREATE INDEX agent_evaluation_runs_lookup_idx ON agent_evaluation_runs(question_split, model, capability, started_at DESC);

CREATE TABLE agent_evaluation_attempts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id UUID NOT NULL REFERENCES agent_evaluation_runs(id) ON DELETE CASCADE,
  question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  replicate INTEGER NOT NULL CHECK (replicate > 0),
  session_id UUID NOT NULL DEFAULT gen_random_uuid(),
  tool_request JSONB,
  raw_tool_selection TEXT,
  raw_tool_output JSONB,
  raw_answer_completion TEXT,
  answer_text TEXT,
  citation TEXT,
  tool_request_valid BOOLEAN NOT NULL DEFAULT false,
  evidence_hit BOOLEAN NOT NULL DEFAULT false,
  citation_correct BOOLEAN NOT NULL DEFAULT false,
  reference_answer_match BOOLEAN NOT NULL DEFAULT false,
  selection_latency_ms INTEGER,
  search_latency_ms INTEGER,
  answer_latency_ms INTEGER,
  end_to_end_latency_ms INTEGER,
  error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ,
  UNIQUE(run_id, question_id, replicate)
);
CREATE INDEX agent_evaluation_attempts_run_question_idx ON agent_evaluation_attempts(run_id, question_id, replicate);

CREATE TABLE agent_evaluation_retrievals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  attempt_id UUID NOT NULL REFERENCES agent_evaluation_attempts(id) ON DELETE CASCADE,
  source_unit_id UUID NOT NULL REFERENCES source_units(id) ON DELETE CASCADE,
  rank INTEGER NOT NULL CHECK (rank > 0),
  score DOUBLE PRECISION NOT NULL,
  UNIQUE(attempt_id, rank)
);
CREATE INDEX agent_evaluation_retrievals_attempt_idx ON agent_evaluation_retrievals(attempt_id, rank);
