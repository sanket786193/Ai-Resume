-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS resume_analyses CASCADE;
DROP TABLE IF EXISTS job_descriptions CASCADE;
DROP TABLE IF EXISTS ocr_jobs CASCADE;
DROP TABLE IF EXISTS document_upload CASCADE;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
