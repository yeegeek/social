-- Create blacklist table
CREATE TABLE IF NOT EXISTS blacklist (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, blocked_user_id)
);

CREATE INDEX idx_blacklist_user ON blacklist(user_id);
CREATE INDEX idx_blacklist_blocked ON blacklist(blocked_user_id);
