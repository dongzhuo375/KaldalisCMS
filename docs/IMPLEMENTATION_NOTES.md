# 实现点（Implementation Notes）

本文记录**已落地且代码中可追溯**的关键实现点。每条尽量附带“对应文件/函数”，方便维护与排查。

## 错误平台化与观测闭环 - [2026-04-19 新增]

### 统一错误出口

- handler/middleware/router 全部通过 `internal/api/errorx/responses.go` 写出错误响应。
- 统一注入 `details.request_id`，并同步响应头 `X-Request-Id`，便于前端与日志联查。
- 出站前执行 `details` 白名单与敏感键剔除（`password/token/secret/authorization`）。

代表文件：
- `internal/core/error.go`
- `internal/api/errorx/responses.go`
- `internal/api/middleware/auth.go`
- `internal/api/middleware/casbin.go`
- `internal/router/swagger_routes_enabled.go`

### 统一请求上下文与恢复策略

- 新增 `RequestContext` 中间件：透传或生成 request_id 并写入 Gin context。
- 新增 `RecoverAsContract`：panic 统一转为错误契约包络，避免默认 recovery 的非契约输出。

代表文件：
- `internal/api/middleware/observability.go`
- `internal/router/router.go`

### 指标与结构化日志

- 新增 Prometheus 指标：
  - `kaldalis_http_requests_total{method,route,status,code}`
  - `kaldalis_http_request_duration_seconds{method,route}`
- 新增结构化访问日志字段：`request_id/method/route/status/duration_ms/code/ip`。

推荐告警基线（5 分钟窗口）：

- 5xx 比例异常：
  `sum(rate(kaldalis_http_requests_total{status=~"5.."}[5m])) / sum(rate(kaldalis_http_requests_total[5m])) > 0.05`
- 特定错误码激增（示例：`INTERNAL_ERROR`）：
  `sum(rate(kaldalis_http_requests_total{code="INTERNAL_ERROR"}[5m])) > 2`

代表文件：
- `internal/api/middleware/observability.go`
- `internal/router/router.go`

## 系统初始化与安装模式（Setup Mode）- [2026-03-06 新增]

### 自动模式切换与安装探测
系统启动时会进行“自感知”检查。如果满足以下任一条件，系统将自动进入 **SETUP MODE**（安装向导模式），仅暴露安装相关的 API：
- 数据库连接失败（配置缺失或错误）。
- 数据库连接成功，但 `system_settings` 表不存在或 `installed` 字段为 `false`。

代表文件：
- `cmd/server/main.go` (`BootstrapApp` 探测逻辑)

### 数据库自动拨备 (Auto-Provisioning)
在安装阶段，系统具备自动创建数据库的能力。用户只需提供 Postgres 实例的管理员账密：
1. 先连接到默认的 `postgres` 管理库。
2. 校验目标数据库名（正则校验防止注入）。
3. 如果目标库不存在，自动执行 `CREATE DATABASE`。
4. 随后再切换到目标库执行 `AutoMigrate`。

代表文件：
- `internal/service/setup_service.go` (`ValidateDatabase`)

### 数据库预检机制 (Pre-flight Check)
前端 Step 1 引入了强制预检。后端提供专门的 `check-db` 接口进行“探路”测试，确保地基稳固后才允许进入后续安装步骤。该接口为幂等操作，不修改任何物理配置。

代表文件：
- `internal/api/v1/setup.go` (`CheckDB`)
- `web/src/app/[locale]/setup/page.tsx` (前端测试按钮逻辑)

### 增量配置持久化
安装成功后，系统会更新 `config.yaml`。采用 **Patch（补丁）更新模式**：
- 仅修改 `database` 节点下的字段。
- 使用 `viper.WriteConfigAs` 确保即使初始状态无文件也能生成新文件。
- 保留文件中已有的其他配置（如 `jwt`, `media` 等），不执行覆写。

代表文件：
- `cmd/server/config.go` (`SaveDatabaseConfig`)

### Root Health/Ready 探针 - [2026-03-23 新增]

系统在根路径提供统一探针接口，供 Docker/K8s/反向代理进行存活与就绪判定。

- `GET /healthz`：仅表示进程存活，恒定返回 `200`。
- `GET /readyz`：表示服务可接流量；App Mode 下至少检查 DB `PingContext`；Setup Mode 下固定返回 `503`（未就绪）。

响应结构固定且可扩展：

```json
{
  "status": "ok|not_ready",
  "mode": "app|setup",
  "checks": {
    "database": {
      "status": "ok|fail|skip",
      "detail": "optional"
    }
  }
}
```

