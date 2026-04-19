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

### P0-2) 生产级 Session/Cookie/CSRF 安全配置收口   **[2026-03-09 完成]**

目标：确保登录态在生产环境符合最小安全要求。

范围（代表文件）：
- `internal/infra/auth/session.go`
- `pkg/security/csrf.go`
- `internal/api/middleware/auth.go`

完成标准（DoD）：
- **[2026-03-09 完成]** 已将 Cookie 策略升级为 Strict/Lax，添加了弱密钥检测，并完善了文档。
- Cookie: `HttpOnly`、`Secure`（https）、`SameSite` 策略明确；
- 密钥/盐来自环境变量或配置，禁止默认弱值；
- CSRF 校验在受保护写接口一致生效，并有前端对接说明。

### P0-4) `/media` 静态目录暴露收敛 **[2026-03-09 完成]**

目标：避免 uploadDir 下非媒体文件被意外暴露。

范围（代表文件）：
- `internal/router/router.go`（`r.Static("/media", uploadDir)`）
- `cmd/server/config.go`（media.upload_dir 配置）

完成标准（DoD）：
- 只暴露必要的子路径（例如仅 `/media/a` 或等价安全策略）；
- 文档说明 uploadDir 的目录约束；
- 覆盖常见路径穿越/误放文件风险。

### P0-5) Health/Ready 探活接口 **[2026-03-23 完成]**

目标：支持部署平台（容器/反向代理）探活与就绪判断。

范围（代表文件）：
- `internal/router/router.go`（注册新路由）
- `internal/api/v1/system.go` 或新增 `internal/api/v1/health.go`

完成标准（DoD）：
- `/healthz`：进程存活即可 200；
- `/readyz`：至少检查 DB 连接可用；
- 响应结构固定，便于监控。

本次改造进度：
- [x] 新增根路径探针：`GET /healthz` 与 `GET /readyz`。
- [x] App Mode `readyz` 接入 DB `PingContext` 检查，失败返回 `503`。
- [x] Setup Mode `readyz` 语义固定为未就绪（`503`）。
- [x] 响应结构统一为可扩展 `status/mode/checks`。
- [x] 提供 Docker `HEALTHCHECK` 与 K8s `liveness/readinessProbe` 对接示例（见 `docs/IMPLEMENTATION_NOTES.md`）。

新增代做（探针能力扩展）：
- [ ] P1：在 `checks` 中增加 `redis`、`queue`、`storage` 多依赖聚合检查。
- [x] P1：加入轻量级短 TTL 缓存（如 200~500ms）以削峰高频探针流量。
- [x] P1：补充探针结果指标上报（Prometheus counter/gauge）与告警阈值文档。
- [ ] P1：补充 `Dockerfile`/`docker-compose` 实际探针配置，并完成一次容器内探针联调记录。
- [ ] P1：补充 K8s `startupProbe`（冷启动保护）及 `liveness/readiness` 推荐阈值基线文档。

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

本次改造进度（仅后端）：
- [x] `internal/core/error.go`：新增稳定错误码常量与 `ErrorCodeOf` 映射。
- [x] `internal/api/v1/dto/common_dto.go`：统一错误 DTO 为 `{code,message,details}`。
- [x] `internal/api/v1/responses.go`：新增统一错误/成功响应写入器。
- [x] `internal/api/v1/post.go`：统一公共文章接口错误返回格式。
- [x] `internal/api/v1/admin_post.go`：统一错误返回并统一 Update/Delete 成功返回为 `200 + message`。
- [x] `internal/api/v1/user.go`：统一鉴权/注册/登录错误返回格式与 code。
- [x] `internal/api/v1/system.go`：统一 setup/status 错误映射。
- [x] `internal/api/v1/setup.go`：统一安装流程错误返回格式。
- [x] `internal/api/v1/media.go`：统一上传/删除/列表错误返回格式。
- [x] OpenAPI（Swagger 注释与 `internal/docs/docs.go`）全量同步更新（已补齐 `media` 相关路由注解并重新生成文档产物）。

