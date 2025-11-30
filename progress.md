## Project Progress (项目进度)

As of 2025年11月30日星期日:

### Completed Tasks (已完成任务):
-   **Models Defined (模型定义):** Core models for `User`, `Category`, `Tag`, and `Post` are defined in `internal/model`. (核心模型，包括 `User`、`Category`、`Tag` 和 `Post` 已在 `internal/model` 中定义。)
-   **GORM and PostgreSQL Setup (GORM 和 PostgreSQL 设置):** GORM with the PostgreSQL driver is configured. `internal/repository/db.go` initializes a PostgreSQL database connection using GORM. (GORM 及其 PostgreSQL 驱动已配置。`internal/repository/db.go` 文件已成功使用 GORM 初始化 PostgreSQL 数据库连接。)
-   **Database Dependency Injection:** Global database connection refactored to use dependency injection, improving testability and clarity. (全局数据库连接已重构为使用依赖注入，提高了可测试性和清晰度。)
-   **PostRepository Refactoring for Entity Interface (Post仓库层接口实体化重构):**
    -   `PostRepository` (in `internal/repository/postgres/post_repo.go`) has been refactored to use `entity.Post` for its public interface. (Post仓库层已重构，其公共接口现在使用 `entity.Post`。)
    -   Mapper functions (`toEntity`, `toModel`) have been added to `post_repo.go` for internal conversion between `model.Post` and `entity.Post`. (已在 `post_repo.go` 中添加映射函数 `toEntity` 和 `toModel`，用于 `model.Post` 和 `entity.Post` 之间的内部转换。)
    -   All CRUD methods (`GetAll`, `GetByID`, `Create`, `Update`) in `post_repo.go` now accept/return `entity.Post`. (Post仓库层中的所有CRUD方法现在接受/返回 `entity.Post`。)
-   **Enhanced Error Handling in Repository Layer:** Custom error types (`core.ErrNotFound`, `core.ErrDuplicate`, `core.ErrInternalError`) are defined and used. Error wrapping with `fmt.Errorf("%w", err)` is implemented in `post_repo.go` to preserve error context. Postgres-specific unique constraint error detection is in place. (`post_repo.go` 中实现了增强的仓库层错误处理：定义并使用了自定义错误类型，通过 `fmt.Errorf("%w", err)` 实现错误包装以保留上下文，并支持 PostgreSQL 特有的唯一约束错误检测。)


### Current Discussions & Decisions (当前讨论与决策):
-   **Repository `Create` Method Signature:** Decided that `PostRepository.Create` should return `(int, error)` to provide the ID of the newly created post, enabling the service layer to fetch the complete entity if needed. (已决定 `PostRepository.Create` 方法应返回 `(int, error)` 以提供新创建文章的 ID，从而使服务层在需要时可以获取完整的实体。)
-   **Service Layer Adaptation:** Adapting `PostService.CreatePost` and `PostService.UpdatePost` to the new repository method signatures and implementing the "command-then-query" pattern in the service layer. (正在调整 `PostService.CreatePost` 和 `PostService.UpdatePost` 以适应新的仓库方法签名，并在服务层实现“命令-查询”模式。)

### Next Steps (下一步计划):
-   **Finalize Repository Interface and Implementation for `Create`:** Complete the change of `PostRepository.Create` to `(int, error)` in `core/repository.go` and `internal/repository/postgres/post_repo.go`. (完成 `core/repository.go` 和 `internal/repository/postgres/post_repo.go` 中 `PostRepository.Create` 接口和实现的修改，使其返回 `(int, error)`。)
-   **Update Service Layer (`post_service.go`):** Implement the "command-then-query" pattern for `CreatePost` and adapt `UpdatePost` to the new `PostRepository.Update` signature. (更新服务层 (`post_service.go`)：为 `CreatePost` 实现“命令-查询”模式，并调整 `UpdatePost` 以适应新的 `PostRepository.Update` 签名。)
-   **Update API Layer (`api/v1/post.go`)**: Refine error handling in API handlers to explicitly check for specific service/repository errors (e.g., `ErrNotFound`) and return appropriate HTTP status codes (e.g., 404 vs 500). (更新API层 (`api/v1/post.go`)：优化API处理程序中的错误处理，明确检查特定的服务/仓库错误（例如 `ErrNotFound`）并返回适当的HTTP状态码（例如404 vs 500）。)
-   **Context Propagation:** Implement `context.Context` passing throughout the application's call stack (API -> Service -> Repository). (在应用的调用栈中（API -> Service -> Repository）实现 `context.Context` 的传递。)
-   **Safe Database Migrations:** Replace GORM's `AutoMigrate` with a dedicated database migration tool. (用专门的数据库迁移工具替换 GORM 的 `AutoMigrate`。)
-   **Testing:** Begin implementing unit tests, starting with the service layer. (开始实施单元测试，从服务层开始。)