实现细节：
- `readyz` 的 DB 检查使用 `context.WithTimeout(..., 2s)` + `PingContext`，避免探针阻塞工作线程。
- `readyz` 增加进程内短 TTL 缓存以削峰：成功结果默认缓存 `400ms`，失败结果默认缓存 `250ms`，在高频探针下可显著降低 DB `ping` 压力并保持较快恢复感知。
- 不额外引入第三方健康检查库，优先保持依赖面最小；后续可在 `checks` 中继续扩展 `redis`、`queue`、`storage` 等依赖状态。

指标与告警：
- 新增 Prometheus 指标：
  - `kaldalis_probe_requests_total{probe,mode,result,cache}`：探针请求总量计数（包含缓存命中/未命中）。
  - `kaldalis_probe_ready_state{mode}`：当前就绪状态（`1=ready`，`0=not_ready`）。
- 指标抓取端点：`GET /metrics`。
- 推荐告警基线（按 5 分钟窗口）：
  - 就绪状态持续异常：`max_over_time(kaldalis_probe_ready_state{mode="app"}[5m]) < 1`。
  - 未就绪占比过高：
    `sum(rate(kaldalis_probe_requests_total{probe="readyz",mode="app",result="not_ready"}[5m])) / sum(rate(kaldalis_probe_requests_total{probe="readyz",mode="app"}[5m])) > 0.2`。

代表文件：
- `internal/api/v1/health.go`
- `internal/service/system_service.go` (`CheckDatabase`)
- `internal/router/router.go`（App/Setup 路由统一注册，含 `/metrics`）

容器探针示例：

```dockerfile
HEALTHCHECK --interval=10s --timeout=3s --start-period=15s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1
```

```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
readinessProbe:
  httpGet:
    path: /readyz
    port: 8080
```

---

## 媒体库（Media Library）

### 公共访问路径与物理存储

- 公共访问：`/media/a/{assetID}/{stored_name}`（静态目录映射到 `MEDIA_UPLOAD_DIR`，默认 `./data/uploads`）
- 物理路径：`{MEDIA_UPLOAD_DIR}/a/{assetID}/{stored_name}`

代表文件：
- `internal/router/router.go`（静态资源路由挂载）
- `internal/utils/env.go`（`MEDIA_*` 环境变量解析/默认值）

### 文件系统与数据库一致性（状态机 + 最终一致性）

采用“先落库、后写盘、成功后再置为可见”的状态机流程（见 `CreateAssetFromUpload`）：

- **PENDING (0)**：上传开始时先在 DB 创建 `media_assets` 记录，状态为 `PENDING`。
- **写入文件**：尝试把文件写入磁盘 `{upload_dir}/a/{assetID}/{stored_name}`。
- **UPLOADED (1)**：仅当文件写入成功后，才更新 DB 字段（`object_key/url/width/height`）并将状态更新为 `UPLOADED`。
- **FAILED (2)**：若创建目录/打开文件/写入失败，会把 DB 状态更新为 `FAILED`，并尽力删除部分写入的文件。

列表接口默认只返回 `UPLOADED` 资产（Repo `List` 中 `WHERE status = 1`），从而避免“DB 已有记录但文件未就绪”的资源被用户看到。

代表文件：
- `internal/service/media_service.go`（`CreateAssetFromUpload` 状态机与落盘逻辑）
- `internal/infra/repository/postgres/media_repo.go`（列表过滤 `status = UPLOADED`）

### 清理任务（Pending/软删除的最终一致性 GC）

存在后台 **GC 清理任务**（在路由初始化时启动，默认每小时执行一次 `CleanupStaleMedia`）：

- 清理“超过 1 小时仍为 PENDING 的记录”，执行物理删除（先删文件、再硬删 DB）。
- 删除采用“软删除优先、异步最终清理”：API 删除只设 `deleted_at`；GC 扫描“软删除超过 1 小时”的记录执行物理删除：
    1) 先删除文件（失败则返回，下一轮 GC 重试）
    2) 再对 DB 做硬删除（`Unscoped()`）

代表文件：
- `internal/router/router.go`（ticker/定时任务启动）
- `internal/service/media_service.go`（`CleanupStaleMedia` / 物理删除实现）

### 媒体引用同步（Best-Effort + 超时保护）

- Post Create/Update 会解析 Markdown 内容/封面 URL 并同步 `post_assets`（`PostService` 调用 `MediaService.SyncPostReferences`）。
- 引用同步是 **best-effort**：同步失败只记录日志，不影响发帖/更新成功（不会回滚 Post）。
- 具备 **超时保护**：
    - `CreatePost`：`context.WithTimeout(请求ctx, 10s)`
    - `UpdatePost`：`context.WithTimeout(context.Background(), 5s)`（独立于请求 ctx，避免请求取消导致同步完全跳过）

> 正则说明：用于从 Markdown/URL 中提取 assetID 的正则为 `reAssetURL = /media/a/(\d+)/[^)\s]+`，只要 URL 路径中包含 `/media/a/{id}/...`（即使有 CDN 域名）即可提取 ID；若未来调整 URL 结构，需要同步更新该正则。

