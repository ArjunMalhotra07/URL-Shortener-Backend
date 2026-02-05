-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id, email, tier, created_at;

-- name: CreateGoogleUser :one
INSERT INTO users (email, password_hash, login_type, google_id, name, avatar_url)
VALUES ($1, '', 1, $2, $3, $4)
RETURNING id, email, name, avatar_url, tier, created_at;

-- name: GetUserByGoogleID :one
SELECT id, email, login_type, google_id, name, avatar_url, tier, subscription_ends_at, created_at, updated_at
FROM users
WHERE google_id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, login_type, google_id, name, avatar_url, tier, subscription_ends_at, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, login_type, name, avatar_url, tier, subscription_ends_at, created_at, updated_at
FROM users
WHERE id = $1;

-- name: UpdateUserTier :exec
UPDATE users
SET tier = $2, subscription_ends_at = $3
WHERE id = $1;

-- name: UpdateStripeCustomerID :exec
UPDATE users
SET stripe_customer_id = $2
WHERE id = $1;

-- name: GetUserByStripeCustomerID :one
SELECT id, email, tier, subscription_ends_at
FROM users
WHERE stripe_customer_id = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, expires_at, created_at;

-- name: GetRefreshToken :one
SELECT id, user_id, token_hash, expires_at, created_at
FROM refresh_tokens
WHERE token_hash = $1 AND expires_at > NOW();

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteAllUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at <= NOW();
