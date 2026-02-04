# SNS 社交网络 - API 架构设计

## API 版本和基础路径

- **版本**: v1
- **基础路径**: `/api/v1`
- **认证方式**: JWT Bearer Token
- **响应格式**: JSON

---

## API 模块划分

### 1. 认证模块 (`/api/v1/auth`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/auth/register` | 用户注册 | ❌ |
| POST | `/auth/login` | 用户登录 | ❌ |
| POST | `/auth/logout` | 用户登出 | ✅ |
| POST | `/auth/refresh` | 刷新令牌 | ❌ |
| POST | `/auth/forgot-password` | 忘记密码 | ❌ |
| POST | `/auth/reset-password` | 重置密码 | ❌ |
| POST | `/auth/verify-email` | 验证邮箱 | ❌ |
| POST | `/auth/send-verification-code` | 发送验证码 | ❌ |
| POST | `/auth/oauth/facebook` | Facebook 登录 | ❌ |
| POST | `/auth/oauth/wechat` | 微信登录 | ❌ |
| GET | `/auth/me` | 获取当前用户信息 | ✅ |
| PUT | `/auth/me` | 更新当前用户信息 | ✅ |
| PUT | `/auth/me/password` | 修改密码 | ✅ |

### 2. 用户模块 (`/api/v1/users`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/users` | 获取用户列表（分页、筛选） | ✅ |
| GET | `/users/:id` | 获取用户详情 | ✅ |
| PUT | `/users/:id` | 更新用户信息（管理员） | ✅ Admin |
| DELETE | `/users/:id` | 删除用户（管理员） | ✅ Admin |
| GET | `/users/search` | 搜索用户 | ✅ |
| GET | `/users/recommended` | 推荐用户 | ✅ |
| GET | `/users/online` | 在线用户列表 | ✅ |
| GET | `/users/birthday-today` | 今日生日用户 | ✅ |
| GET | `/users/latest` | 最新注册用户 | ✅ |
| POST | `/users/:id/report` | 举报用户 | ✅ |
| GET | `/users/:id/reviews` | 获取用户评价 | ✅ |
| POST | `/users/:id/review` | 评价用户 | ✅ |

### 3. 好友模块 (`/api/v1/friends`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/friends` | 获取好友列表 | ✅ |
| POST | `/friends/request` | 发送好友请求 | ✅ |
| GET | `/friends/requests` | 获取好友请求列表 | ✅ |
| POST | `/friends/requests/:id/accept` | 接受好友请求 | ✅ |
| POST | `/friends/requests/:id/reject` | 拒绝好友请求 | ✅ |
| DELETE | `/friends/:id` | 删除好友 | ✅ |
| PUT | `/friends/:id/remark` | 设置好友备注 | ✅ |
| PUT | `/friends/:id/group` | 设置好友分组 | ✅ |
| GET | `/friends/groups` | 获取好友分组列表 | ✅ |
| POST | `/friends/:id/block` | 拉黑用户 | ✅ |
| DELETE | `/friends/:id/unblock` | 取消拉黑 | ✅ |
| GET | `/friends/blocked` | 获取黑名单列表 | ✅ |

### 4. 消息模块 (`/api/v1/messages`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/messages/conversations` | 获取对话列表 | ✅ |
| GET | `/messages/conversations/:id` | 获取对话详情 | ✅ |
| POST | `/messages` | 发送消息 | ✅ |
| GET | `/messages/:conversation_id` | 获取对话消息列表 | ✅ |
| PUT | `/messages/:id/read` | 标记消息已读 | ✅ |
| DELETE | `/messages/:id` | 删除消息 | ✅ |
| POST | `/messages/:id/recall` | 撤回消息 | ✅ |
| POST | `/messages/:id/translate` | 翻译消息 | ✅ |
| GET | `/messages/unread-count` | 获取未读消息数 | ✅ |
| POST | `/messages/drafts` | 保存草稿 | ✅ |
| GET | `/messages/drafts/:user_id` | 获取草稿 | ✅ |
| DELETE | `/messages/drafts/:user_id` | 删除草稿 | ✅ |

### 5. 动态模块 (`/api/v1/posts`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/posts` | 获取动态列表（时间线） | ✅ |
| GET | `/posts/:id` | 获取动态详情 | ✅ |
| POST | `/posts` | 发布动态 | ✅ |
| PUT | `/posts/:id` | 更新动态 | ✅ |
| DELETE | `/posts/:id` | 删除动态 | ✅ |
| GET | `/posts/user/:user_id` | 获取用户动态列表 | ✅ |
| POST | `/posts/:id/like` | 点赞动态 | ✅ |
| DELETE | `/posts/:id/like` | 取消点赞 | ✅ |
| GET | `/posts/:id/likes` | 获取点赞列表 | ✅ |
| POST | `/posts/:id/comments` | 评论动态 | ✅ |
| GET | `/posts/:id/comments` | 获取评论列表 | ✅ |
| DELETE | `/posts/comments/:id` | 删除评论 | ✅ |