代表文件：
- `internal/service/post_service.go`（发帖/更新 -> 同步引用）
- `internal/service/media_service.go`（`SyncPostReferences` 与引用写入逻辑）

---

## 安全机制（Auth & Security）

### CSRF 保护与前后端对接
本系统采用 **Stateless CSRF**（Double Submit Cookie）与指纹绑定相结合的方案。

#### 后端行为：
1. **登录成功时 (EstablishSession)**：
    - 生成 uuid 作为 `kaldalis_csrf` Cookie（HTTPOnly=false，前端可读）。
    - 同时将该 uuid 计算 HMAC 哈希后，存入 JWT Payload (`csrf_h`)。
    - 这实现了 CSRF Token 与当前登录 Session 的强绑定。

2. **请求校验 (Middleware: CSRFCheck)**：
    - 仅针对受保护的 **写操作接口**（POST/PUT/DELETE /posts, /media...）生效。
    - 读取 Header `X-CSRF-Token`。
    - 读取 Cookie `kaldalis_csrf`。
    - 校验 1：Header 值必须等于 Cookie 值（防跨域伪造）。
    - 校验 2：Header 值计算哈希后必须等于 JWT 中的 `csrf_h`（防 Session 劫持）。

#### 前端对接指南：
所有写操作请求必须携带 `X-CSRF-Token` 头，值需从 Cookie `kaldalis_csrf` 中读取。
参考实现：`web/src/lib/api.ts` (Axios Interceptor 自动注入)。

代表文件：
- `internal/infra/auth/session.go` (`EstablishSession/ValidateCSRF`)
- `internal/api/middleware/auth.go` (`CSRFCheck`)
- `web/src/lib/api.ts` ( 前端请求拦截器)

---

## 媒体库安全策略 - [2026-03-09 新增]

### 目录暴露收敛 (Directory Exposure Convergence)
为防止 `uploadDir` 根目录下的敏感文件（如备份、日志等）被误暴露，我们在路由层实际上仅公开了 `uploadDir/a` (assets) 子目录。

- **实施文件**: `internal/router/router.go`
- **原逻辑**: `r.Static("/media", uploadDir)` -> 暴露整个目录。
- **现逻辑**: `r.Static("/media/a", filepath.Join(uploadDir, "a"))` -> 仅暴露媒体资产目录。

这意味着所有上传的媒体文件 URL 均形如 `/media/a/{id}/{filename}`。若未来需新增公开目录（如 `avatars/`），需显式在 Router 中注册。

---

## API 文档分层与构建策略 - [2026-03-13 新增]

### 分层职责

- 注释层（接口契约源）：`internal/api/v1/*.go` 与 `internal/api/v1/dto/*.go`
  - 在 Handler 上维护 `@Summary/@Router/@Param/@Success/@Failure`。
  - DTO 作为 API contract，供 Swag 扫描生成 schema。
- 产物层（生成文档）：`internal/docs/`
  - `docs.go`（注册 swagger 文档）
  - `swagger.json` / `swagger.yaml`（生成产物）
- 路由集成层（文档胶水层）：`internal/router/`
  - `swagger_options.go`：文档开关与元信息配置对象
  - `swagger_routes_enabled.go`：`//go:build swagger`，挂载 UI 与 OpenAPI3 JSON
  - `swagger_routes_disabled.go`：`//go:build !swagger`，空实现（生产可剔除）
- 应用装配层：`cmd/server/main.go` + `cmd/server/config.go`
  - 从配置注入 `SwaggerOptions`，并在 `NewAppRouter/NewSetupRouter` 统一注册。

### 运行时开关（config + env）

配置源：`cmd/configs/config.yaml`

```yaml
swagger:
  enabled: false
  path: /swagger
  title: KaldalisCMS API
  version: dev
  description: KaldalisCMS backend API documentation
```

支持环境变量覆盖（Viper 自动映射）：

- `SWAGGER_ENABLED`
- `SWAGGER_PATH`
- `SWAGGER_TITLE`
- `SWAGGER_VERSION`
- `SWAGGER_DESCRIPTION`

### 编译标签策略（开发/生产分离）

- 开发/测试文档能力：使用 `-tags swagger` 构建，包含 `gin-swagger/swaggo` 集成与文档路由。
- 生产最小二进制：默认不带 `swagger` 标签构建，Swagger 相关代码不参与编译与链接。

### 对外文档路由

- Swagger UI：`<swagger.path>`（默认 `/swagger`）
- OpenAPI 3 JSON：`<swagger.path>-openapi3.json`（默认 `/swagger-openapi3.json`）

### 文档刷新命令

```powershell
Set-Location -Path D:\project\KaldalisCMS
go generate -tags swagger ./internal/docs
```
