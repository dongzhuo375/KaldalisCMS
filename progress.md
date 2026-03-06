## Project Progress (项目进度)

### [2026-03-06 星期五] - 健壮初始化与预检系统

-   **Robust System Setup (健壮的安装模式):** 实现了自动拨备数据库逻辑，系统现在能自动创建缺失的 PostgreSQL 数据库。
-   **Pre-flight Check (安装预检):** 在前端 Step 1 增加了“Test Connection”按钮，后端提供专门的预检接口，确保数据库地基稳固后才进行安装。
-   **System Self-Awareness (启动自感知):** 强化了 `main.go` 启动逻辑，能够自动检测数据库连接状态和 `installed` 标记，不符合条件时强制降级到 Setup Mode。
-   **Incremental Config Persistence (增量配置持久化):** 优化了 `SaveDatabaseConfig`，更新数据库配置时不再覆写整个文件，而是保留原有的 `jwt`, `media` 等配置。

---

### [2025-12-03 星期三] - 核心模型与 Post 模块

-   **Models Defined (模型定义):** Core models for `User`, `Category`, `Tag`, and `Post` are defined in `internal/model`.
-   **GORM and PostgreSQL Setup (GORM 和 PostgreSQL 设置):** GORM with the PostgreSQL driver is configured. `internal/repository/db.go` initializes a PostgreSQL database connection using GORM.
-   **Database Dependency Injection:** Global database connection refactored to use dependency injection.
-   **PostRepository Refactoring (Post仓库层重构):** Refactored to use `entity.Post` for its public interface with internal mappers (`toEntity`, `toModel`).
-   **Enhanced Error Handling:** Custom error types (`core.ErrNotFound`, `core.ErrDuplicate`, `core.ErrInternalError`) are defined and used in `post_repo.go`.
-   **PostService & API Implementation:** `PostService` and API handlers for posts are fully implemented and integrated.
-   **Router Setup (路由设置):** Configured for post endpoints.
-   **User Module Scaffolding (用户模块脚手架):** Entity, Repository, and basic API (`Register`, `Login`) for Users are in place.
-   **Database Model Optimization:** Handled empty strings and optimized JSON output with `omitempty` tags.
