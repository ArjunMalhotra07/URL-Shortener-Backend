-- +goose Up
CREATE TYPE owner_type_enum AS ENUM ('user', 'anonymous');

-- +goose StatementBegin
CREATE TABLE users (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email         TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE refresh_tokens (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash  TEXT NOT NULL,
  expires_at  TIMESTAMPTZ NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);

-- +goose StatementBegin
CREATE TABLE short_urls (
  id          BIGSERIAL PRIMARY KEY,
  code        VARCHAR(10) NOT NULL UNIQUE,
  long_url    TEXT NOT NULL,
  owner_type  owner_type_enum NOT NULL DEFAULT 'anonymous',
  owner_id    TEXT NOT NULL,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  is_deleted  BOOLEAN NOT NULL DEFAULT FALSE,
  expires_at  TIMESTAMPTZ NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

CREATE INDEX idx_short_urls_owner ON short_urls (owner_type, owner_id);
CREATE INDEX idx_short_urls_expires_at ON short_urls (expires_at);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_short_urls_set_updated_at
BEFORE UPDATE ON short_urls
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;
DROP TRIGGER IF EXISTS trg_short_urls_set_updated_at ON short_urls;
DROP FUNCTION IF EXISTS set_updated_at;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS short_urls;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS owner_type_enum;
