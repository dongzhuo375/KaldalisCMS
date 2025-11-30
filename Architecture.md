/KaldalisCMS
├── 应用入口
│   ├── cmd/
│   │   └── server/
│   │       ├── main.go                    # 应用主入口
│   │       └── config.go                  # 配置初始化
│   ├── configs/
│   │   ├── config.yaml                    # 主配置文件
│   │   ├── config.prod.yaml              # 生产环境配置
│   │   └── config.dev.yaml               # 开发环境配置
│   └── Makefile                          # 构建脚本
│
├── 后端核心
│   ├── internal/
│   │   ├── api/
│   │   │   ├── v1/
│   │   │   │   ├── theme.go              # 主题管理API
│   │   │   │   ├── post.go               # 文章管理API
│   │   │   │   ├── user.go               # 用户管理API
│   │   │   │   ├── upload.go             # 文件上传API
│   │   │   │   ├── system.go             # 系统管理API
│   │   │   │   └── plugin.go             # 插件管理API
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go               # 认证中间件
│   │   │   │   ├── cors.go               # 跨域中间件
│   │   │   │   ├── logger.go             # 日志中间件
│   │   │   │   ├── theme.go              # 主题上下文中间件
│   │   │   │   └── plugin.go             # 插件中间件
│   │   │   └── router.go                 # 路由配置
│   │   │
│   │   ├── core/
│   │   │   ├── entity.go                 # 领域实体定义
│   │   │   ├── service.go                # 服务接口
│   │   │   ├── repository.go             # 仓储接口
│   │   │   └── plugin.go                 # 插件核心接口定义
│   │   │
│   │   ├── service/
│   │   │   ├── theme_service.go          # 主题业务服务
│   │   │   ├── post_service.go           # 文章业务服务
│   │   │   ├── user_service.go           # 用户业务服务
│   │   │   ├── file_service.go           # 文件服务
│   │   │   ├── system_service.go         # 系统服务
│   │   │   └── plugin_service.go         # 插件生命周期管理服务
│   │   │
│   │   ├── repository/
│   │   │   ├── postgres/
│   │   │   │   ├── theme_repo.go         # 主题数据访问
│   │   │   │   ├── post_repo.go          # 文章数据访问
│   │   │   │   ├── user_repo.go          # 用户数据访问
│   │   │   │   └── migration/
│   │   │   │       └── 001_init.sql      # 初始迁移
│   │   │   └── redis/
│   │   │       └── cache.go              # 缓存服务
│   │   │
│   │   ├── model/
│   │   │   ├── theme.go                  # 主题模型
│   │   │   ├── post.go                   # 文章模型
│   │   │   ├── user.go                   # 用户模型
│   │   │   ├── category.go               # 分类模型
│   │   │   ├── setting.go                # 系统设置模型
│   │   │   └── plugin.go                 # 插件元数据模型
│   │   │
│   │   ├── theme/
│   │   │   ├── manager.go                # 主题管理器
│   │   │   ├── loader.go                 # 主题加载器
│   │   │   ├── registry.go               # 主题注册表
│   │   │   ├── validator.go              # 主题验证器
│   │   │   └── cache.go                  # 主题缓存
│   │   │
│   │   ├── plugin/                       # 插件系统核心
│   │   │   ├── manager.go                # 插件管理器
│   │   │   ├── registry.go               # 插件注册表
│   │   │   ├── loader.go                 # 插件加载器
│   │   │   ├── dispatcher.go             # 事件分发器
│   │   │   ├── service.go                # 插件RPC服务
│   │   │   └── hook/                     # 钩子系统
│   │   │       ├── hook.go               # 钩子定义
│   │   │       ├── manager.go            # 钩子管理器
│   │   │       └── types.go              # 钩子类型
│   │   │
│   │   └── utils/
│   │       ├── crypto.go                 # 加密工具
│   │       ├── file.go                   # 文件工具
│   │       └── validator.go              # 验证工具
│   │
│   ├── pkg/
│   │   ├── plugin/                       # 插件SDK
│   │   │   ├── sdk.go                    # 插件SDK主入口
│   │   │   ├── interfaces.go             # 插件接口定义
│   │   │   ├── types.go                  # 共享类型定义
│   │   │   └── util.go                   # SDK工具函数
│   │   │
│   │   └── database/                     # 数据库工具包（如有）
│   │       ├── mysql.go
│   │       └── redis.go
│   │
│   └── plugins/                          # 插件目录
│       ├── example/                      # 示例插件
│       │   ├── main.go                   # 插件主文件
│       │   ├── go.mod                    # 插件模块定义
│       │   └── implementation.go         # 插件实现
│       └── .gitkeep                      # 保持目录结构
│
├── 前端应用 (Vue 3)
│   ├── web/
│   │   ├── public/
│   │   │   ├── index.html                # HTML模板
│   │   │   └── favicon.ico               # 网站图标
│   │   │
│   │   ├── src/
│   │   │   ├── main.js                   # 应用入口
│   │   │   ├── App.vue                   # 根组件
│   │   │   │
│   │   │   ├── router/
│   │   │   │   ├── index.js              # 路由主文件
│   │   │   │   └── routes.js             # 路由定义
│   │   │   │
│   │   │   ├── stores/
│   │   │   │   ├── index.js              # Store主文件
│   │   │   │   ├── theme.js              # 主题状态
│   │   │   │   ├── auth.js               # 认证状态
│   │   │   │   ├── post.js               # 文章状态
│   │   │   │   ├── app.js                # 应用状态
│   │   │   │   └── plugin.js             # 插件状态管理
│   │   │   │
│   │   │   ├── components/
│   │   │   │   ├── common/               # 通用组件
│   │   │   │   │   ├── Header.vue
│   │   │   │   │   ├── Sidebar.vue
│   │   │   │   │   └── Footer.vue
│   │   │   │   ├── theme/                # 主题相关组件
│   │   │   │   │   ├── ThemeList.vue
│   │   │   │   │   ├── ThemeCard.vue
│   │   │   │   │   └── ThemeUpload.vue
│   │   │   │   ├── post/                 # 文章相关组件
│   │   │   │   │   ├── PostList.vue
│   │   │   │   │   ├── PostEditor.vue
│   │   │   │   │   └── PostCard.vue
│   │   │   │   ├── user/                 # 用户相关组件
│   │   │   │   │   ├── UserList.vue
│   │   │   │   │   ├── UserProfile.vue
│   │   │   │   │   └── LoginForm.vue
│   │   │   │   └── plugins/              # 插件相关组件
│   │   │   │       ├── PluginManager.vue # 插件管理界面
│   │   │   │       └── PluginStore.vue   # 插件商店界面
│   │   │   │
│   │   │   ├── views/                    # 页面视图
│   │   │   │   ├── Dashboard.vue         # 仪表板
│   │   │   │   ├── ThemeManagement.vue   # 主题管理
│   │   │   │   ├── PostManagement.vue    # 文章管理
│   │   │   │   ├── UserManagement.vue    # 用户管理
│   │   │   │   ├── SystemSettings.vue    # 系统设置
│   │   │   │   └── PluginManagement.vue  # 插件管理
│   │   │   │
│   │   │   ├── composables/
│   │   │   │   ├── useApi.js             # API调用
│   │   │   │   ├── useTheme.js           # 主题相关逻辑
│   │   │   │   ├── useAuth.js            # 认证相关逻辑
│   │   │   │   └── usePlugin.js          # 插件组合式函数
│   │   │   │
│   │   │   ├── assets/                   # 静态资源
│   │   │   │   ├── css/
│   │   │   │   └── images/
│   │   │   │
│   │   │   └── utils/                    # 前端工具函数
│   │   │       ├── request.js            # HTTP请求封装
│   │   │       ├── storage.js            # 本地存储
│   │   │       └── validate.js           # 表单验证
│   │   │
│   │   ├── package.json
│   │   ├── vite.config.js                # Vite配置
│   │   └── index.html
│   │
│   └── themes/                           # 前端主题目录
│       ├── default/                      # 默认主题
│       │   ├── assets/
│       │   ├── components/
│       │   ├── layouts/
│       │   └── style.css
│       └── .gitkeep
│
├── go.mod
├── go.sum
└── README.md