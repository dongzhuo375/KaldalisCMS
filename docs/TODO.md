# TODO / Roadmap

本文是可执行的待办清单。每条建议包含：目标、范围（涉及文件/模块）、完成标准（DoD）。

## 上线最小要求（P0：阻塞上线）

### P0-1) 固化权限策略来源（不要在 main 里硬编码 AddPolicy）

目标：生产环境权限策略可追溯、可重复部署、可持久化（多实例/重启一致）。

范围（代表文件）：
- `cmd/server/main.go`（当前启动时 `AddPolicy` 写死）
- `internal/infra/auth/casbin.go`（Casbin adapter/初始化）

完成标准（DoD）：
- 策略来源明确（DB/配置/迁移之一），启动时不再散落硬编码；
- 重启后策略保持一致；
- 有最小文档说明如何变更策略。

### P0-2) 生产级 Session/Cookie/CSRF 安全配置收口

目标：确保登录态在生产环境符合最小安全要求。

范围（代表文件）：
- `internal/infra/auth/session.go`
- `pkg/security/csrf.go`
- `internal/api/middleware/auth.go`

完成标准（DoD）：
- Cookie: `HttpOnly`、`Secure`（https）、`SameSite` 策略明确；
- 密钥/盐来自环境变量或配置，禁止默认弱值；
- CSRF 校验在受保护写接口一致生效，并有前端对接说明。

### P0-3) Setup 写配置可靠性（首次部署可完成初始化）

目标：没有 `cmd/configs/config.yaml` 时也能成功完成 Setup 写入 DB 配置并重启到 App Mode。

范围（代表文件）：
- `cmd/server/config.go`（`InitConfig` / `SaveDatabaseConfig`）
- `internal/api/v1/setup.go`、`internal/service/setup_service.go`

完成标准（DoD）：
- 首次运行无配置文件时，Setup 能成功写出配置文件；
- Setup 成功后能热重载进入 App Mode；
- 错误信息明确（写文件失败/权限不足/路径不存在）。

### P0-4) `/media` 静态目录暴露收敛

目标：避免 uploadDir 下非媒体文件被意外暴露。

范围（代表文件）：
- `internal/router/router.go`（`r.Static("/media", uploadDir)`）
- `cmd/server/config.go`（media.upload_dir 配置）

完成标准（DoD）：
- 只暴露必要的子路径（例如仅 `/media/a` 或等价安全策略）；
- 文档说明 uploadDir 的目录约束；
- 覆盖常见路径穿越/误放文件风险。

### P0-5) Health/Ready 探活接口

目标：支持部署平台（容器/反向代理）探活与就绪判断。

范围（代表文件）：
- `internal/router/router.go`（注册新路由）
- `internal/api/v1/system.go` 或新增 `internal/api/v1/health.go`

完成标准（DoD）：
- `/healthz`：进程存活即可 200；
- `/readyz`：至少检查 DB 连接可用；
- 响应结构固定，便于监控。

### P0-6) 文章发布最小工作流（Draft/Published）

目标：实现“草稿→发布→下线”的最小闭环，避免所有内容默认公开。

范围（代表文件）：
- `internal/infra/model/post.go`（增加状态字段/索引）
- `internal/service/post_service.go`、`internal/api/v1/post.go`（写入/筛选逻辑）

完成标准（DoD）：
- 公共 GET 只返回 Published；
- 受保护接口可创建 Draft、发布、下线；
- 列表/详情的返回字段包含状态；
- 最小测试覆盖（public 过滤 + 权限/角色）。

### P0-7) 错误响应格式统一（前端可稳定处理）

目标：所有 API 错误返回统一结构与 code，减少前端兜底与线上排障成本。

范围（代表文件）：
- `internal/api/v1/*.go`（各 handler 的错误返回）
- `internal/core/error.go`（领域错误）

完成标准（DoD）：
- 统一错误结构（例如 `{code,message,details}`）并落地；
- 关键错误（not found / permission / duplicate / validation）映射稳定 code；
- Update/Delete 等端点返回体与状态码约定一致。

## 上线前强烈建议（P1）

### P1-1) 分类/标签最小可用（至少一个维度）

目标：文章具备最基本的组织能力（分类或标签），后台可管理。

范围（代表文件）：
- `internal/service/tag_service.go` / `internal/infra/repository/postgres/tag_repo.go`
- `internal/api/v1`（补齐 `tag`/`category` 的 handler 路由）

完成标准（DoD）：
- 标签（或分类）CRUD API 可用；
- Post 可关联并在列表/详情中返回；
- 最小权限策略与后台入口对齐。

