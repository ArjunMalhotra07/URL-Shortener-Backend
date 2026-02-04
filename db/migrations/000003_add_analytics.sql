-- +goose Up
CREATE TABLE clicks (
    id BIGSERIAL PRIMARY KEY,
    short_url_id BIGINT NOT NULL REFERENCES short_urls(id) ON DELETE CASCADE,
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_hash TEXT NOT NULL,
    country TEXT,
    city TEXT,
    region TEXT,
    browser TEXT,
    os TEXT,
    device_type TEXT,
    referrer TEXT,
    referrer_domain TEXT,
    utm_source TEXT,
    utm_medium TEXT,
    utm_campaign TEXT,
    is_unique BOOLEAN DEFAULT FALSE,
    is_bot BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_clicks_short_url_id ON clicks(short_url_id);
CREATE INDEX idx_clicks_clicked_at ON clicks(clicked_at);
CREATE INDEX idx_clicks_country ON clicks(country) WHERE country IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS clicks;
