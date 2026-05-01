-- name: AdminListUsers :many
SELECT
    u.id,
    u.email,
    u.name,
    u.tier,
    u.login_type,
    u.created_at,
    COUNT(DISTINCT s.id)::bigint AS url_count,
    COALESCE(SUM(click_counts.total), 0)::bigint AS total_clicks
FROM users u
LEFT JOIN short_urls s ON s.owner_type = 'user' AND s.owner_id = u.id::text AND s.is_deleted = false
LEFT JOIN (
    SELECT short_url_id, COUNT(*)::bigint AS total
    FROM clicks
    GROUP BY short_url_id
) click_counts ON click_counts.short_url_id = s.id
GROUP BY u.id
ORDER BY total_clicks DESC
LIMIT $1 OFFSET $2;

-- name: AdminCountUsers :one
SELECT COUNT(*)::bigint FROM users;

-- name: AdminGetUserURLs :many
SELECT
    s.id,
    s.code,
    s.long_url,
    s.name,
    s.is_active,
    s.is_deleted,
    s.created_at,
    s.expires_at,
    COALESCE(COUNT(c.id), 0)::bigint AS click_count
FROM short_urls s
LEFT JOIN clicks c ON c.short_url_id = s.id
WHERE s.owner_type = 'user' AND s.owner_id = $1
GROUP BY s.id
ORDER BY click_count DESC
LIMIT $2 OFFSET $3;

-- name: AdminCountUserURLs :one
SELECT COUNT(*)::bigint FROM short_urls
WHERE owner_type = 'user' AND owner_id = $1;

-- name: AdminGetPlatformStats :one
SELECT
    (SELECT COUNT(*)::bigint FROM users) AS total_users,
    (SELECT COUNT(*)::bigint FROM short_urls WHERE is_deleted = false) AS total_urls,
    (SELECT COUNT(*)::bigint FROM clicks) AS total_clicks;

-- name: AdminGetUsersByTier :many
SELECT tier, COUNT(*)::bigint AS count
FROM users
GROUP BY tier;
