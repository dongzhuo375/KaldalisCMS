/KaldalisCMS
├── 应用入口
│   ├── cmd/
│   │   ├── server/
│   │   │   ├── main.go                    # 应用主入口 [已实现]
│   │   │   └── config.go                  # 配置初始化（Viper）[已实现]
│   │   └── configs/                       # [差异更新] 当前为 cmd/configs（不在 cmd/configs 根同级）
│   │       ├── config.yaml                # 主配置文件 [已实现]
│   │       └── casbin_model.conf          # Casbin 模型 [已实现]
│   └── Makefile                          # 构建脚本 [预留]（当前仓库根目录未提供）
│
├── 后端核心
│   ├── internal/
│   │   ├── api/
│   │   │   ├── v1/
│   │   │   │   ├── post.go               # 文章管理 API [已实现]
│   │   │   │   ├── user.go               # 用户管理 API [已实现]
│   │   │   │   ├── system.go             # 系统管理 API [已实现]
│   │   │   │   ├── media.go              # 媒体库 API [已实现]
│   │   │   │   ├── theme.go              # 主题管理 API [预留]
│   │   │   │   ├── upload.go             # 通用文件上传 API [预留]（当前以 media.go 为主）
│   │   │   │   ├── plugin.go             # 插件管理 API [预留]
│   │   │   │   └── dto/
│   │   │   │       ├── post_dto.go       # [已实现]
│   │   │   │       ├── user_dto.go       # [已实现]
│   │   │   │       ├── system_dto.go     # [已实现]
│   │   │   │       ├── tag_dto.go        # [已实现]
│   │   │   │       └── media_dto.go      # 媒体 DTO（显式 API contract）[已实现]
│   │   │   └── middleware/
│   │   │       ├── auth.go               # Cookie/JWT/CSRF 认证 [已实现]
│   │   │       ├── casbin.go             # Casbin 权限中间件 [已实现]
│   │   │       ├── cors.go               # CORS [已实现]
│   │   │       ├── logger.go             # 日志中间件 [预留]
│   │   │       ├── theme.go              # 主题上下文中间件 [预留]
│   │   │       └── plugin.go             # 插件中间件 [预留]
│   │   │
│   │   ├── router/
│   │   │   ├── router.go                 # 路由配置/依赖注入 [已实现]
│   │   │   └── env_parse.go              # env 解析（已迁移到 internal/utils/env.go，保留占位）[已实现]
│   │   │
│   │   ├── core/
│   │   │   ├── entity/
│   │   │   │   ├── post.go               # 文章实体 [已实现]
│   │   │   │   ├── user.go               # 用户实体 [已实现]
│   │   │   │   ├── tag.go                # 标签实体 [已实现]
│   │   │   │   └── media_asset.go        # 媒体资产领域实体 [已实现]
│   │   │   ├── service.go                # 服务接口 [已实现]
│   │   │   ├── repository.go             # 仓储接口（含 MediaRepository）[已实现]
│   │   │   ├── error.go                  # core.ErrPermission 等 [已实现]
│   │   │   └── plugin.go                 # 插件核心接口定义 [预留]
│   │   │
│   │   ├── infra/
│   │   │   ├── auth/
│   │   │   │   ├── casbin.go             # Casbin 初始化 [已实现]
│   │   │   │   └── session.go            # Session/JWT/CSRF [已实现]
│   │   │   ├── repository/
│   │   │   │   ├── postgres/
│   │   │   │   │   ├── db.go             # InitDB + AutoMigrate [已实现]
│   │   │   │   │   ├── post_repo.go      # 文章数据访问 [已实现]
│   │   │   │   │   ├── user_repo.go      # 用户数据访问 [已实现]
│   │   │   │   │   ├── tag_repo.go       # 标签数据访问 [已实现]
│   │   │   │   │   ├── system_repo.go    # 系统数据访问 [已实现]
│   │   │   │   │   ├── media_repo.go     # 媒体仓储（model<->entity mapper）[已实现]
│   │   │   │   │   ├── theme_repo.go     # 主题数据访问 [预留]
│   │   │   │   │   └── migration/
│   │   │   │   │       └── 001_init.sql  # 初始迁移 [已实现]
│   │   │   │   └── redis/
│   │   │   │       └── cache.go          # 缓存服务 [预留]
│   │   │   └── model/
│   │   │       ├── post.go               # 文章模型 [已实现]
│   │   │       ├── user.go               # 用户模型 [已实现]
│   │   │       ├── tag.go                # 标签模型 [已实现]
│   │   │       ├── category.go           # 分类模型 [已实现]
│   │   │       ├── setting.go            # 系统设置模型 [已实现]
│   │   │       ├── media_asset.go        # 媒体模型 [已实现]
│   │   │       ├── post_asset.go         # 引用关系模型 [已实现]
│   │   │       ├── theme.go              # 主题模型 [预留]
│   │   │       └── plugin.go             # 插件元数据模型 [预留]
│   │   │
│   │   ├── service/
│   │   │   ├── post_service.go           # 文章业务服务（含引用同步）[已实现]
│   │   │   ├── user_service.go           # 用户业务服务 [已实现]
│   │   │   ├── tag_service.go            # 标签业务服务 [已实现]
│   │   │   ├── system_service.go         # 系统服务 [已实现]
│   │   │   ├── media_service.go          # 媒体库服务 [已实现]
│   │   │   ├── media_imageutil.go        # 图片宽高工具 [已实现]
│   │   │   ├── theme_service.go          # 主题业务服务 [预留]
│   │   │   ├── file_service.go           # 通用文件服务 [预留]
│   │   │   └── plugin_service.go         # 插件生命周期管理服务 [预留]
│   │   │
│   │   ├── theme/                        # 主题系统核心 [预留]
│   │   │   ├── manager.go
│   │   │   ├── loader.go
│   │   │   ├── registry.go
│   │   │   ├── validator.go
│   │   │   └── cache.go
│   │   │
│   │   ├── plugin/                       # 插件系统核心 [预留]
│   │   │   ├── manager.go
│   │   │   ├── registry.go
│   │   │   ├── loader.go
│   │   │   ├── dispatcher.go
│   │   │   ├── service.go
│   │   │   └── hook/
│   │   │       ├── hook.go
│   │   │       ├── manager.go
│   │   │       └── types.go
│   │   │
│   │   └── utils/
│   │       ├── env.go                    # env 工具（解析 MEDIA_*）[已实现]
│   │       ├── crypto.go                 # 加密工具 [预留]
│   │       ├── file.go                   # 文件工具 [预留]
│   │       └── validator.go              # 验证工具 [预留]
│   │
│   ├── pkg/
│   │   ├── auth/
│   │   │   └── jwt.go                    # JWT 工具 [已实现]
│   │   ├── security/
│   │   │   └── csrf.go                   # CSRF 工具 [已实现]
│   │   ├── plugin/                       # 插件 SDK [预留]
│   │   │   ├── sdk.go
│   │   │   ├── interfaces.go
│   │   │   ├── types.go
│   │   │   └── util.go
│   │   └── database/                     # 数据库工具包 [预留]
│   │       ├── mysql.go
│   │       └── redis.go
│   │
│   └── plugins/                          # 后端插件目录（Go源码）[预留]
│       ├── example/
│       │   ├── main.go
│       │   ├── go.mod
│       │   └── implementation.go
│       └── .gitkeep
│
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

