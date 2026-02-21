-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ats_users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255),
  name VARCHAR(255) NOT NULL,
  role user_role NOT NULL DEFAULT 'CANDIDATE',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_ats_users_email ON ats_users(email);
CREATE INDEX IF NOT EXISTS idx_ats_users_role ON ats_users(role);
CREATE INDEX IF NOT EXISTS idx_ats_users_deleted_at ON ats_users(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ats_users;
-- +goose StatementEnd
