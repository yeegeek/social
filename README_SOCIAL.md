# Social Network - SNS 社交网络平台

一个功能完整的社交网络平台，基于 Go + React 构建，支持好友系统、消息、动态、支付等核心社交功能。

## 项目概述

本项目是一个现代化的社交网络服务（SNS），提供完整的社交功能，包括用户管理、好友系统、私信、动态发布、视频通话、支付系统等。基于 [uyou-go-api-starter](https://github.com/yeegeek/uyou-go-api-starter) 脚手架构建。

### 技术栈

**后端：**
- Go 1.24+
- Gin Web Framework
- PostgreSQL（关系型数据）
- MongoDB（文档型数据）
- Redis（缓存和队列）
- JWT 认证
- Swagger API 文档

**前端：**
- React 18+
- TypeScript
- shadcn/ui
- TailwindCSS
- React Router
- Axios

## 功能模块

### 已实现功能 ✅

1. **用户认证与授权**
   - 用户注册/登录
   - JWT Token 认证
   - 刷新令牌机制
   - 密码加密（bcrypt）

2. **用户管理**
   - 用户资料管理
   - 头像上传
   - 个人信息编辑
   - VIP 会员系统
   - 在线状态管理
   - 虚拟货币余额

3. **好友系统**
   - 发送好友请求
   - 接受/拒绝好友请求
   - 删除好友
   - 好友备注
   - 好友分组
   - 黑名单管理

4. **数据库设计**
   - 用户表（扩展字段）
   - OAuth 提供商表
   - 验证码表
   - 好友关系表
   - 黑名单表

### 计划实现功能 🚧

5. **消息系统**
   - 私信发送/接收
   - 对话线程管理
   - 消息翻译
   - 草稿保存
   - 消息撤回

6. **动态系统**
   - 发布/删除动态
   - 点赞/评论
   - 图片/视频分享
   - 地理位置

7. **媒体管理**
   - 文件上传（图片/视频）
   - 付费内容
   - 内容审核

8. **支付系统**
   - 虚拟货币充值
   - VIP 升级
   - 付费内容购买
   - 交易记录

9. **通知系统**
   - 实时推送
   - 离线推送
   - 通知历史

10. **视频通话**
    - 通话发起/接受
    - 通话记录
    - 通话计费

11. **翻译服务**
    - 多语言支持
    - 实时翻译

12. **推广系统**
    - 邀请码
    - 推荐奖励

13. **客服支持**
    - 用户投诉
    - 客服回复

14. **管理员系统**
    - 用户管理
    - 内容审核
    - 数据统计

## 快速开始

### 前置要求

- Go 1.24+
- PostgreSQL 14+
- MongoDB 6.0+
- Redis 6.0+
- Node.js 18+ (前端开发)

### 后端启动

1. **配置环境变量**
```bash
export JWT_SECRET="social_network_jwt_secret_key_2026_very_secure_random_string_32_chars_min"
export DATABASE_HOST="localhost"
export DATABASE_USER="social_user"
export DATABASE_PASSWORD="social_pass123"
export DATABASE_NAME="social_network"
```

2. **安装依赖**
```bash
go mod download
```

3. **运行数据库迁移**
```bash
go run cmd/migrate/main.go up
```

4. **启动服务器**
```bash
go run cmd/server/main.go
```

服务器将在 `http://localhost:8080` 启动。

### 访问 API 文档

启动服务器后，访问：
- Swagger UI: http://localhost:8080/swagger/index.html
- 健康检查: http://localhost:8080/health

## API 端点

### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新令牌
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/me` - 获取当前用户信息

### 用户相关
- `GET /api/v1/users/:id` - 获取用户信息
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户

### 好友相关
- `GET /api/v1/friends` - 获取好友列表
- `POST /api/v1/friends/request` - 发送好友请求
- `GET /api/v1/friends/requests` - 获取好友请求列表
- `POST /api/v1/friends/requests/:id/accept` - 接受好友请求
- `POST /api/v1/friends/requests/:id/reject` - 拒绝好友请求
- `DELETE /api/v1/friends/:id` - 删除好友
- `PUT /api/v1/friends/:id/remark` - 设置好友备注
- `PUT /api/v1/friends/:id/group` - 设置好友分组
- `POST /api/v1/friends/:id/block` - 拉黑用户
- `DELETE /api/v1/friends/:id/unblock` - 取消拉黑
- `GET /api/v1/friends/blocked` - 获取黑名单列表

### 管理员相关
- `GET /api/v1/admin/users` - 获取用户列表（管理员）
- `PUT /api/v1/admin/users/:id` - 更新用户信息（管理员）
- `DELETE /api/v1/admin/users/:id` - 删除用户（管理员）

## 项目结构

```
.
├── cmd/                    # 应用程序入口
│   ├── server/            # API 服务器
│   ├── migrate/           # 数据库迁移工具
│   └── scheduler/         # 定时任务调度器
├── internal/              # 内部应用代码
│   ├── auth/             # 认证服务
│   ├── user/             # 用户模块
│   ├── friend/           # 好友模块
│   ├── health/           # 健康检查
│   ├── middleware/       # 中间件
│   ├── errors/           # 错误处理
│   ├── config/           # 配置管理
│   ├── db/               # PostgreSQL 数据库
│   ├── mongodb/          # MongoDB 数据库
│   ├── redis/            # Redis 缓存
│   └── server/           # 服务器路由
├── migrations/            # 数据库迁移文件
├── configs/              # 配置文件
├── api/                  # API 文档和 Swagger
├── docs/                 # 项目文档
│   ├── database_design.md      # 数据库设计文档
│   └── api_architecture.md     # API 架构设计文档
├── frontend/             # 前端代码
├── .env                  # 环境变量配置
├── go.mod                # Go 模块依赖
└── README_SOCIAL.md      # 项目说明
```

## 数据库设计

详细的数据库设计请查看：[数据库设计文档](docs/database_design.md)

### PostgreSQL 表
- `users` - 用户表（扩展字段：username, phone, avatar_url, gender, birthday, country, city, bio, language, is_vip, vip_expires_at, is_online, last_active_at, status, coins, fingerprint）
- `oauth_providers` - OAuth 提供商
- `verification_codes` - 验证码
- `friendships` - 好友关系
- `blacklist` - 黑名单
- `refresh_tokens` - 刷新令牌
- `roles` - 角色
- `user_roles` - 用户角色关系

### MongoDB 集合（计划）
- `messages` - 私信
- `conversations` - 对话列表
- `posts` - 动态/帖子
- `notifications` - 通知
- `visitor_logs` - 访客记录

## 开发指南

### 添加新模块

1. 在 `internal/` 目录下创建新模块目录
2. 实现以下文件：
   - `model.go` - 数据模型
   - `dto.go` - 数据传输对象
   - `repository.go` - 数据访问层
   - `service.go` - 业务逻辑层
   - `handler.go` - HTTP 处理层

3. 在 `internal/server/router.go` 中注册路由
4. 在 `cmd/server/main.go` 中初始化模块

### 数据库迁移

创建新迁移：
```bash
go run cmd/migrate/main.go create add_new_table
```

执行迁移：
```bash
go run cmd/migrate/main.go up
```

回滚迁移：
```bash
go run cmd/migrate/main.go down
```

### 生成 Swagger 文档

```bash
swag init -g cmd/server/main.go -o api/docs
```

## 测试

运行测试：
```bash
go test ./...
```

运行测试并查看覆盖率：
```bash
go test -cover ./...
```

## 安全特性

- JWT Token 认证
- 密码加密（bcrypt）
- IP 限流
- 输入验证
- SQL 注入防护
- XSS 防护
- CSRF 防护

## 性能优化

- 数据库连接池
- Redis 缓存
- 查询优化
- 索引优化
- 分页查询
- 异步任务处理

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 联系方式

- 项目主页: https://github.com/yeegeek/social
- 问题反馈: https://github.com/yeegeek/social/issues

## 致谢

- [uyou-go-api-starter](https://github.com/yeegeek/uyou-go-api-starter) - 项目脚手架
- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [GORM](https://gorm.io/) - ORM 库
- [shadcn/ui](https://ui.shadcn.com/) - UI 组件库
