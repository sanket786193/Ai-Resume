-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ats_sessions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES ats_users(id) ON DELETE CASCADE,
  refresh_token_hash VARCHAR(255) NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ats_sessions_user_id ON ats_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_ats_sessions_expires_at ON ats_sessions(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ats_sessions;
-- +goose StatementEnd
