-- +goose Up
-- +goose StatementBegin
-- Experience mismatch warnings and keyword vs semantic match view for HR
ALTER TABLE ats_records
  ADD COLUMN IF NOT EXISTS experience_warnings JSONB,
  ADD COLUMN IF NOT EXISTS keyword_matches JSONB,
  ADD COLUMN IF NOT EXISTS semantic_matches JSONB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ats_records
  DROP COLUMN IF EXISTS experience_warnings,
  DROP COLUMN IF EXISTS keyword_matches,
  DROP COLUMN IF EXISTS semantic_matches;
-- +goose StatementEnd
