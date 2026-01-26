-- +goose Up
CREATE TABLE short_urls (
  id          BIGSERIAL PRIMARY KEY,
  code        VARCHAR(10) NOT NULL UNIQUE,
  long_url    TEXT NOT NULL,
  owner_id    BIGINT NULL,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  expires_at  TIMESTAMPTZ NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_short_urls_owner_id ON short_urls (owner_id);
CREATE INDEX idx_short_urls_expires_at ON short_urls (expires_at);

-- trigger to update updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_short_urls_set_updated_at
BEFORE UPDATE ON short_urls
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_short_urls_set_updated_at ON short_urls;
DROP FUNCTION IF EXISTS set_updated_at;
DROP TABLE IF EXISTS short_urls;
