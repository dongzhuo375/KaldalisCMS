# KaldalisCMS 架构文档

> 无头 CMS（Headless CMS）：Go 后端（Gin + GORM）+ Next.js 16 前端管理面板。
> 后端专注于 RESTful Content API 输出，前端为独立管理界面，第三方消费端通过 API 对接。

标注说明：`[已实现]` = 代码已落地 | `[预留]` = 目录/文件尚未创建，属于规划中

---

```
/KaldalisCMS
│
├── cmd/
│   ├── server/
│   │   ├── main.go                        # 应用主入口（双模式切换：Setup / App）[已实现]
│   │   └── config.go                      # Viper 配置初始化 [已实现]
│   └── configs/
│       ├── config.yaml                    # 主配置文件 [已实现]
│       └── casbin_model.conf              # Casbin RBAC 模型定义 [已实现]
│
├── internal/
│   ├── api/
│   │   ├── errorx/
│   │   │   └── responses.go              # 统一错误/成功响应写入器 [已实现]
│   │   ├── v1/
│   │   │   ├── post.go                   # 公共文章 API（已发布内容分发）[已实现]
│   │   │   ├── admin_post.go             # 后台文章管理 API（草稿/发布工作流）[已实现]
│   │   │   ├── user.go                   # 用户认证 API（登录/注册/登出）[已实现]
│   │   │   ├── system.go                 # 系统管理 API（站点状态/设置）[已实现]
│   │   │   ├── health.go                 # 探活接口（/healthz /readyz）[已实现]
│   │   │   ├── media.go                  # 媒体库 API（上传/列表/删除）[已实现]
│   │   │   ├── setup.go                  # 安装向导 API [已实现]
│   │   │   ├── tag.go                    # 标签管理 API [预留]
│   │   │   ├── category.go              # 分类管理 API [预留]
│   │   │   ├── admin_user.go            # 用户管理 API（CRUD/角色分配）[预留]
│   │   │   └── dto/
│   │   │       ├── post_dto.go           # [已实现]
│   │   │       ├── user_dto.go           # [已实现]
│   │   │       ├── system_dto.go         # [已实现]
│   │   │       ├── tag_dto.go            # [已实现]
│   │   │       ├── media_dto.go          # [已实现]
│   │   │       ├── auth_dto.go           # [已实现]
│   │   │       ├── setup_dto.go          # [已实现]
│   │   │       └── common_dto.go         # 通用响应 DTO（ErrorResponse/MessageResponse）[已实现]
│   │   └── middleware/
│   │       ├── auth.go                   # Cookie/JWT 认证 + CSRF 校验 + OptionalAuth/RequireAuth [已实现]
│   │       ├── casbin.go                 # Casbin 路由级权限中间件（Authorize）[已实现]
│   │       ├── cors.go                   # CORS [已实现]
│   │       └── observability.go          # 请求上下文/结构化访问日志/Prometheus 指标/panic 契约化恢复
│   │
│   ├── router/
│   │   ├── router.go                     # 路由注册 + 依赖注入（NewAppRouter / NewSetupRouter）[已实现]
│   │   ├── env_parse.go                  # 空壳（已迁移到 internal/utils/env.go）[已实现]
│   │   ├── swagger_options.go            # Swagger 运行时选项 [已实现]
│   │   ├── swagger_routes_enabled.go     # Swagger 路由注册（build tag: swagger）[已实现]
│   │   └── swagger_routes_disabled.go    # Swagger 空实现（生产构建）[已实现]
│   │
│   ├── docs/
│   │   └── docs.go                       # Swagger/OpenAPI 生成产物 [已实现]
│   │
│   ├── core/
│   │   ├── entity/
│   │   │   ├── post.go                   # 文章领域实体（Status: Draft/Published）[已实现]
│   │   │   ├── user.go                   # 用户领域实体 [已实现]
│   │   │   ├── tag.go                    # 标签领域实体 [已实现]
│   │   │   └── media_asset.go            # 媒体资产领域实体 [已实现]
│   │   ├── service.go                    # 服务层接口定义（PostService/UserService 等）[已实现]
│   │   ├── repository.go                 # 仓储层接口定义（含 MediaRepository）[已实现]
│   │   ├── auth.go                       # 认证相关接口（SessionManager 等）[已实现]
│   │   ├── authorization.go              # 授权接口（PostAuthorizer + PostPermission 常量）[已实现]
│   │   └── error.go                      # 领域错误定义 + ErrorCode 映射 [已实现]
│   │
│   ├── infra/
│   │   ├── auth/
│   │   │   ├── casbin.go                 # Casbin Enforcer 初始化（GORM Adapter）[已实现]
│   │   │   ├── post_authorizer.go        # CasbinPostAuthorizer 实现（能力策略检查）[已实现]
│   │   │   ├── session.go                # Session/JWT 签发与校验 + CSRF [已实现]
│   │   │   ├── casbin_test.go            # 权限策略矩阵测试（80+ 子用例）[已实现]
│   │   │   └── post_authorizer_test.go   # PostAuthorizer 单元测试 [已实现]
│   │   ├── repository/
│   │   │   └── postgres/
│   │   │       ├── db.go                 # InitDB + AutoMigrate [已实现]
│   │   │       ├── post_repo.go          # 文章数据访问 [已实现]
│   │   │       ├── user_repo.go          # 用户数据访问 [已实现]
│   │   │       ├── tag_repo.go           # 标签数据访问 [已实现]
│   │   │       ├── system_repo.go        # 系统设置数据访问 [已实现]
│   │   │       └── media_repo.go         # 媒体仓储（model↔entity 映射）[已实现]
│   │   └── model/
│   │       ├── post.go                   # 文章 GORM 模型 [已实现]
│   │       ├── user.go                   # 用户 GORM 模型 [已实现]
│   │       ├── tag.go                    # 标签 GORM 模型 [已实现]
│   │       ├── category.go              # 分类 GORM 模型 [已实现]
│   │       ├── setting.go                # 系统设置模型 [已实现]
│   │       ├── media_asset.go            # 媒体资产模型 [已实现]
│   │       └── post_asset.go             # 文章-媒体引用关系模型 [已实现]
│   │
│   ├── service/
│   │   ├── post_service.go               # 文章业务（公共/管理视图 + 发布工作流 + 引用同步）[已实现]
│   │   ├── user_service.go               # 用户认证业务 [已实现]
│   │   ├── tag_service.go                # 标签业务 [已实现]
│   │   ├── system_service.go             # 系统/站点业务 [已实现]
│   │   ├── setup_service.go              # 安装流程（建库/迁移/权限初始化）[已实现]
│   │   ├── media_service.go              # 媒体库业务（上传/引用/孤儿清理）[已实现]
│   │   ├── media_imageutil.go            # 图片宽高检测工具 [已实现]
│   │   └── error_semantics.go            # service 层错误语义归一化 [已实现]
│   │
│   └── utils/
│       ├── env.go                        # 环境变量解析工具（MEDIA_* 等）[已实现]
│       └── ticker.go                     # 定时任务工具（RunTicker）[已实现]
│
├── pkg/
│   ├── auth/
│   │   └── jwt.go                        # JWT 工具 [已实现]
│   └── security/
│       └── csrf.go                       # CSRF 工具 [已实现]
│
├── web/                                   # 前端管理面板（Next.js 16 + React 19）
│   ├── src/
│   │   ├── app/
│   │   │   ├── page.tsx                  # 根页面（重定向到默认 locale）[已实现]
│   │   │   ├── api/v1/[...path]/
│   │   │   │   └── route.ts             # API 反向代理（Next.js Route Handler → Go 后端）[已实现]
│   │   │   │
│   │   │   └── [locale]/                 # i18n 动态路由（next-intl）
│   │   │       ├── layout.tsx            # locale 根布局 [已实现]
│   │   │       ├── page.tsx              # locale 首页 [已实现]
│   │   │       ├── setup/
│   │   │       │   └── page.tsx          # 安装向导页 [已实现]
│   │   │       │
│   │   │       ├── (auth)/               # [路由组] 认证（无侧边栏布局）
│   │   │       │   ├── layout.tsx        # 居中卡片布局 [已实现]
│   │   │       │   ├── login/
│   │   │       │   │   └── page.tsx      # /login [已实现]
│   │   │       │   └── register/
│   │   │       │       └── page.tsx      # /register [已实现]
│   │   │       │
│   │   │       ├── (admin)/admin/        # [路由组] 后台管理（Sidebar 布局）
│   │   │       │   ├── layout.tsx        # 后台布局（Sidebar + Header）[已实现]
│   │   │       │   ├── dashboard/
│   │   │       │   │   └── page.tsx      # /admin/dashboard [已实现]
│   │   │       │   ├── posts/
│   │   │       │   │   ├── page.tsx      # /admin/posts（文章列表）[已实现]
│   │   │       │   │   ├── new/
│   │   │       │   │   │   └── page.tsx  # /admin/posts/new（新建文章）[已实现]
│   │   │       │   │   └── [id]/edit/
│   │   │       │   │       └── page.tsx  # /admin/posts/:id/edit（编辑文章）[已实现]
│   │   │       │   ├── media/
│   │   │       │   │   └── page.tsx      # /admin/media（媒体库）[已实现]
│   │   │       │   ├── users/
│   │   │       │   │   └── page.tsx      # /admin/users（用户管理）[已实现，后端 API 预留]
│   │   │       │   ├── analytics/
│   │   │       │   │   └── page.tsx      # /admin/analytics（数据统计）[已实现，后端 API 预留]
│   │   │       │   └── settings/
│   │   │       │       └── page.tsx      # /admin/settings（站点设置）[已实现]
│   │   │       │
│   │   │       └── (public)/             # [路由组] 前台展示
│   │   │           ├── layout.tsx        # 前台布局（SiteHeader + Footer）[已实现]
│   │   │           ├── page.tsx          # / 首页 [已实现]
│   │   │           └── posts/
│   │   │               ├── page.tsx      # /posts（文章列表）[已实现]
│   │   │               └── [id]/
│   │   │                   └── page.tsx  # /posts/:id（文章详情）[已实现]
│   │   │
│   │   ├── components/
│   │   │   ├── ui/                       # Shadcn 基础组件
│   │   │   │   ├── avatar.tsx            # [已实现]
│   │   │   │   ├── badge.tsx             # [已实现]
│   │   │   │   ├── button.tsx            # [已实现]
│   │   │   │   ├── card.tsx              # [已实现]
│   │   │   │   ├── checkbox.tsx          # [已实现]
│   │   │   │   ├── dropdown-menu.tsx     # [已实现]
│   │   │   │   ├── input.tsx             # [已实现]
│   │   │   │   ├── label.tsx             # [已实现]
│   │   │   │   ├── select.tsx            # [已实现]
│   │   │   │   ├── skeleton.tsx          # [已实现]
│   │   │   │   ├── sonner.tsx            # Toast 通知 [已实现]
│   │   │   │   └── table.tsx             # [已实现]
│   │   │   ├── admin/
│   │   │   │   └── post-editor.tsx       # Markdown 文章编辑器 [已实现]
│   │   │   ├── site/
│   │   │   │   ├── site-header.tsx       # 前台顶栏 [已实现]
│   │   │   │   └── sun-wave-background.tsx # Three.js 动效背景 [已实现]
│   │   │   ├── providers/
│   │   │   │   └── query-provider.tsx    # TanStack React Query Provider [已实现]
│   │   │   ├── system-status-guard.tsx   # 系统状态守卫（Setup/App 模式切换）[已实现]
│   │   │   ├── LanguageSwitcher.tsx      # 语言切换组件 [已实现]
│   │   │   ├── theme-provider.tsx        # 明暗主题 Provider（next-themes）[已实现]
│   │   │   └── theme-toggle.tsx          # 明暗主题切换按钮 [已实现]
│   │   │
│   │   ├── services/                     # API 调用封装层
│   │   │   ├── auth-service.ts           # 认证相关请求 [已实现]
│   │   │   ├── post-service.ts           # 文章相关请求 [已实现]
│   │   │   ├── media-service.ts          # 媒体相关请求 [已实现]
│   │   │   └── system-service.ts         # 系统相关请求 [已实现]
│   │   │
│   │   ├── lib/
│   │   │   ├── api.ts                    # Axios 实例（拦截器/CSRF/Cookie）[已实现]
│   │   │   ├── utils.ts                  # Shadcn cn() 工具 [已实现]
│   │   │   └── types.ts                  # TS 类型定义（Post, User 等）[已实现]
│   │   │
│   │   ├── store/
│   │   │   └── useAuthStore.ts           # Zustand 用户状态（Persist 持久化）[已实现]
│   │   │
│   │   ├── i18n/
│   │   │   ├── routing.ts                # i18n 路由配置（locale 列表/默认值）[已实现]
│   │   │   └── request.ts                # next-intl 请求配置 [已实现]
│   │   │
│   │   ├── messages/
│   │   │   ├── en.json                   # 英文翻译 [已实现]
│   │   │   └── zh-CN.json               # 简体中文翻译 [已实现]
│   │   │
│   │   └── middleware.ts                 # Next.js 中间件（路由守卫 / i18n 重定向）[已实现]
│   │
│   ├── next.config.ts                    # Next.js 配置 [已实现]
│   ├── components.json                   # Shadcn 配置 [已实现]
│   └── package.json
│
├── docs/
│   ├── ERROR_CONTRACT.md                  # API 错误响应规范 [已实现]
│   ├── IMPLEMENTATION_NOTES.md            # 媒体库 + API 文档实现细节 [已实现]
│   └── TODO.md                            # Roadmap 与待办 [已实现]
│
├── go.mod
├── go.sum
├── CLAUDE.md                              # Claude Code 项目指引
└── README.md
```

