-- 创建消息表
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text', -- text, image, video, audio, gift
    media_url TEXT,
    translated_content TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITHOUT TIME ZONE,
    is_recalled BOOLEAN DEFAULT FALSE,
    recalled_at TIMESTAMP WITHOUT TIME ZONE,
    is_draft BOOLEAN DEFAULT FALSE,
    gift_id INTEGER,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);
CREATE INDEX idx_messages_is_read ON messages(is_read);
CREATE INDEX idx_messages_is_draft ON messages(is_draft);
