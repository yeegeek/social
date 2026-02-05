# SNS 社交网络 - 数据库设计文档

## 数据库选型

### PostgreSQL（关系型数据）
用于存储结构化数据，需要事务支持和复杂查询的模块：
- 用户信息、认证
- 好友关系
- 角色权限
- 交易记录
- 订单信息

### MongoDB（文档型数据）
用于存储半结构化、高频读写的数据：
- 消息记录（私信、对话）
- 动态/帖子内容
- 通知记录
- 访客记录
- 日志数据

### Redis（缓存和队列）
用于：
- 会话缓存
- 在线用户列表
- 热点数据缓存
- 消息队列
- 分布式锁

---

## PostgreSQL 表结构设计

### 1. 用户认证与授权模块

#### users 表（已存在，需扩展）
```sql
-- 基础字段已存在于脚手架
-- 需要添加的字段：
ALTER TABLE users ADD COLUMN IF NOT EXISTS username VARCHAR(50) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS gender VARCHAR(10);
ALTER TABLE users ADD COLUMN IF NOT EXISTS birthday DATE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS country VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS city VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS language VARCHAR(10) DEFAULT 'en';
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_vip BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS vip_expires_at TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_online BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_active_at TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'active';
ALTER TABLE users ADD COLUMN IF NOT EXISTS coins INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS fingerprint JSONB;
```

#### oauth_providers 表
```sql
CREATE TABLE oauth_providers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- facebook, wechat
    provider_user_id VARCHAR(255) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_user_id)
);
CREATE INDEX idx_oauth_user_id ON oauth_providers(user_id);
```

#### verification_codes 表
```sql
CREATE TABLE verification_codes (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255),
    phone VARCHAR(20),
    code VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL, -- register, reset_password, login
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_verification_email ON verification_codes(email, type);
CREATE INDEX idx_verification_phone ON verification_codes(phone, type);
```

### 2. 好友系统

#### friendships 表
```sql
CREATE TABLE friendships (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL, -- pending, accepted, rejected, blocked
    group_name VARCHAR(100),
    remark VARCHAR(255),
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, friend_id)
);
CREATE INDEX idx_friendship_user ON friendships(user_id, status);
CREATE INDEX idx_friendship_friend ON friendships(friend_id, status);
```

#### blacklist 表
```sql
CREATE TABLE blacklist (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, blocked_user_id)
);
CREATE INDEX idx_blacklist_user ON blacklist(user_id);
```

### 3. 媒体文件管理

#### media_files 表
```sql
CREATE TABLE media_files (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_type VARCHAR(20) NOT NULL, -- image, video, audio
    file_url TEXT NOT NULL,
    thumbnail_url TEXT,
    file_size BIGINT,
    width INTEGER,
    height INTEGER,
    duration INTEGER, -- for video/audio in seconds
    is_private BOOLEAN DEFAULT FALSE,
    is_paid BOOLEAN DEFAULT FALSE,
    price INTEGER DEFAULT 0, -- in coins
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    audit_result JSONB,
    face_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_media_user ON media_files(user_id);
CREATE INDEX idx_media_status ON media_files(status);
```

#### media_purchases 表
```sql
CREATE TABLE media_purchases (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    media_id BIGINT NOT NULL REFERENCES media_files(id) ON DELETE CASCADE,
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, media_id)
);
CREATE INDEX idx_purchase_user ON media_purchases(user_id);
```

### 4. 支付系统

#### transactions 表
```sql
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- recharge, purchase, vip_upgrade, gift, withdrawal
    amount DECIMAL(10, 2) NOT NULL,
    coins INTEGER,
    currency VARCHAR(10) DEFAULT 'USD',
    payment_method VARCHAR(50), -- paypal, wechat, stripe
    payment_id VARCHAR(255),
    status VARCHAR(20) DEFAULT 'pending', -- pending, completed, failed, refunded
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_transaction_user ON transactions(user_id);
CREATE INDEX idx_transaction_status ON transactions(status);
```

#### products 表
```sql
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- coins, vip, gift
    price DECIMAL(10, 2) NOT NULL,
    coins INTEGER,
    vip_days INTEGER,
    image_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_product_type ON products(type, is_active);
```

#### orders 表
```sql
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT REFERENCES products(id),
    recipient_id BIGINT REFERENCES users(id), -- for gifts
    order_number VARCHAR(100) UNIQUE NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, paid, delivered, cancelled
    payment_method VARCHAR(50),
    payment_id VARCHAR(255),
    delivered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_order_user ON orders(user_id);
CREATE INDEX idx_order_number ON orders(order_number);
```

### 5. 评价系统

#### reviews 表
```sql
CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    reviewer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reviewed_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(reviewer_id, reviewed_id)
);
CREATE INDEX idx_review_reviewed ON reviews(reviewed_id);
```

