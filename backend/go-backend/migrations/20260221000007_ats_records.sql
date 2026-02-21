-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ats_records (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
  candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
  resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE RESTRICT,
  status ats_status_enum NOT NULL DEFAULT 'APPLIED',
  skill_match_score DECIMAL(5, 4),
  ranking_score DECIMAL(5, 4),
  qualified BOOLEAN,
  ai_processed_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  UNIQUE(job_id, candidate_id)
);

CREATE INDEX IF NOT EXISTS idx_ats_records_job_id ON ats_records(job_id);
CREATE INDEX IF NOT EXISTS idx_ats_records_candidate_id ON ats_records(candidate_id);
CREATE INDEX IF NOT EXISTS idx_ats_records_status ON ats_records(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ats_records;
-- +goose StatementEnd
