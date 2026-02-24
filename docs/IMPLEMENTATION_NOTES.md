# 实现点（Implementation Notes）

本文记录**已落地且代码中可追溯**的关键实现点。每条尽量附带“对应文件/函数”，方便维护与排查。

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

