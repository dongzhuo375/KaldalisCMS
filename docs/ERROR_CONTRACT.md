# API Error Contract

本文档是 KaldalisCMS 错误处理的单一事实来源（SSOT），覆盖后端语义、HTTP 映射、出站包络与前端消费约束。

## 1) 统一响应结构（冻结）

所有 API 错误必须返回：

```json
{
  "code": "NOT_FOUND",
  "message": "resource not found",
  "details": {
    "request_id": "8d9f6e8d-7b7d-45d0-8e74-f8f7e9e7fbb6"
  }
}
```

字段约定：

- `code`：稳定、机器可读错误码。
- `message`：对外可读文本，默认由错误码策略决定。
- `details`：附加上下文；字段始终存在（可为空对象）。
- `details.request_id`：请求链路追踪 ID（后端统一注入）。

代码来源：

- 错误码、默认映射、details 脱敏策略：`internal/core/error.go`
- 统一写出器与 request_id 注入：`internal/api/errorx/responses.go`
- 错误语义归一化：`internal/service/error_semantics.go`

## 2) 错误码与默认 HTTP 映射（冻结）

| code | 默认 HTTP | 默认 message |
| --- | --- | --- |
| `VALIDATION_FAILED` | `400` | `request validation failed` |
| `UNAUTHORIZED` | `401` | `unauthorized` |
| `FORBIDDEN` | `403` | `permission denied` |
| `NOT_FOUND` | `404` | `resource not found` |
| `DUPLICATE_RESOURCE` | `409` | `resource already exists` |
| `CONFLICT` | `409` | `request conflict` |
| `TIMEOUT` | `504` | `request timed out` |
| `INTERNAL_ERROR` | `500` | `internal server error` |

变更规则（冻结）：

- 新增/变更错误码必须同步更新本文档 + 测试并走评审。
- 禁止在 handler/middleware 私自定义“同义新码”。

## 3) 允许覆盖与禁止覆盖

允许覆盖：

- `message` 可在局部场景做业务友好化（如 `failed to build openapi spec`）。
- `status` 仅允许 `INTERNAL_ERROR` 等非业务语义码按场景调整。

禁止覆盖：

- 不允许修改既有业务语义码的 canonical HTTP（如 `NOT_FOUND` 必须是 `404`）。
- 不允许绕开统一写出器直接返回 `gin.H{"error": ...}` 或任意非契约结构。

## 4) details 脱敏与白名单

- 采用按错误码白名单过滤，未知字段默认不出站。
- 下列敏感键会被强制剔除：包含 `password`/`token`/`secret`/`authorization`。
- 推荐写入：`field`、`resource`、`references`、`request_id`。

## 5) 分层约束

- `handler`/`middleware`/`router` 仅使用 `internal/api/errorx` 出站。
- `service` 仅返回 `core.Err*` 语义错误（或可被 `errors.Is` 识别）。
- `repository` 原生错误需在 `service` 完成语义翻译后再向上返回。

## 6) 回归与 CI 阻断

当前契约测试：

- `internal/core/error_test.go`：`error -> code`、`code -> status/message`
- `internal/api/v1/responses_test.go`：统一写出器映射与 `request_id` 注入
- `internal/api/middleware/error_contract_integration_test.go`：中间件 4xx/5xx 包络一致性
- `internal/api/v1/write_error_contract_integration_test.go`：关键写接口错误包络

建议在 CI 中将上述测试作为阻断项，禁止跳过。
