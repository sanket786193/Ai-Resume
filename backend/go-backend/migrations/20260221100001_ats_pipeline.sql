-- +goose Up
-- +goose StatementBegin
-- Parsed resume data: raw text, structured fields, cleaned text (best practice: keep separately)
CREATE TABLE IF NOT EXISTS resume_parsed_data (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
  raw_text TEXT,
  parsed_json JSONB,
  cleaned_text TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(resume_id)
);
CREATE INDEX IF NOT EXISTS idx_resume_parsed_data_resume_id ON resume_parsed_data(resume_id);

-- Full AI evaluation on ats_records (ATS score 0-100, missing skills, summary, model version)
ALTER TABLE ats_records
  ADD COLUMN IF NOT EXISTS ats_score INTEGER,
  ADD COLUMN IF NOT EXISTS skill_match_pct INTEGER,
  ADD COLUMN IF NOT EXISTS missing_skills JSONB,
  ADD COLUMN IF NOT EXISTS experience_match TEXT,
  ADD COLUMN IF NOT EXISTS ai_summary TEXT,
  ADD COLUMN IF NOT EXISTS model_version VARCHAR(64);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ats_records
  DROP COLUMN IF EXISTS ats_score,
  DROP COLUMN IF EXISTS skill_match_pct,
  DROP COLUMN IF EXISTS missing_skills,
  DROP COLUMN IF EXISTS experience_match,
  DROP COLUMN IF EXISTS ai_summary,
  DROP COLUMN IF EXISTS model_version;
DROP TABLE IF EXISTS resume_parsed_data;
-- +goose StatementEnd
