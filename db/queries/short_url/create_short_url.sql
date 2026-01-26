-- name: CreateShortURL :one
INSERT INTO short_urls (code, long_url, owner_id)
VALUES ($1, $2, $3)
RETURNING id, code, long_url, owner_id, created_at, updated_at;

-- name: UpdateShortURLCode :one
UPDATE short_urls
SET code = $2
WHERE id = $1
RETURNING id, code, long_url, owner_id, created_at, updated_at;

-- name: GetShortURLByCode :one
SELECT id, code, long_url, owner_id, is_active, expires_at, created_at, updated_at
FROM short_urls
WHERE code = $1;
