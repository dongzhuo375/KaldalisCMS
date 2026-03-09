# TODO / Roadmap

本文是可执行的待办清单。每条建议包含：目标、范围（涉及文件/模块）、完成标准（DoD）。

## 上线最小要求（P0：阻塞上线）

### P0-1) 固化权限策略来源（不要在 main 里硬编码 AddPolicy）

目标：生产环境权限策略可追溯、可重复部署、可持久化（多实例/重启一致）。

范围（代表文件）：
- `cmd/server/main.go`（当前启动时 `AddPolicy` 写死）
- `internal/infra/auth/casbin.go`（Casbin adapter/初始化）

完成标准（DoD）：
- **[2026-03-06 更新]** 已将策略初始化抽离到 `setupPolicies` 函数，但仍需研究如何迁入数据库持久化。
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

### P0-3) Setup 流程健壮性 (已完成 - 2026-03-06)

目标：没有 `cmd/configs/config.yaml` 或数据库未初始化时也能成功完成 Setup 写入 DB 配置并重启到 App Mode。

范围（代表文件）：
- `cmd/server/config.go`（`InitConfig` / `SaveDatabaseConfig`）
- `internal/api/v1/setup.go`、`internal/service/setup_service.go`

完成标准（DoD）：
- **[已落地]** 自动创建数据库逻辑；
- **[已落地]** 前端分步预检（Test Connection）机制；
- **[已落地]** 增量保存配置并触发热重启。

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

#

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

### P1-3) 密码策略与安全基线

目标：密码以 bcrypt 存储并增加最小口令策略。

### P1-4) 数据库迁移策略（可重复部署）

目标：明确生产环境 schema 变更方式，避免 AutoMigrate 隐式变更。

### P1-5) 备份与恢复（Postgres + UploadDir）

目标：具备最小可恢复能力。

## 媒体库（Media）

### 1) （增强）批量删除 / 批量查询接口

目标：为媒体库增加批量操作接口，并保持 owner 权限与“被引用禁止删除（409）”等约束一致。

### 2) （增强）增加 SHA256 去重策略

目标：支持（同用户/全站）去重策略，并定义重复上传处理策略（复用已有资源/拒绝/生成新版本）。

## 运维一致性

### 3) （增强）后台一致性扫描接口

目标：扫描上传目录与数据库记录的一致性（孤儿文件/孤儿记录），提供只读报告或修复选项。

## 前端

### 4) （前端待做）帖子编辑器的媒体选择/插入

目标：在编辑器中提供媒体库选择、上传、插入 Markdown 链接的体验。

## 远期模块草案（从架构文档迁移）

以下为“设计方向”，当前仓库结构中未落地对应目录/文件：

- Theme 主题系统：API、repo、service、中间件、前端动态主题组件。
- Plugin 插件系统：后端插件加载、hook/dispatcher、以及 pkg 级 SDK。
