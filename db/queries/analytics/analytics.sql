-- name: InsertClick :exec
INSERT INTO clicks (short_url_id, ip_hash, country, city, region, browser, os, device_type, referrer, referrer_domain, utm_source, utm_medium, utm_campaign, is_unique, is_bot)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);

-- name: GetClicksByShortURLID :many
SELECT * FROM clicks WHERE short_url_id = $1 ORDER BY clicked_at DESC LIMIT $2 OFFSET $3;

-- name: CountClicksByShortURLID :one
SELECT COUNT(*) FROM clicks WHERE short_url_id = $1;

-- name: CountUniqueClicksByShortURLID :one
SELECT COUNT(*) FROM clicks WHERE short_url_id = $1 AND is_unique = TRUE;

-- name: GetClicksByCountry :many
SELECT country, COUNT(*) as clicks FROM clicks WHERE short_url_id = $1 AND clicked_at >= $2 GROUP BY country ORDER BY clicks DESC LIMIT $3;

-- name: GetClicksByCity :many
SELECT city, country, COUNT(*) as clicks FROM clicks WHERE short_url_id = $1 AND clicked_at >= $2 AND city IS NOT NULL GROUP BY city, country ORDER BY clicks DESC LIMIT $3;

-- name: GetClicksByDevice :many
SELECT device_type, COUNT(*) as clicks FROM clicks WHERE short_url_id = $1 AND clicked_at >= $2 GROUP BY device_type;

-- name: GetClicksByBrowser :many
SELECT browser, COUNT(*) as clicks FROM clicks WHERE short_url_id = $1 AND clicked_at >= $2 GROUP BY browser ORDER BY clicks DESC;

-- name: GetClicksByOS :many
SELECT os, COUNT(*) as clicks FROM clicks WHERE short_url_id = $1 AND clicked_at >= $2 GROUP BY os ORDER BY clicks DESC;

-- name: GetClicksTimeseries :many
SELECT DATE_TRUNC('day', clicked_at)::TIMESTAMPTZ as date, COUNT(*) as clicks FROM clicks WHERE short_url_id = $1 AND clicked_at >= $2 GROUP BY DATE_TRUNC('day', clicked_at) ORDER BY date;

-- name: GetClicksTimeseriesHourly :many
SELECT DATE_TRUNC('hour', clicked_at)::TIMESTAMPTZ as date, COUNT(*) as clicks FROM clicks WHERE short_url_id = $1 AND clicked_at >= $2 GROUP BY DATE_TRUNC('hour', clicked_at) ORDER BY date;

-- name: GetTopReferrers :many
SELECT referrer_domain, COUNT(*) as clicks FROM clicks WHERE short_url_id = $1 AND clicked_at >= $2 AND referrer_domain IS NOT NULL GROUP BY referrer_domain ORDER BY clicks DESC LIMIT $3;

-- name: GetTopCampaigns :many
SELECT utm_campaign, COUNT(*) as clicks FROM clicks WHERE short_url_id = $1 AND clicked_at >= $2 AND utm_campaign IS NOT NULL GROUP BY utm_campaign ORDER BY clicks DESC LIMIT $3;

-- name: GetClicksSummary :one
SELECT
    COUNT(*) as total_clicks,
    COUNT(*) FILTER (WHERE is_unique = TRUE) as unique_clicks,
    COUNT(*) FILTER (WHERE is_bot = TRUE) as bot_clicks
FROM clicks
WHERE short_url_id = $1 AND clicked_at >= $2;

-- name: GetClicksSummaryAllTime :one
SELECT
    COUNT(*) as total_clicks,
    COUNT(*) FILTER (WHERE is_unique = TRUE) as unique_clicks,
    COUNT(*) FILTER (WHERE is_bot = TRUE) as bot_clicks
FROM clicks
WHERE short_url_id = $1;