新增代做（P0-7 收尾）：
- [x] 为 `internal/api/v1/media.go` 补齐 Swagger 注解（`@Summary/@Success/@Failure/@Router`）。
- [x] 重新生成并校验 `internal/docs/docs.go`、`internal/docs/swagger.json`、`internal/docs/swagger.yaml`，确认包含 `/media`、`/media/{id}`、`/posts/{id}/media`。
- [x] 增加最小错误响应合约测试（覆盖 `400/403/404/409` 的 `{code,message,details}` 结构）。

后端剩余 error 处理清单（施工前盘点）：
- [x] P0：统一中间件错误返回结构与 code（`internal/api/middleware/auth.go`、`internal/api/middleware/casbin.go`，替换 `gin.H{"error":...}` 为统一 `ErrorResponse`）。
- [x] P0：统一路由层非业务错误返回结构（`internal/router/swagger_routes_enabled.go` 的 OpenAPI 构建失败返回）。
- [x] P0：收敛 service 层主流程错误到 `core` 语义（`internal/service/system_service.go` 的安装冲突、`internal/service/post_service.go` 的 validation/conflict 分支、`internal/service/media_service.go` 的 not found/conflict/validation 映射）。
- [x] P1：补齐错误映射单测（`error -> code -> status`），覆盖 `core.ErrorCodeOf` 与 `respondErrorByCore` 关键分支。
- [x] P1：增加中间件与关键写接口的集成测试，校验所有 4xx/5xx 响应均为 `{code,message,details}`。

新增代做（service 层收敛拆分）：
- [x] 将 `internal/service/system_service.go` 的安装冲突收敛为 `core.ErrConflict` 语义，并由 API 统一映射。
- [x] 将 `internal/service/post_service.go` 中关键 `errors.New/fmt.Errorf` 分支替换为可识别的 `core` 语义错误（validation/conflict）。
- [x] 定义并落地 `internal/service/media_service.go` 业务错误到 `core.ErrorCode` 的映射策略（保持 409/400 语义稳定）。

### P0-8) 双轨日志体系（Operation Log + Audit Event）

目标：建设“操作日志（运营/客服排障）+ 审计事件（安全/合规追踪）”双轨体系，并采用“service 层统一采集 + outbox 异步投递 + 关键事件同步兜底”模式。

意义（为什么现在做）：
- 将排障从“看零散日志猜问题”升级为“按 request_id / actor / action 直接追踪链路”。
- 将安全追责从“事后人工拼接”升级为“结构化事件可检索、可导出、可告警”。
- 为后续用户管理、内容治理、合规报表提供统一数据底座。

现状收益（本轮已完成基础能力）：
- [x] `request_id` 注入与透传：`internal/api/middleware/observability.go`、`internal/api/errorx/responses.go`
- [x] 统一 panic 收敛出口：`RecoverAsContract`
- [x] 结构化访问日志字段：`request_id/method/route/status/duration_ms/code/ip`
- [x] HTTP 指标：`kaldalis_http_requests_total`、`kaldalis_http_request_duration_seconds`
- [x] 统一错误契约与 details 脱敏基线：`docs/ERROR_CONTRACT.md`

当前差距（与工业目标对比）：
- [ ] P0：缺统一 `AuditEvent` 事件模型（`event_id/occurred_at/actor/action/resource/result/risk_level/...`）。
- [ ] P0：缺审计落库与检索能力（当前以访问日志为主，不等于审计轨）。
- [ ] P0：缺 outbox 表/repo/worker，尚未实现“业务写入 + 事件写入”同事务原子性。
- [ ] P0：缺关键事件同步兜底（异步失败时关键事件仍需同步落库）。
- [ ] P1：缺审计查询 API、导出能力、告警基线与值班手册。
- [ ] P2：缺不可篡改增强（append-only 约束、哈希链/签名校验）。

