-- 创建交易记录表
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- recharge, vip_upgrade, purchase, withdraw, gift
    amount DECIMAL(10, 2) NOT NULL,
    balance_before DECIMAL(10, 2) NOT NULL,
    balance_after DECIMAL(10, 2) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending', -- pending, completed, failed, cancelled
    payment_method VARCHAR(50), -- paypal, wechat, alipay, stripe
    payment_id VARCHAR(255),
    related_id INTEGER, -- 关联的订单/内容ID
    related_type VARCHAR(50), -- post, vip, gift
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_payment_id ON transactions(payment_id);
