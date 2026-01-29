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
SELECT id, code, long_url, owner_type, owner_id, is_active, is_deleted, expires_at, created_at, updated_at
FROM short_urls
WHERE code = $1 AND is_deleted = FALSE;

-- name: GetShortURLsByOwner :many
SELECT id, code, long_url, owner_type, owner_id, is_active, is_deleted, expires_at, created_at, updated_at
FROM short_urls
WHERE owner_type = $1 AND owner_id = $2 AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountURLsByOwner :one
SELECT COUNT(*) FROM short_urls
WHERE owner_type = $1 AND owner_id = $2 AND is_deleted = FALSE;

-- name: GetURLByCodeAndOwner :one
SELECT id, code, long_url, owner_type, owner_id, is_active, is_deleted, expires_at, created_at, updated_at
FROM short_urls
WHERE code = $1 AND owner_type = $2 AND owner_id = $3 AND is_deleted = FALSE;

-- name: ToggleURLActive :exec
UPDATE short_urls
SET is_active = NOT is_active
WHERE code = $1 AND owner_type = $2 AND owner_id = $3 AND is_deleted = FALSE;

-- name: SoftDeleteURL :exec
UPDATE short_urls
SET is_deleted = TRUE
WHERE code = $1 AND owner_type = $2 AND owner_id = $3 AND is_deleted = FALSE;

-- name: TransferAnonymousURLsToUser :exec
UPDATE short_urls
SET owner_type = 'user', owner_id = $2
WHERE owner_type = 'anonymous' AND owner_id = $1;

-- name: TransferAnonymousURLsToUserWithLimit :exec
UPDATE short_urls
SET owner_type = 'user', owner_id = $2
WHERE id IN (
    SELECT s.id FROM short_urls s
    WHERE s.owner_type = 'anonymous' AND s.owner_id = $1 AND s.is_deleted = FALSE
    ORDER BY s.created_at ASC
    LIMIT $3
);

-- name: CountURLsCreatedThisMonth :one
SELECT COUNT(*) FROM short_urls
WHERE owner_id = $1
AND created_at >= date_trunc('month', CURRENT_TIMESTAMP)
AND is_deleted = FALSE;
