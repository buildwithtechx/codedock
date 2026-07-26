ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

CREATE TABLE IF NOT EXISTS refresh_token_revocations (
	id TEXT PRIMARY KEY,
	token_hash TEXT UNIQUE NOT NULL,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	expires_at DATETIME NOT NULL,
	revoked_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_refresh_token_revocations_token_hash ON refresh_token_revocations (token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_token_revocations_user_id ON refresh_token_revocations (user_id);
