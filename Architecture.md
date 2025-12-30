/KaldalisCMS
├── 应用入口
│   ├── cmd/
│   │   └── server/
│   │       ├── main.go                    # 应用主入口
│   │       └── config.go                  # 配置初始化
│   ├── configs/
│   │   ├── config.yaml                    # 主配置文件
│   │   ├── config.prod.yaml              # 生产环境配置
│   │   └── config.dev.yaml               # 开发环境配置
│   └── Makefile                          # 构建脚本
│
├── 后端核心
│   ├── internal/
│   │   ├── api/
│   │   │   ├── v1/
│   │   │   │   ├── theme.go              # 主题管理API
│   │   │   │   ├── post.go               # 文章管理API
│   │   │   │   ├── user.go               # 用户管理API
│   │   │   │   ├── upload.go             # 文件上传API
│   │   │   │   ├── system.go             # 系统管理API
│   │   │   │   ├── plugin.go             # 插件管理API
│   │   │   │   └── dto/
│   │   │   │       ├── post_dto.go
│   │   │   │       └── user_dto.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go               # [修改] 需支持从 Cookie 读取 Token (适配SSR)
│   │   │   │   ├── casbin.go             # Casbin权限中间件
│   │   │   │   ├── cors.go               # [修改] 更新 AllowedOrigins (如 localhost:3000)
│   │   │   │   ├── logger.go             # 日志中间件
│   │   │   │   ├── theme.go              # 主题上下文中间件
│   │   │   │   └── plugin.go             # 插件中间件
│   │   ├── router/
│   │   │   └── router.go                 # 路由配置
│   │   │
│   │   ├── core/
│   │   │   ├── entity/   
│   │   │   │   ├── post.go               # 文章实体
│   │   │   │   └── user.go               # 用户实体
│   │   │   │
│   │   │   ├── service.go                # 服务接口
│   │   │   ├── repository.go             # 仓储接口
│   │   │   └── plugin.go                 # 插件核心接口定义
│   │   │
│   │   ├── infra/
│   │   │   ├── auth/
│   │   │   │   ├── casbin.go
│   │   │   │   └── session.go
│   │   │   ├── repository/
│   │   │   │   ├── postgres/
│   │   │   │   │   ├── theme_repo.go         # 主题数据访问
│   │   │   │   │   ├── post_repo.go          # 文章数据访问
│   │   │   │   │   ├── user_repo.go          # 用户数据访问
│   │   │   │   │   └── migration/
│   │   │   │   │       └── 001_init.sql      # 初始迁移
│   │   │   │   │
│   │   │   │   └── redis/
│   │   │   │       └── cache.go              # 缓存服务
│   │   │   │
│   │   │   └── model/
│   │   │       ├── theme.go                  # 主题模型
│   │   │       ├── post.go                   # 文章模型
│   │   │       ├── user.go                   # 用户模型
│   │   │       ├── category.go               # 分类模型
│   │   │       ├── setting.go                # 系统设置模型
│   │   │       └── plugin.go                 # 插件元数据模型
│   │   │    
│   │   ├── service/
│   │   │   ├── theme_service.go          # 主题业务服务
│   │   │   ├── post_service.go           # 文章业务服务
│   │   │   ├── user_service.go           # 用户业务服务
│   │   │   ├── file_service.go           # 文件服务
│   │   │   ├── system_service.go         # 系统服务
│   │   │   └── plugin_service.go         # 插件生命周期管理服务
│   │   │
│   │   ├── theme/
│   │   │   ├── manager.go                # [修改] 逻辑调整：不再负责前端文件分发，仅管理元数据
│   │   │   ├── loader.go                 # 主题加载器
│   │   │   ├── registry.go               # 主题注册表
│   │   │   ├── validator.go              # [修改] 验证规则：检查 .tsx/.jsx 文件而非 .vue
│   │   │   └── cache.go                  # 主题缓存
│   │   │
│   │   ├── plugin/                       # 插件系统核心
│   │   │   ├── manager.go                # 插件管理器
│   │   │   ├── registry.go               # 插件注册表
│   │   │   ├── loader.go                 # 插件加载器
│   │   │   ├── dispatcher.go             # 事件分发器
│   │   │   ├── service.go                # 插件RPC服务
│   │   │   └── hook/                     # 钩子系统
│   │   │       ├── hook.go               # 钩子定义
│   │   │       ├── manager.go            # 钩子管理器
│   │   │       └── types.go              # 钩子类型
│   │   │
│   │   └── utils/
│   │       ├── crypto.go                 # 加密工具
│   │       ├── file.go                   # 文件工具
│   │       └── validator.go              # 验证工具
│   │
│   ├── pkg/
│   │   ├── plugin/                       # 插件SDK
│   │   │   ├── sdk.go                    # 插件SDK主入口
│   │   │   ├── interfaces.go             # 插件接口定义
│   │   │   ├── types.go                  # 共享类型定义
│   │   │   └── util.go                   # SDK工具函数
│   │   │
│   │   └── database/                     # 数据库工具包
│   │       ├── mysql.go
│   │       └── redis.go
│   │
│   └── plugins/                          # 后端插件目录 (Go源码)
│       ├── example/                      
│       │   ├── main.go                   
│       │   ├── go.mod                    
│       │   └── implementation.go         
│       └── .gitkeep                      
├── 前端应用 (Next.js 14+)
│   ├── web/
│   │   ├── public/                       # 静态资源
│   │   │   ├── favicon.ico
│   │   │   └── images/
│   │   │
│   │   ├── src/
│   │   │   ├── app/                      # [核心] App Router 路由
│   │   │   │   ├── layout.tsx            # [关键] 根布局 (包含 <html><body>)
│   │   │   │   ├── globals.css           # 全局样式 (Tailwind v4 配置)
│   │   │   │   │
│   │   │   │   ├── (auth)/               # [路由组] 认证相关 (无侧边栏布局)
│   │   │   │   │   ├── layout.tsx        # 居中卡片布局
│   │   │   │   │   ├── login/
│   │   │   │   │   │   └── page.tsx      # URL: /login
│   │   │   │   │   └── register/
│   │   │   │   │       └── page.tsx      # URL: /register
│   │   │   │   │
│   │   │   │   ├── (admin)/              # [路由组] 后台管理 (权限隔离)
│   │   │   │   │   └── admin/            # 实体路径前缀 /admin
│   │   │   │   │       ├── layout.tsx    # 后台布局 (Sidebar + Header + 登出逻辑)
│   │   │   │   │       ├── dashboard/
│   │   │   │   │       │   └── page.tsx  # URL: /admin/dashboard
│   │   │   │   │       ├── posts/
│   │   │   │   │       │   ├── page.tsx  # URL: /admin/posts (文章列表)
│   │   │   │   │       │   └── [id]/     # 文章编辑/详情
│   │   │   │   │       └── themes/       # 主题管理
│   │   │   │   │           └── page.tsx
│   │   │   │   │
│   │   │   │   └── (public)/             # [路由组] 前台展示 (游客/普通用户)
│   │   │   │       ├── layout.tsx        # 前台布局 (SiteHeader + Footer)
│   │   │   │       ├── page.tsx          # URL: / (首页/欢迎页)
│   │   │   │       └── posts/
│   │   │   │           └── [slug]/
│   │   │   │               └── page.tsx  # URL: /posts/xxx (文章详情)
│   │   │   │
│   │   │   ├── components/               # 组件库
│   │   │   │   ├── ui/                   # [Shadcn] 基础原子组件 (Button, Input, Table...)
│   │   │   │   ├── admin/                # 后台业务组件 (Sidebar, AdminHeader)
│   │   │   │   ├── site/                 # 前台业务组件 (SiteHeader, Hero, FeatureCard)
│   │   │   │   └── themes/               # 动态主题组件 (按需加载)
│   │   │   │
│   │   │   ├── lib/                      # 工具库
│   │   │   │   ├── api.ts                # [核心] Axios 封装 (拦截器, CSRF, Cookie)
│   │   │   │   ├── utils.ts              # Shadcn cn() 工具
│   │   │   │   └── types.ts              # TS 类型定义 (Post, User)
│   │   │   │
│   │   │   ├── store/                    # 状态管理
│   │   │   │   └── useAuthStore.ts       # [核心] Zustand (用户状态 + Persist持久化)
│   │   │   │
│   │   │   ├── hooks/                    # 自定义 Hooks
│   │   │   │   └── use-toast.ts          # Shadcn Toast Hook
│   │   │   │
│   │   │   └── middleware.ts             # [核心] 路由守卫 (未登录拦截 / 角色分流)
│   │   │
│   │   ├── components.json               # Shadcn 配置文件
│   │   ├── next.config.mjs               # Next.js 配置
│   │   ├── package.json
│   │   ├── postcss.config.mjs
│   │   └── .env.local                    # 环境变量 (NEXT_PUBLIC_API_URL)
│
├── go.mod
├── go.sum
└── README.md
