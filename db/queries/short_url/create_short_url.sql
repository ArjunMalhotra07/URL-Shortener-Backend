-- name: CreateShortURL :one
INSERT INTO short_urls (code, long_url, owner_type, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING id, code, long_url, owner_type, owner_id, created_at, updated_at;

-- name: UpdateShortURLCode :one
UPDATE short_urls
SET code = $2
WHERE id = $1
RETURNING id, code, long_url, owner_type, owner_id, created_at, updated_at;

-- name: GetShortURLByCode :one
SELECT id, code, long_url, owner_type, owner_id, is_active, expires_at, created_at, updated_at
FROM short_urls
WHERE code = $1;

-- name: GetShortURLsByOwner :many
SELECT id, code, long_url, owner_type, owner_id, is_active, expires_at, created_at, updated_at
FROM short_urls
WHERE owner_type = $1 AND owner_id = $2
ORDER BY created_at DESC;

-- name: TransferAnonymousURLsToUser :exec
UPDATE short_urls
SET owner_type = 'user', owner_id = $2
WHERE owner_type = 'anonymous' AND owner_id = $1;
