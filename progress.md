## Project Progress (项目进度)

As of 2025年12月3日星期三:

### Completed Tasks (已完成任务):
-   **Models Defined (模型定义):** Core models for `User`, `Category`, `Tag`, and `Post` are defined in `internal/model`. (核心模型，包括 `User`、`Category`、`Tag` 和 `Post` 已在 `internal/model` 中定义。)
-   **GORM and PostgreSQL Setup (GORM 和 PostgreSQL 设置):** GORM with the PostgreSQL driver is configured. `internal/repository/db.go` initializes a PostgreSQL database connection using GORM. (GORM 及其 PostgreSQL 驱动已配置。`internal/repository/db.go` 文件已成功使用 GORM 初始化 PostgreSQL 数据库连接。)
-   **Database Dependency Injection:** Global database connection refactored to use dependency injection, improving testability and clarity. (全局数据库连接已重构为使用依赖注入，提高了可测试性和清晰度。)
-   **PostRepository Refactoring for Entity Interface (Post仓库层接口实体化重构):**
    -   `PostRepository` (in `internal/repository/postgres/post_repo.go`) has been refactored to use `entity.Post` for its public interface. (Post仓库层已重构，其公共接口现在使用 `entity.Post`。)
    -   Mapper functions (`toEntity`, `toModel`) have been added to `post_repo.go` for internal conversion between `model.Post` and `entity.Post`. (已在 `post_repo.go` 中添加映射函数 `toEntity` 和 `toModel`，用于 `model.Post` 和 `entity.Post` 之间的内部转换。)
    -   All CRUD methods (`GetAll`, `GetByID`, `Create`, `Update`) in `post_repo.go` now accept/return `entity.Post`. (Post仓库层中的所有CRUD方法现在接受/返回 `entity.Post`。)
-   **Enhanced Error Handling in Repository Layer:** Custom error types (`core.ErrNotFound`, `core.ErrDuplicate`, `core.ErrInternalError`) are defined and used. Error wrapping with `fmt.Errorf("%w", err)` is implemented in `post_repo.go` to preserve error context. Postgres-specific unique constraint error detection is in place. (`post_repo.go` 中实现了增强的仓库层错误处理：定义并使用了自定义错误类型，通过 `fmt.Errorf("%w", err)` 实现错误包装以保留上下文，并支持 PostgreSQL 特有的唯一约束错误检测。)
-   **PostService Implementation (Post服务层实现):** The `PostService` in `internal/service/post_service.go` has been implemented, integrating with the `PostRepository` to provide business logic for post-related operations. (`internal/service/post_service.go` 中的 `PostService` 已实现，与 `PostRepository` 集成，为帖子相关操作提供业务逻辑。)
-   **Post API Handlers (Post API 处理器):** API handlers for posts are implemented in `internal/api/v1/post.go`, handling requests and responses, and interacting with the `PostService`. (`internal/api/v1/post.go` 中已实现帖子相关的 API 处理器，负责处理请求和响应，并与 `PostService` 交互。)
-   **Router Setup (路由设置):** The application router in `internal/router/router.go` has been configured to define API routes for the post endpoints. (`internal/router/router.go` 中的应用程序路由已配置，用于定义帖子端点的 API 路由。)
-   **User Module Scaffolding (用户模块脚手架):**
    -   **Entity and Model (实体与模型):** `User` entity (`internal/core/entity/user.go`) and GORM model (`internal/infra/model/user.go`) are defined.
    -   **Repository (仓库层):** `UserRepository` (`internal/infra/repository/postgres/user_repo.go`) is fully implemented with CRUD operations, error handling, and entity-model mapping.
    -   **API Layer (API 层):** Basic structure for `UserAPI` (`internal/api/v1/user.go`) is in place with `Register` and `Login` endpoints defined.
-   **Database Model Optimization (数据库模型优化):**
    -   **Handled empty strings (处理了空字符串):** Added `check` constraints to model fields to prevent empty or whitespace-only strings.
    -   **Optimized JSON output (优化了JSON输出):** Used `json` tags (`-`, `omitempty`) to control the JSON output for frontend, preventing oversized payloads and circular references.
    -   **Added relationships (添加了对应关系):** Defined `one-to-many` and `many-to-many` relationships between models.


### To-Do (待办事项):
-   **Implement `UserService` Business Logic (实现 `UserService` 业务逻辑):**
    -   Implement password hashing using `bcrypt` before creating a user.
    -   Implement user creation and verification logic in `internal/service/user_service.go`.
-   **Complete `UserAPI` Integration (完成 `UserAPI` 集成):**
    -   Integrate the `UserAPI` with the completed `UserService`.
    -   Implement JWT-based authentication in the `Login` handler.
-   **Implement `Category` and `Tag` Modules (实现 `Category` 和 `Tag` 模块):**
    -   Follow the same layered architecture pattern (Entity, Repository, Service, API) to implement the `Category` and `Tag` modules.
-   **Refine Configuration Management (优化配置管理):**
    -   Load configuration from `config.yaml` files.
-   **Implement Database Migrations (实现数据库迁移):**
    -   Set up a proper database migration system.