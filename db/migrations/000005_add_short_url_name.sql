-- +goose Up
ALTER TABLE short_urls
ADD COLUMN name TEXT;

-- +goose Down
ALTER TABLE short_urls
DROP COLUMN IF EXISTS name;