---

## 关键设计模式

### 双模式启动

`cmd/server/main.go` 根据 `config.yaml` 是否存在有效数据库配置决定运行模式：
- **Setup Mode**：仅注册安装向导路由（`NewSetupRouter`），引导用户完成首次配置
- **App Mode**：注册完整 CMS 路由（`NewAppRouter`），包含认证/授权/内容管理全链路

### 两层 Casbin 授权

| 层级 | 策略类型 | 检查位置 | 示例 |
|------|---------|---------|------|
| 路由层 | `(role, path, method)` | `middleware/casbin.go` | `("admin", "/api/v1/admin/posts", "GET")` |
| 能力层 | `(role, resource, action)` | `infra/auth/post_authorizer.go` | `("admin", "post", "publish")` |

Matcher 硬编码 `r.sub == "super_admin"` 跳过策略检查，角色继承链：`super_admin → admin → user`。

### 错误响应契约

所有 API 错误统一返回 `{code, message, details}` 结构。领域错误定义在 `core/error.go`，HTTP 映射在 `api/errorx/responses.go`。详见 `docs/ERROR_CONTRACT.md`。

### 可观测性链路（本轮新增）

- 入口中间件：`internal/api/middleware/observability.go`
  - `RequestContext`：透传/生成 `request_id`（响应头 `X-Request-Id`）
  - `ObserveHTTP`：结构化访问日志 + HTTP 请求指标
  - `RecoverAsContract`：panic 统一收敛为错误契约输出