范围（代表文件）：
- `internal/core/`（新增 `audit.go`：事件模型、风险等级、采集接口）
- `internal/service/`（`post/media/user/setup` 统一采集埋点）
- `internal/api/middleware/`（请求上下文字段与 actor 提取补强）
- `internal/infra/model/`、`internal/infra/repository/postgres/`（审计表 + outbox 表 + repo）
- `internal/router/router.go`（outbox relay worker 生命周期挂载）
- `docs/IMPLEMENTATION_NOTES.md`、`docs/ERROR_CONTRACT.md`（规则与运维手册）

分阶段施工与 DoD：

阶段 1（P0，打地基）
- [ ] 定义 `AuditEvent` 统一模型并冻结字段规范（含脱敏规则）。
- [ ] 落地本地审计表与 outbox 表（最小可用）。
- [ ] 核心写操作（user/post/media/setup）都能产出结构化事件并带 `request_id`。
- [ ] DoD：单次请求可通过 `request_id` 串联 API 错误、访问日志、审计事件。

阶段 2（P1，一致性与覆盖）
- [ ] 在 `post/media/user/setup` service 层完成 P0 事件清单全覆盖（成功/失败都记录）。
- [ ] Outbox relay 支持重试、退避、失败计数；关键事件同步兜底。
- [ ] 失败分类统一对齐 `core.ErrorCode`，沉淀事件 `result/error_code`。
- [ ] DoD：P0 清单覆盖率 >= 95%，写链路性能增量可控（P95 增量 < 5%）。

阶段 3（P2，治理与合规）
- [ ] 审计查询 API（仅 `super_admin`/安全角色）+ CSV/JSON 导出。
- [ ] 留存与归档策略（冷热分层）+ 告警规则（高风险行为/登录失败爆发/越权尝试）。
- [ ] append-only 与篡改检测（哈希链/签名批校验至少一种）。
- [ ] DoD：可按用户/资源/时间/动作追溯，并可生成审计报表。

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

## 后端待实现（P1：前端已占位）

### P1-6) 用户管理 API

目标：提供用户 CRUD 及角色分配接口，供前端用户管理页面对接。

范围（代表文件）：
- `internal/api/v1/admin_user.go`（新增管理接口）
- `internal/service/user_service.go`（扩展业务逻辑）
- `internal/infra/repository/postgres/user_repo.go`（扩展查询）

完成标准（DoD）：
- [ ] 用户列表接口（分页、搜索）：`GET /api/v1/admin/users`
- [ ] 用户详情：`GET /api/v1/admin/users/:id`
- [ ] 创建用户：`POST /api/v1/admin/users`
- [ ] 更新用户信息/角色：`PUT /api/v1/admin/users/:id`
- [ ] 删除/禁用用户：`DELETE /api/v1/admin/users/:id`
- [ ] Casbin 权限策略同步（角色变更时自动更新策略）
- [ ] 最小测试覆盖（CRUD + 权限校验）

### P1-7) 数据统计/分析 API

目标：提供站点基础统计数据，供前端 Analytics 页面对接。

范围（代表文件）：
- `internal/api/v1/analytics.go`（新增）
- `internal/service/analytics_service.go`（新增）
- `internal/infra/model/page_view.go`（新增，访问记录表）

完成标准（DoD）：
- [ ] 页面浏览量追踪中间件或埋点接口
- [ ] 统计概览接口：`GET /api/v1/admin/analytics/overview`（总浏览量、独立访客、热门文章）
- [ ] 时间序列数据：`GET /api/v1/admin/analytics/traffic?range=7d`（按日/小时聚合）
- [ ] 来源统计：`GET /api/v1/admin/analytics/sources`（Referrer 聚合）
- [ ] 最小测试覆盖

## 远期模块草案（从架构文档迁移）

以下为“设计方向”，当前仓库结构中未落地对应目录/文件：

- Theme 主题系统：API、repo、service、中间件、前端动态主题组件。
- Plugin 插件系统：后端插件加载、hook/dispatcher、以及 pkg 级 SDK。