---

## 媒体库（Media Library）当前实现要点

- 公共访问：`/media/a/{assetID}/{stored_name}`（静态目录映射到 `MEDIA_UPLOAD_DIR`，默认 `./data/uploads`）
- 物理路径：`{MEDIA_UPLOAD_DIR}/a/{assetID}/{stored_name}`
- 列表权限：admin 看全站；普通用户只看自己的资源。
- 删除权限：普通用户“只能删自己上传的”；否则返回 `core.ErrPermission`（API 映射 403）。
- 删除硬限制：被帖子引用（post_assets）则禁止删除（409）。
- 引用同步：Post Create/Update 解析 Markdown 中的媒体 URL，写入 post_assets。
- API：
  - `POST   /api/v1/media`
  - `GET    /api/v1/media`
  - `DELETE /api/v1/media/:id`
  - `GET    /api/v1/posts/:id/media`

---

## TODO（当前）

1. （可选/增强）统一错误响应格式：例如 `{ code, message, details }`，并为常见错误建立规范化 code。
2. （可选/增强）为媒体库增加“批量删除/批量查询”接口，并保持 owner 权限与引用硬限制一致。
3. （可选/增强）为 media_assets 增加 SHA256 去重策略（同用户/全站可选）与重复上传处理策略。
4. （可选/增强）增加后台运维接口：扫描上传目录与数据库记录的一致性（孤儿文件/孤儿记录）。
5. （前端待做）实现帖子编辑器中的媒体选择/插入体验（本次暂不施工）。