- 错误出口：`internal/api/errorx/responses.go` 会把 `request_id` 注入 `details.request_id`
- 路由装配：`internal/router/router.go` 在 App/Setup 模式统一挂载上述中间件

### 媒体生命周期

上传 → 数据库记录 + 磁盘文件 → 文章引用同步（`post_asset`）→ 孤儿清理（每小时 ticker）。

### 前端架构

- **API 代理**：`app/api/v1/[...path]/route.ts` 将前端请求转发到 Go 后端，避免 CORS 问题
- **状态管理**：Zustand + Persist（用户认证态）+ TanStack React Query（服务端数据缓存）
- **国际化**：`next-intl` 驱动，`[locale]` 动态路由，支持 `zh-CN` / `en`

---

## 前端主要依赖

| 依赖 | 用途 |
|------|------|
| Next.js 16 + React 19 | 框架 |
| next-intl | i18n 国际化 |
| TanStack React Query | 异步数据管理 |
| Zustand | 客户端状态管理 |
| Axios | HTTP 请求 |
| Shadcn/ui (Radix) | UI 组件库 |
| @uiw/react-md-editor | Markdown 编辑器 |
| react-hook-form + Zod | 表单验证 |
| Three.js + @react-three | 3D 动效 |
| framer-motion | 动画 |
| next-themes | 明暗主题 |

