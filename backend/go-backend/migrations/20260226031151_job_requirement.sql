-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
  CREATE TYPE experience_level_enum AS ENUM ('FRESHER', 'EXPERIENCED', 'ANY');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

ALTER TABLE jobs
  ADD COLUMN IF NOT EXISTS experience_level experience_level_enum DEFAULT 'ANY',
  ADD COLUMN IF NOT EXISTS qualification VARCHAR(512),
  ADD COLUMN IF NOT EXISTS skills TEXT[] DEFAULT '{}',
  ADD COLUMN IF NOT EXISTS vacancy_limits JSONB DEFAULT '[]';

CREATE INDEX IF NOT EXISTS idx_jobs_experience_level ON jobs(experience_level);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_jobs_experience_level;
ALTER TABLE jobs
  DROP COLUMN IF EXISTS experience_level,
  DROP COLUMN IF EXISTS qualification,
  DROP COLUMN IF EXISTS skills,
  DROP COLUMN IF EXISTS vacancy_limits;
DROP TYPE IF EXISTS experience_level_enum;
-- +goose StatementEnd