### 6. 媒体模块 (`/api/v1/media`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/media/upload` | 上传文件 | ✅ |
| GET | `/media/:id` | 获取文件信息 | ✅ |
| DELETE | `/media/:id` | 删除文件 | ✅ |
| GET | `/media/user/:user_id` | 获取用户媒体列表 | ✅ |
| POST | `/media/:id/purchase` | 购买付费媒体 | ✅ |
| GET | `/media/:id/access` | 访问付费媒体 | ✅ |
| PUT | `/media/:id/visibility` | 设置媒体可见性 | ✅ |

### 7. 通知模块 (`/api/v1/notifications`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/notifications` | 获取通知列表 | ✅ |
| GET | `/notifications/unread-count` | 获取未读通知数 | ✅ |
| PUT | `/notifications/:id/read` | 标记通知已读 | ✅ |
| PUT | `/notifications/read-all` | 标记所有通知已读 | ✅ |
| DELETE | `/notifications/:id` | 删除通知 | ✅ |

### 8. 访客模块 (`/api/v1/visitors`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/visitors` | 获取访客列表 | ✅ |
| POST | `/visitors/:user_id` | 记录访问 | ✅ |

### 9. 支付模块 (`/api/v1/payments`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/payments/products` | 获取商品列表 | ✅ |
| POST | `/payments/recharge` | 充值 | ✅ |
| POST | `/payments/vip/upgrade` | 升级VIP | ✅ |
| GET | `/payments/transactions` | 获取交易记录 | ✅ |
| GET | `/payments/transactions/:id` | 获取交易详情 | ✅ |
| POST | `/payments/webhook/paypal` | PayPal 回调 | ❌ |
| POST | `/payments/webhook/wechat` | 微信支付回调 | ❌ |

### 10. 订单模块 (`/api/v1/orders`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/orders` | 获取订单列表 | ✅ |
| GET | `/orders/:id` | 获取订单详情 | ✅ |
| POST | `/orders` | 创建订单 | ✅ |
| PUT | `/orders/:id/cancel` | 取消订单 | ✅ |
| PUT | `/orders/:id/deliver` | 发货（管理员） | ✅ Admin |

### 11. 礼物模块 (`/api/v1/gifts`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/gifts` | 获取礼物列表 | ✅ |
| POST | `/gifts/send` | 发送礼物 | ✅ |
| GET | `/gifts/received` | 获取收到的礼物 | ✅ |
| GET | `/gifts/sent` | 获取发送的礼物 | ✅ |

### 12. 视频通话模块 (`/api/v1/calls`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/calls/initiate` | 发起通话 | ✅ |
| POST | `/calls/:id/accept` | 接受通话 | ✅ |
| POST | `/calls/:id/reject` | 拒绝通话 | ✅ |
| POST | `/calls/:id/end` | 结束通话 | ✅ |
| GET | `/calls/history` | 获取通话记录 | ✅ |
| GET | `/calls/:id` | 获取通话详情 | ✅ |

### 13. 翻译模块 (`/api/v1/translate`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/translate` | 翻译文本 | ✅ |
| GET | `/translate/languages` | 获取支持的语言列表 | ✅ |

### 14. 推广模块 (`/api/v1/promotions`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/promotions/codes` | 生成邀请码 | ✅ |
| GET | `/promotions/codes` | 获取邀请码列表 | ✅ |
| GET | `/promotions/referrals` | 获取推荐记录 | ✅ |
| GET | `/promotions/stats` | 获取推广统计 | ✅ |

### 15. 客服模块 (`/api/v1/support`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/support/complaints` | 提交投诉 | ✅ |
| GET | `/support/complaints` | 获取投诉列表 | ✅ |
| GET | `/support/complaints/:id` | 获取投诉详情 | ✅ |
| POST | `/support/complaints/:id/messages` | 发送客服消息 | ✅ |
| GET | `/support/complaints/:id/messages` | 获取客服消息列表 | ✅ |
| PUT | `/support/complaints/:id/status` | 更新投诉状态（管理员） | ✅ Admin |

### 16. 员工管理模块 (`/api/v1/employees`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/employees` | 获取员工列表（管理员） | ✅ Admin |
| GET | `/employees/:id` | 获取员工详情（管理员） | ✅ Admin |
| POST | `/employees` | 创建员工（管理员） | ✅ Admin |
| PUT | `/employees/:id` | 更新员工信息（管理员） | ✅ Admin |
| DELETE | `/employees/:id` | 删除员工（管理员） | ✅ Admin |
| GET | `/employees/:id/profiles` | 获取员工档案列表 | ✅ |
| POST | `/employees/:id/profiles` | 申请绑定档案 | ✅ |
| PUT | `/employees/profiles/:id/approve` | 审批档案（管理员） | ✅ Admin |

