# KaldalisCMS

This is Readme.
这是 Readme。

## API Documentation (Swagger + OpenAPI 3)

Swagger integration is layered in `internal/router` and controlled by **both** runtime config and compile-time build tags.

- Compile-time gate: `swagger` build tag
- Runtime gate: `swagger.enabled` in `cmd/configs/config.yaml` (or `SWAGGER_ENABLED` env)
- UI path: `swagger.path` (default `/swagger`)
- OpenAPI 3 JSON path: `<swagger.path>-openapi3.json` (default `/swagger-openapi3.json`)

### Config keys

```yaml
swagger:
  enabled: false
  path: /swagger
  title: KaldalisCMS API
  version: dev
  description: KaldalisCMS backend API documentation
```

Env override examples:

- `SWAGGER_ENABLED=true`
- `SWAGGER_PATH=/swagger`
- `SWAGGER_TITLE=KaldalisCMS API`
- `SWAGGER_VERSION=v1`
- `SWAGGER_DESCRIPTION=Internal API docs`

### Build strategy

- Production build (exclude swagger code/deps from binary): build **without** `swagger` tag.
- Development/doc build (include swagger routes): build **with** `-tags swagger`.

### Regenerate docs

```powershell
Set-Location -Path D:\project\KaldalisCMS
go generate -tags swagger ./internal/docs
```