### 6. 推广系统

#### invitation_codes 表
```sql
CREATE TABLE invitation_codes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(20) DEFAULT 'user', -- user, partner, employee
    max_uses INTEGER DEFAULT 1,
    used_count INTEGER DEFAULT 0,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_invitation_code ON invitation_codes(code);
CREATE INDEX idx_invitation_user ON invitation_codes(user_id);
```

#### referrals 表
```sql
CREATE TABLE referrals (
    id BIGSERIAL PRIMARY KEY,
    referrer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referred_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    invitation_code_id BIGINT REFERENCES invitation_codes(id),
    reward_coins INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_referral_referrer ON referrals(referrer_id);
```

### 7. 员工管理

#### employees 表
```sql
CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    employee_number VARCHAR(50) UNIQUE,
    department VARCHAR(100),
    position VARCHAR(100),
    salary_type VARCHAR(20), -- fixed, commission
    base_salary DECIMAL(10, 2),
    commission_rate DECIMAL(5, 2),
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, suspended
    hired_at DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_employee_user ON employees(user_id);
CREATE INDEX idx_employee_status ON employees(status);
```

#### employee_profiles 表
```sql
CREATE TABLE employee_profiles (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    profile_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected, bound
    bound_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_employee_profile_employee ON employee_profiles(employee_id);
```

### 8. 视频通话

#### video_calls 表
```sql
CREATE TABLE video_calls (
    id BIGSERIAL PRIMARY KEY,
    caller_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    callee_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel_id VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'initiated', -- initiated, ringing, active, ended, missed, rejected
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    duration INTEGER DEFAULT 0, -- in seconds
    cost_coins INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_call_caller ON video_calls(caller_id);
CREATE INDEX idx_call_callee ON video_calls(callee_id);
CREATE INDEX idx_call_status ON video_calls(status);
```

### 9. 客服支持

#### complaints 表
```sql
CREATE TABLE complaints (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL, -- harassment, fraud, inappropriate_content, other
    description TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, investigating, resolved, closed
    assigned_to BIGINT REFERENCES users(id),
    resolution TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_complaint_user ON complaints(user_id);
CREATE INDEX idx_complaint_status ON complaints(status);
```

#### support_messages 表
```sql
CREATE TABLE support_messages (
    id BIGSERIAL PRIMARY KEY,
    complaint_id BIGINT NOT NULL REFERENCES complaints(id) ON DELETE CASCADE,
    sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    is_staff BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_support_msg_complaint ON support_messages(complaint_id);
```

### 10. 系统配置

#### system_configs 表
```sql
CREATE TABLE system_configs (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,
    value TEXT,
    type VARCHAR(50) DEFAULT 'string', -- string, number, boolean, json
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_config_key ON system_configs(key);
```

---

## MongoDB 集合设计

### 1. messages 集合（私信）
```javascript
{
    _id: ObjectId,
    conversation_id: String, // 对话唯一标识
    sender_id: Long,
    receiver_id: Long,
    message_type: String, // text, image, video, audio, gift, location
    content: String,
    media_url: String,
    media_thumbnail: String,
    gift_id: Long,
    is_translated: Boolean,
    translations: {
        en: String,
        zh: String,
        ja: String,
        // ...
    },
    is_read: Boolean,
    read_at: Date,
    is_deleted_by_sender: Boolean,
    is_deleted_by_receiver: Boolean,
    is_recalled: Boolean,
    recalled_at: Date,
    created_at: Date,
    updated_at: Date
}
// 索引
db.messages.createIndex({ conversation_id: 1, created_at: -1 })
db.messages.createIndex({ sender_id: 1, created_at: -1 })
db.messages.createIndex({ receiver_id: 1, is_read: 1 })
```

### 2. conversations 集合（对话列表）
```javascript
{
    _id: ObjectId,
    conversation_id: String,
    participants: [Long], // [user_id_1, user_id_2]
    last_message: {
        content: String,
        sender_id: Long,
        created_at: Date
    },
    unread_count: {
        user_id_1: Number,
        user_id_2: Number
    },
    is_archived: {
        user_id_1: Boolean,
        user_id_2: Boolean
    },
    created_at: Date,
    updated_at: Date
}
// 索引
db.conversations.createIndex({ conversation_id: 1 }, { unique: true })
db.conversations.createIndex({ participants: 1 })
```

### 3. message_drafts 集合（草稿）
```javascript
{
    _id: ObjectId,
    user_id: Long,
    receiver_id: Long,
    content: String,
    created_at: Date,
    updated_at: Date
}
// 索引
db.message_drafts.createIndex({ user_id: 1, receiver_id: 1 }, { unique: true })
```

