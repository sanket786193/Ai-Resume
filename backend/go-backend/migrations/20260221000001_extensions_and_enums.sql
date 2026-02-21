-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
  CREATE TYPE user_role AS ENUM ('HR', 'CANDIDATE');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
  CREATE TYPE ats_status_enum AS ENUM (
    'APPLIED', 'SCREENING', 'SHORTLISTED', 'INTERVIEW', 'REJECTED', 'HIRED'
  );
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
  CREATE TYPE job_status_enum AS ENUM ('DRAFT', 'PUBLISHED', 'CLOSED');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS job_status_enum;
DROP TYPE IF EXISTS ats_status_enum;
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
