-- name: CreateShortURL :one
INSERT INTO short_urls (code, long_url, owner_id)
VALUES ($1, $2, $3)
RETURNING id, code, long_url, owner_id, created_at, updated_at;

-- name: UpdateShortURLCode :one
UPDATE short_urls
SET code = $2
WHERE id = $1
RETURNING id, code, long_url, owner_id, created_at, updated_at;
