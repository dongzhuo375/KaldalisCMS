# API Error Contract

本文件定义 KaldalisCMS 后端统一错误响应契约，作为前后端联调与排障的单一事实来源。

## 1) 统一响应结构

所有 API 错误响应统一为：

```json
{
  "code": "NOT_FOUND",
  "message": "resource not found",
  "details": null
}
```

字段约定：

- `code`：稳定、机器可读错误码。
- `message`：面向客户端的可读错误信息。
- `details`：附加信息；可为对象或 `null`，字段始终存在。

代码来源：

- 错误码定义与归一：`internal/core/error.go`
- HTTP 映射与默认 message：`internal/api/v1/responses.go`
- DTO 结构：`internal/api/v1/dto/common_dto.go`

## 2) 标准错误码与默认 HTTP 映射

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

说明：

- 映射函数为 `respondErrorByCore`。
- 业务接口可在少量场景覆盖默认 message 或 HTTP（例如上传超限返回 `413`）。
- 覆盖时仍必须保持 `{code,message,details}` 结构。

## 3) `details` 使用规则

- `details` 用于附加可诊断上下文，不承载主语义。
- 建议内容：参数名、失败原因、冲突数量、追踪辅助信息。
- 禁止内容：明文敏感信息（密码、token、密钥、连接串）。

示例：

```json
{
  "code": "VALIDATION_FAILED",
  "message": "invalid request body",
  "details": {
    "reason": "Key: 'UserRegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
  }
}
```

```json
{
  "code": "CONFLICT",
  "message": "asset is referenced",
  "details": {
    "references": 3
  }
}
```

## 4) 层级约束

- handler / middleware / router 层出现错误时，均必须返回统一错误包络。
- service 层优先返回 `core.Err*`（或可被 `errors.Is` 识别到 `core.Err*` 的包装错误）。
- repository 层自定义错误应在 service 层收敛为 `core` 语义后再出站。

## 5) 回归检查

当前最小回归覆盖：

- `internal/core/error_test.go`：`error -> code`
- `internal/api/v1/responses_test.go`：`error -> code -> status/message/details`
- `internal/api/middleware/error_contract_integration_test.go`：中间件 4xx/5xx 包络
- `internal/api/v1/write_error_contract_integration_test.go`：关键写接口错误包络

建议在变更错误码或 message 前，先更新本文件并同步测试。
