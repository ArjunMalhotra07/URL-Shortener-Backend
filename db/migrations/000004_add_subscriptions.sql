-- +goose Up
CREATE TYPE subscription_tier AS ENUM ('free', 'pro', 'business');

ALTER TABLE users
ADD COLUMN tier subscription_tier NOT NULL DEFAULT 'free',
ADD COLUMN stripe_customer_id TEXT,
ADD COLUMN subscription_ends_at TIMESTAMPTZ;

CREATE INDEX idx_users_stripe_customer_id ON users(stripe_customer_id) WHERE stripe_customer_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_users_stripe_customer_id;

ALTER TABLE users
DROP COLUMN IF EXISTS tier,
DROP COLUMN IF EXISTS stripe_customer_id,
DROP COLUMN IF EXISTS subscription_ends_at;

DROP TYPE IF EXISTS subscription_tier;
