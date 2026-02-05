-- 创建对话表
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    user1_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_message_id INTEGER,
    last_message_at TIMESTAMP WITHOUT TIME ZONE,
    user1_unread_count INTEGER DEFAULT 0,
    user2_unread_count INTEGER DEFAULT 0,
    user1_deleted BOOLEAN DEFAULT FALSE,
    user2_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT conversations_users_unique UNIQUE (user1_id, user2_id),
    CONSTRAINT conversations_users_order CHECK (user1_id < user2_id)
);

-- 创建索引
CREATE INDEX idx_conversations_user1_id ON conversations(user1_id);
CREATE INDEX idx_conversations_user2_id ON conversations(user2_id);
CREATE INDEX idx_conversations_last_message_at ON conversations(last_message_at DESC);