### 17. 管理员模块 (`/api/v1/admin`)

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/admin/dashboard` | 获取仪表盘数据 | ✅ Admin |
| GET | `/admin/users` | 获取用户列表（高级筛选） | ✅ Admin |
| PUT | `/admin/users/:id/status` | 更新用户状态 | ✅ Admin |
| GET | `/admin/media/pending` | 获取待审核媒体 | ✅ Admin |
| PUT | `/admin/media/:id/approve` | 审核通过媒体 | ✅ Admin |
| PUT | `/admin/media/:id/reject` | 拒绝媒体 | ✅ Admin |
| GET | `/admin/transactions` | 获取交易记录 | ✅ Admin |
| GET | `/admin/statistics` | 获取统计数据 | ✅ Admin |
| GET | `/admin/logs` | 获取系统日志 | ✅ Admin |
| GET | `/admin/configs` | 获取系统配置 | ✅ Admin |
| PUT | `/admin/configs/:key` | 更新系统配置 | ✅ Admin |

### 18. 健康检查和监控

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/health` | 综合健康检查 | ❌ |
| GET | `/health/live` | 存活探针 | ❌ |
| GET | `/health/ready` | 就绪探针 | ❌ |
| GET | `/metrics` | Prometheus 指标 | ❌ |

---

## 响应格式规范

### 成功响应
```json
{
    "success": true,
    "data": {
        // 响应数据
    },
    "message": "操作成功"
}
```

### 分页响应
```json
{
    "success": true,
    "data": {
        "items": [],
        "pagination": {
            "page": 1,
            "page_size": 20,
            "total": 100,
            "total_pages": 5
        }
    }
}
```

### 错误响应
```json
{
    "success": false,
    "error": {
        "code": "ERROR_CODE",
        "message": "错误描述",
        "details": {}
    }
}
```

---

## 错误码规范

| 错误码 | HTTP 状态码 | 描述 |
|--------|-------------|------|
| `UNAUTHORIZED` | 401 | 未授权 |
| `FORBIDDEN` | 403 | 无权限 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `VALIDATION_ERROR` | 400 | 参数验证失败 |
| `DUPLICATE_ENTRY` | 409 | 资源已存在 |
| `INTERNAL_ERROR` | 500 | 服务器内部错误 |
| `SERVICE_UNAVAILABLE` | 503 | 服务不可用 |
| `RATE_LIMIT_EXCEEDED` | 429 | 请求频率超限 |
| `INSUFFICIENT_BALANCE` | 400 | 余额不足 |
| `PAYMENT_FAILED` | 400 | 支付失败 |
| `MEDIA_LOCKED` | 403 | 媒体已锁定（需购买） |

---

## 认证和授权

### JWT Token 结构
```json
{
    "user_id": 123,
    "email": "user@example.com",
    "roles": ["user", "vip"],
    "exp": 1234567890,
    "iat": 1234567890
}
```

### 角色权限
- **user**: 普通用户
- **vip**: VIP 用户
- **employee**: 员工
- **admin**: 管理员
- **super_admin**: 超级管理员

---

## WebSocket 实时通信

### 连接地址
```
ws://localhost:8080/ws?token={jwt_token}
```

### 消息格式
```json
{
    "type": "message_type",
    "data": {}
}
```

### 消息类型
- `new_message`: 新消息
- `message_read`: 消息已读
- `friend_request`: 好友请求
- `friend_accepted`: 好友请求已接受
- `call_initiated`: 通话发起
- `call_accepted`: 通话已接受
- `call_ended`: 通话结束
- `user_online`: 用户上线
- `user_offline`: 用户下线
- `notification`: 系统通知

---

## 速率限制

### 全局限制
- 每个 IP 每分钟最多 60 次请求

### 特定端点限制
- 登录/注册: 每分钟 5 次
- 发送消息: 每分钟 30 次
- 上传文件: 每分钟 10 次
- 发送验证码: 每分钟 1 次

---

## 文件上传

### 支持的文件类型
- **图片**: jpg, jpeg, png, gif, webp
- **视频**: mp4, mov, avi
- **音频**: mp3, wav, aac

### 文件大小限制
- 图片: 最大 10MB
- 视频: 最大 100MB
- 音频: 最大 20MB

### 上传流程
1. 客户端调用 `/api/v1/media/upload` 获取预签名 URL
2. 客户端直接上传到 S3（或其他存储服务）
3. 上传完成后调用 `/api/v1/media` 创建媒体记录

---

## 国际化

### 支持的语言
- `en`: English
- `zh`: 简体中文
- `zh-TW`: 繁体中文
- `ja`: 日本語
- `ko`: 한국어
- `fr`: Français
- `de`: Deutsch
- `it`: Italiano

### 请求头
```
Accept-Language: zh
```

---

## API 版本管理

### 版本策略
- 使用 URL 路径版本控制：`/api/v1`, `/api/v2`
- 向后兼容：旧版本至少保留 6 个月
- 废弃通知：通过响应头 `X-API-Deprecated` 告知

### 版本生命周期
1. **开发中** (Development)
2. **测试中** (Beta)
3. **稳定版** (Stable)
4. **已废弃** (Deprecated)
5. **已下线** (Retired)