### P1-2) 角色/权限管理入口与初始 Admin 策略

目标：上线后能可控地管理 admin/user，不靠改代码。

范围（代表文件）：
- `internal/service/user_service.go`、`internal/infra/repository/postgres/user_repo.go`
- `cmd/server/main.go`（移除硬编码策略后，初始化 admin/策略的替代方案）

完成标准（DoD）：
- 安装流程确保存在可登录的 admin；
- 具备最小变更角色方式（管理端接口/DB 操作文档皆可）；
- 权限策略与角色含义有文档。

### P1-3) 密码策略与安全基线

目标：密码以 bcrypt 存储并增加最小口令策略。

范围（代表文件）：
- `internal/service/user_service.go`
- `internal/infra/model/user.go`

完成标准（DoD）：
- 注册/安装时密码 bcrypt；
- 最小强度校验（长度/常见弱口令拒绝可选）；
- 登录失败反馈不泄露用户存在性。

### P1-4) 数据库迁移策略（可重复部署）

目标：明确生产环境 schema 变更方式，避免 AutoMigrate 隐式变更。

范围（代表文件）：
- `internal/infra/repository/postgres/db.go`（AutoMigrate 开关/策略）
- `internal/infra/repository/postgres/migration/*`

完成标准（DoD）：
- prod 环境下迁移/AutoMigrate 策略明确（开关 + 文档）；
- 至少能在空库可一键初始化；
- schema 变更可追溯。

### P1-5) 备份与恢复（Postgres + UploadDir）

目标：具备最小可恢复能力。

范围（代表文件）：
- `docs/`（新增运维文档即可）

完成标准（DoD）：
- 说明 Postgres 备份/恢复步骤；
- 说明上传目录备份/恢复步骤；
- 明确恢复顺序与验证方式。

## 媒体库（Media）

### 1) （增强）批量删除 / 批量查询接口

目标：为媒体库增加批量操作接口，并保持 owner 权限与“被引用禁止删除（409）”等约束一致。

范围（代表文件）：
- `internal/api/v1/media.go`（新增 endpoints：批量删除/批量查询）
- `internal/service/media_service.go`（批量权限校验、引用检查、软删策略复用）

完成标准（DoD）：
- 支持批量删除/查询；
- 单个失败如何返回（全成/部分成）有明确约定；
- 权限与引用限制规则与单删一致；
- 有最小覆盖测试（happy path + 引用冲突 + 权限不足）。

### 2) （增强）增加 SHA256 去重策略

目标：支持（同用户/全站）去重策略，并定义重复上传处理策略（复用已有资源/拒绝/生成新版本）。

范围（代表文件）：
- `internal/infra/model/media_asset.go`（增加 sha256 字段、索引/唯一约束策略）
- `internal/infra/repository/postgres/media_repo.go`（按 hash 查询/复用逻辑）

完成标准（DoD）：
- 数据库字段与索引完成；
- 上传流程在可配置策略下复用/拒绝/版本化行为一致；
- 并发上传同一文件不产生重复记录（或可解释的可控重复）。

## 运维一致性

### 3) （增强）后台一致性扫描接口

目标：扫描上传目录与数据库记录的一致性（孤儿文件/孤儿记录），提供只读报告或修复选项。

范围（代表文件）：
- `internal/service/system_service.go`（或新增 `media_audit_service.go`）
- `internal/api/v1/system.go`（新增运维端点）

完成标准（DoD）：
- 输出可读报告（数量、样例、路径）；
- 修复操作（若提供）需有安全开关与权限限制；
- 对大目录有分页/限速/超时策略。

## 前端

### 4) （前端待做）帖子编辑器的媒体选择/插入

目标：在编辑器中提供媒体库选择、上传、插入 Markdown 链接的体验。

范围（代表文件）：
- `web/src/components/admin/post-editor.tsx`
- `web/src/lib/api.ts`

完成标准（DoD）：
- 支持选择/上传媒体并插入到内容；
- 上传失败/权限不足/引用校验提示清晰；
- 与后端媒体 API 对齐。

## 远期模块草案（从架构文档迁移）

以下为“设计方向”，当前仓库结构中未落地对应目录/文件：

- Theme 主题系统：API、repo、service、中间件、前端动态主题组件。
- Plugin 插件系统：后端插件加载、hook/dispatcher、以及 pkg 级 SDK。

当开始实现时建议：先创建最小可运行骨架（接口 + 空实现 + 文档），再逐步扩展。