### 4. posts 集合（动态/帖子）
```javascript
{
    _id: ObjectId,
    user_id: Long,
    content: String,
    media: [{
        type: String, // image, video
        url: String,
        thumbnail: String,
        width: Number,
        height: Number
    }],
    location: {
        type: String,
        coordinates: [Number, Number], // [longitude, latitude]
        address: String
    },
    visibility: String, // public, friends, private
    like_count: Number,
    comment_count: Number,
    share_count: Number,
    is_deleted: Boolean,
    created_at: Date,
    updated_at: Date
}
// 索引
db.posts.createIndex({ user_id: 1, created_at: -1 })
db.posts.createIndex({ created_at: -1 })
db.posts.createIndex({ "location.coordinates": "2dsphere" })
```

### 5. post_likes 集合（点赞）
```javascript
{
    _id: ObjectId,
    post_id: ObjectId,
    user_id: Long,
    created_at: Date
}
// 索引
db.post_likes.createIndex({ post_id: 1, user_id: 1 }, { unique: true })
db.post_likes.createIndex({ user_id: 1, created_at: -1 })
```

### 6. post_comments 集合（评论）
```javascript
{
    _id: ObjectId,
    post_id: ObjectId,
    user_id: Long,
    content: String,
    parent_comment_id: ObjectId, // for replies
    like_count: Number,
    is_deleted: Boolean,
    created_at: Date,
    updated_at: Date
}
// 索引
db.post_comments.createIndex({ post_id: 1, created_at: -1 })
db.post_comments.createIndex({ user_id: 1, created_at: -1 })
```

### 7. notifications 集合（通知）
```javascript
{
    _id: ObjectId,
    user_id: Long,
    type: String, // message, like, comment, friend_request, system
    title: String,
    content: String,
    data: Object, // 额外数据
    is_read: Boolean,
    read_at: Date,
    created_at: Date
}
// 索引
db.notifications.createIndex({ user_id: 1, is_read: 1, created_at: -1 })
```

### 8. visitor_logs 集合（访客记录）
```javascript
{
    _id: ObjectId,
    profile_user_id: Long, // 被访问的用户
    visitor_id: Long, // 访客
    visit_count: Number,
    last_visit_at: Date,
    created_at: Date
}
// 索引
db.visitor_logs.createIndex({ profile_user_id: 1, last_visit_at: -1 })
db.visitor_logs.createIndex({ profile_user_id: 1, visitor_id: 1 }, { unique: true })
```

### 9. api_logs 集合（API日志）
```javascript
{
    _id: ObjectId,
    user_id: Long,
    method: String,
    path: String,
    status_code: Number,
    request_body: Object,
    response_body: Object,
    ip_address: String,
    user_agent: String,
    duration_ms: Number,
    created_at: Date
}
// 索引
db.api_logs.createIndex({ user_id: 1, created_at: -1 })
db.api_logs.createIndex({ created_at: -1 })
db.api_logs.createIndex({ path: 1, created_at: -1 })
```

---

## Redis 数据结构

### 1. 会话缓存
```
Key: session:{user_id}
Type: Hash
Fields: {
    token: string,
    refresh_token: string,
    expires_at: timestamp,
    device_info: json
}
TTL: 7 days
```

### 2. 在线用户
```
Key: online_users
Type: Sorted Set
Score: last_active_timestamp
Member: user_id
```

### 3. 用户缓存
```
Key: user:{user_id}
Type: Hash
Fields: {user object fields}
TTL: 1 hour
```

### 4. 未读消息计数
```
Key: unread_messages:{user_id}
Type: Hash
Fields: {
    conversation_id: count
}
```

### 5. 消息队列
```
Key: queue:{queue_name}
Type: List
Values: json encoded tasks
```

### 6. 速率限制
```
Key: rate_limit:{ip}:{endpoint}
Type: String (counter)
TTL: 1 minute
```

### 7. 验证码
```
Key: verification:{email/phone}:{type}
Type: String
Value: code
TTL: 10 minutes
```

### 8. 分布式锁
```
Key: lock:{resource}
Type: String
Value: lock_id
TTL: 30 seconds
```

---

## 数据迁移策略

1. **PostgreSQL 迁移**：使用 golang-migrate，所有表结构变更都通过 SQL 迁移文件管理
2. **MongoDB 索引**：在应用启动时自动创建索引
3. **Redis**：无需迁移，动态数据

## 数据一致性策略

1. **用户基础信息**：PostgreSQL 为主，Redis 缓存
2. **消息数据**：MongoDB 为主，Redis 缓存未读数
3. **好友关系**：PostgreSQL 为主，Redis 缓存在线好友
4. **事务操作**：使用 PostgreSQL 事务保证一致性
5. **异步任务**：使用 RabbitMQ 保证最终一致性