---

## 作为无头 CMS 尚未实现的核心能力

以下按优先级排列，对应 `docs/TODO.md` 中的具体 DoD：

### 必要（P0）

| 能力 | 现状 | 说明 |
|------|------|------|
| Tag/Category CRUD API | model + service 已有，**缺 HTTP handler 和路由注册** | 文章组织的基础能力 |
| 列表分页 / 搜索 / 筛选 | 所有 list 接口无分页 | 内容量增长后不可用 |
| 公共 API slug 路由 | 仅 `GET /posts/:id`，缺少 `GET /posts/by-slug/:slug` | 前端 SEO 友好路由依赖 |

### 强烈建议（P1）

| 能力 | 现状 | 说明 |
|------|------|------|
| 用户管理 API | 前端页面已占位，后端无实现 | admin 无法管理用户/角色 |
| 数据统计 API | 前端页面已占位，后端无实现 | Dashboard 数据空缺 |
| 权限策略持久化 | 策略散布在 `setup_service` + `router.go` | 需收敛到单一可追溯来源 |
| 缓存层（Redis） | 无实现 | API 性能瓶颈，高频读场景必要 |
| 速率限制 | 无实现 | 公共 API 暴露后防滥用 |
| 数据库迁移策略 | 依赖 AutoMigrate | 生产环境 schema 变更不可控 |

### 远期（P2）

| 能力 | 说明 |
|------|------|
| Webhook / 事件通知 | 内容变更时通知外部系统（SSG rebuild 等），无头 CMS 的核心集成能力 |
| 自定义字段 / 内容类型 | 动态内容模型，超越固定 Post 结构 |
| API Token 认证 | 第三方消费端接入（当前仅 Cookie 认证，不适合无头场景） |
| GraphQL 支持 | 可选的查询层，灵活获取嵌套内容 |
| CDN / 对象存储集成 | 媒体文件外部存储（S3/R2 等） |
