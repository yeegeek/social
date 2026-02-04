# Social Network Frontend

基于 React + TypeScript + shadcn/ui 的社交网络前端应用。

## 技术栈

- **React 18** - UI 框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **TailwindCSS** - 样式框架
- **shadcn/ui** - UI 组件库
- **React Router** - 路由管理
- **Axios** - HTTP 客户端

## 项目结构

```
frontend/
├── src/
│   ├── components/     # 可复用组件
│   ├── pages/          # 页面组件
│   ├── services/       # API 服务
│   ├── hooks/          # 自定义 Hooks
│   ├── types/          # TypeScript 类型定义
│   ├── utils/          # 工具函数
│   ├── App.tsx         # 根组件
│   ├── main.tsx        # 入口文件
│   └── index.css       # 全局样式
├── index.html          # HTML 模板
├── package.json        # 依赖配置
├── vite.config.ts      # Vite 配置
├── tsconfig.json       # TypeScript 配置
├── tailwind.config.js  # TailwindCSS 配置
└── postcss.config.js   # PostCSS 配置
```

## 快速开始

### 安装依赖

```bash
cd frontend
npm install
# 或
pnpm install
```

### 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:3000 启动。

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 已实现功能

### 认证系统
- 用户登录
- 用户注册
- Token 管理
- 自动登出

### 好友系统
- 好友列表
- 发送好友请求
- 接受/拒绝好友请求
- 删除好友
- 拉黑/取消拉黑

### 用户管理
- 个人资料查看
- 个人资料编辑

## 待实现功能

### 页面组件
- [ ] 登录页面 (LoginPage)
- [ ] 注册页面 (RegisterPage)
- [ ] 仪表盘页面 (DashboardPage)
- [ ] 好友页面 (FriendsPage)
- [ ] 个人资料页面 (ProfilePage)
- [ ] 消息页面 (MessagesPage)
- [ ] 动态页面 (FeedPage)

### UI 组件
- [ ] 导航栏 (Navbar)
- [ ] 侧边栏 (Sidebar)
- [ ] 用户卡片 (UserCard)
- [ ] 好友请求卡片 (FriendRequestCard)
- [ ] 消息列表 (MessageList)
- [ ] 动态卡片 (PostCard)

### 功能模块
- [ ] 消息系统
- [ ] 动态发布
- [ ] 通知系统
- [ ] 文件上传
- [ ] 实时通信 (WebSocket)

## API 集成

前端通过 Axios 与后端 API 通信，所有 API 请求都通过 `/api/v1` 代理到后端服务器。

### 认证 API
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/logout` - 登出
- `GET /api/v1/auth/me` - 获取当前用户

### 好友 API
- `GET /api/v1/friends` - 获取好友列表
- `POST /api/v1/friends/request` - 发送好友请求
- `GET /api/v1/friends/requests` - 获取好友请求
- `POST /api/v1/friends/requests/:id/accept` - 接受请求
- `POST /api/v1/friends/requests/:id/reject` - 拒绝请求
- `DELETE /api/v1/friends/:id` - 删除好友
- `POST /api/v1/friends/:id/block` - 拉黑用户
- `DELETE /api/v1/friends/:id/unblock` - 取消拉黑
- `GET /api/v1/friends/blocked` - 获取黑名单

## 开发指南

### 添加新页面

1. 在 `src/pages/` 创建页面组件
2. 在 `src/App.tsx` 中添加路由
3. 根据需要添加权限保护

### 添加新组件

1. 在 `src/components/` 创建组件文件
2. 使用 TypeScript 定义 Props 类型
3. 使用 TailwindCSS 编写样式

### 使用 shadcn/ui 组件

shadcn/ui 组件需要手动添加到项目中：

```bash
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
npx shadcn-ui@latest add dialog
# ... 等等
```

### API 调用示例

```typescript
import { authAPI } from '@/services/api'

// 登录
const handleLogin = async (email: string, password: string) => {
  try {
    const response = await authAPI.login({ email, password })
    const { access_token, user } = response.data
    login(access_token, user)
  } catch (error) {
    console.error('Login failed:', error)
  }
}
```

## 环境变量

创建 `.env.local` 文件配置环境变量：

```env
VITE_API_URL=http://localhost:8080
```

## 样式规范

- 使用 TailwindCSS 工具类
- 遵循 shadcn/ui 设计系统
- 支持深色模式
- 响应式设计

## 代码规范

- 使用 TypeScript 严格模式
- 组件使用函数式组件
- 使用 React Hooks
- 遵循 ESLint 规则

## 测试

```bash
npm run test
```

## 部署

### Vercel 部署

```bash
vercel
```

### Netlify 部署

```bash
netlify deploy
```

### Docker 部署

```bash
docker build -t social-frontend .
docker run -p 3000:3000 social-frontend
```

## 许可证

MIT License
