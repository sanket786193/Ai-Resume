-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
  CREATE TYPE resume_status_enum AS ENUM ('PENDING', 'PROCESSED');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

ALTER TABLE resumes
  ADD COLUMN IF NOT EXISTS status resume_status_enum NOT NULL DEFAULT 'PENDING';

ALTER TABLE interviews
  ADD COLUMN IF NOT EXISTS candidate_confirmed_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_resumes_status ON resumes(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_resumes_status;
ALTER TABLE resumes DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS resume_status_enum;
ALTER TABLE interviews DROP COLUMN IF EXISTS candidate_confirmed_at;
-- +goose StatementEnd
