ALTER TABLE users ADD COLUMN plan_type text NOT NULL DEFAULT 'free';
ALTER TABLE users ADD COLUMN stripe_customer_id text;
ALTER TABLE users ADD COLUMN stripe_subscription_id text;
ALTER TABLE users ADD COLUMN stripe_price_id text;
