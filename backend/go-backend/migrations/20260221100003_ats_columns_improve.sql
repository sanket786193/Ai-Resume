-- +goose Up
-- +goose StatementBegin
-- Improve score columns: NUMERIC(5,4) with valid range 0–1; ensure AI columns exist
ALTER TABLE ats_records
  ALTER COLUMN skill_match_score TYPE NUMERIC(5, 4),
  ALTER COLUMN ranking_score TYPE NUMERIC(5, 4);

ALTER TABLE ats_records
  ADD CONSTRAINT chk_ats_skill_match_score CHECK (skill_match_score IS NULL OR (skill_match_score >= 0 AND skill_match_score <= 1));
ALTER TABLE ats_records
  ADD CONSTRAINT chk_ats_ranking_score CHECK (ranking_score IS NULL OR (ranking_score >= 0 AND ranking_score <= 1));

ALTER TABLE ats_records
  ADD COLUMN IF NOT EXISTS model_version VARCHAR(64),
  ADD COLUMN IF NOT EXISTS experience_warnings JSONB,
  ADD COLUMN IF NOT EXISTS keyword_matches JSONB,
  ADD COLUMN IF NOT EXISTS semantic_matches JSONB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ats_records
  DROP CONSTRAINT IF EXISTS chk_ats_skill_match_score,
  DROP CONSTRAINT IF EXISTS chk_ats_ranking_score;
-- +goose StatementEnd
