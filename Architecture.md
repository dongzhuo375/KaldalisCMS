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
│   │   │   │   └── plugin.go             # 插件管理API
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go               # [修改] 需支持从 Cookie 读取 Token (适配SSR)
│   │   │   │   ├── cors.go               # [修改] 更新 AllowedOrigins (如 localhost:3000)
│   │   │   │   ├── logger.go             # 日志中间件
│   │   │   │   ├── theme.go              # 主题上下文中间件
│   │   │   │   └── plugin.go             # 插件中间件
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
│
├── 前端应用 (Next.js 14+)
│   ├── web/
│   │   ├── public/                       # 静态资源 (Next.js标准)
│   │   │   ├── favicon.ico
│   │   │   └── images/
│   │   │
│   │   ├── src/
│   │   │   ├── app/                      # [核心] App Router 目录 (替代 router/)
│   │   │   │   ├── (admin)/              # 后台管理路由组 (CSR模式)
│   │   │   │   │   ├── dashboard/
│   │   │   │   │   │   └── page.tsx      # 对应原 Dashboard.vue
│   │   │   │   │   ├── themes/
│   │   │   │   │   │   └── page.tsx      # 对应原 ThemeManagement.vue
│   │   │   │   │   ├── posts/
│   │   │   │   │   │   ├── page.tsx      # 对应原 PostManagement.vue
│   │   │   │   │   │   └── editor/
│   │   │   │   │   │       └── page.tsx  # 对应原 PostEditor.vue
│   │   │   │   │   ├── login/
│   │   │   │   │   │   └── page.tsx      # 对应原 LoginForm.vue
│   │   │   │   │   └── layout.tsx        # 后台通用布局 (Sidebar/Header)
│   │   │   │   │
│   │   │   │   ├── (public)/             # 前台展示路由组 (SSR/SSG模式)
│   │   │   │   │   ├── [slug]/           # 文章详情动态路由
│   │   │   │   │   │   └── page.tsx      # 对应原主题的文章页逻辑
│   │   │   │   │   ├── page.tsx          # 首页
│   │   │   │   │   └── layout.tsx        # 前台布局 (负责动态加载主题)
│   │   │   │   │
│   │   │   │   ├── api/auth/             # NextAuth 或自定义API路由
│   │   │   │   ├── global.css            # 全局样式
│   │   │   │   └── layout.tsx            # 根布局
│   │   │   │
│   │   │   ├── lib/                      # 工具库 (替代 utils/)
│   │   │   │   ├── api.ts                # HTTP请求 (对应原 request.js，需处理Server/Client)
│   │   │   │   ├── store.ts              # 状态管理 (Zustand/Jotai 替代 Pinia)
│   │   │   │   └── utils.ts              # 通用工具
│   │   │   │
│   │   │   ├── hooks/                    # React Hooks (替代 composables/)
│   │   │   │   ├── use-theme.ts          # 对应原 useTheme.js
│   │   │   │   ├── use-auth.ts           # 对应原 useAuth.js
│   │   │   │   └── use-plugin.ts         # 对应原 usePlugin.js
│   │   │   │
│   │   │   ├── components/               # 组件目录
│   │   │   │   ├── ui/                   # 通用UI库 (Button, Input等)
│   │   │   │   ├── admin/                # 后台专用组件
│   │   │   │   │   ├── ThemeCard.tsx     # 对应原 ThemeCard.vue
│   │   │   │   │   ├── UserList.tsx      # 对应原 UserList.vue
│   │   │   │   │   └── PluginList.tsx    # 对应原 PluginManager.vue
│   │   │   │   └── plugins/              # 前端插件插槽组件
│   │   │   │       └── Registry.tsx      # 插件组件注册表
│   │   │   │
│   │   │   └── themes/                   # [关键] 前端主题目录 (React组件)
│   │   │       ├── default/              # 默认主题
│   │   │       │   ├── components/       # 主题特有组件
│   │   │       │   ├── layouts/          # 主题布局
│   │   │       │   └── theme.config.ts   # 主题配置
│   │   │       └── .gitkeep
│   │   │
│   │   ├── package.json
│   │   ├── next.config.js                # Next.js 配置 (替代 vite.config.js)
│   │   ├── tailwind.config.ts            # 样式配置 (推荐)
│   │   └── tsconfig.json                 # TS 配置
│   │
│   └── .env.local                        # 前端环境变量
│
├── go.mod
├── go.sum
└── README.md
