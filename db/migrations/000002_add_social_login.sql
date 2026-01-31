-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN login_type SMALLINT NOT NULL DEFAULT 0,
ADD COLUMN google_id TEXT UNIQUE,
ADD COLUMN name TEXT,
ADD COLUMN avatar_url TEXT;
-- +goose StatementEnd

COMMENT ON COLUMN users.login_type IS '0 = email/password, 1 = google';

ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_users_google_id;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
ALTER TABLE users
DROP COLUMN IF EXISTS login_type,
DROP COLUMN IF EXISTS google_id,
DROP COLUMN IF EXISTS name,
DROP COLUMN IF EXISTS avatar_url;
