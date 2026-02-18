# RBAC 权限管理系统

## 项目概述

RBAC（基于角色的访问控制）是一个完整的权限管理系统，使用 Go 语言开发，提供用户管理、角色管理、权限分配等核心功能。系统采用分层架构设计，具有良好的可扩展性和维护性。

## 技术栈

- **语言**：Go 1.25.1
- **Web 框架**：Gin v1.11.0
- **ORM 框架**：GORM v1.31.1
- **数据库**：MySQL
- **缓存**：Redis v8.11.5
- **认证**：JWT v4.5.2
- **配置管理**：TOML
- **日志**：Zap v1.27.1

## 目录结构

```
rbac/
├── api/               # API 相关代码
│   ├── error/         # 错误定义
│   ├── jwt/           # JWT 认证
│   └── sql.sql        # SQL 脚本
├── app/               # 应用入口
│   ├── config/        # 配置文件
│   └── main.go        # 主入口文件
├── conf/              # 配置管理
├── dao/               # 数据访问层
├── manager/           # 管理层
├── model/             # 数据模型
├── server/            # 服务器
│   └── http/          # HTTP 服务器
├── service/           # 业务逻辑层
├── README.md          # 项目文档
├── build.sh           # 构建脚本
├── go.mod             # Go 模块文件
├── go.sum             # 依赖校验文件
└── rbac.service       # 系统服务配置
```

## 核心功能

### 1. 用户管理
- 用户列表查询（支持分页）
- 用户详情查询
- 添加用户
- 更新用户信息
- 删除用户
- 用户登录认证（JWT token 生成）

### 2. 角色管理
- 角色列表查询（支持分页）
- 添加角色
- 修改角色信息
- 删除角色

### 3. 页面/菜单管理
- 页面列表查询（支持分页）
- 添加页面/菜单
- 修改页面/菜单信息
- 删除页面/菜单

### 4. 权限分配
- 角色与页面/菜单关联
- 用户与角色关联
- 菜单树构建（支持多级菜单）

### 5. 认证与授权
- JWT token 生成与验证
- 基于角色的权限控制

## 系统架构

系统采用分层架构设计，各层职责明确：

1. **API 层**：处理 HTTP 请求和响应，参数验证，错误处理
2. **Service 层**：实现核心业务逻辑
3. **DAO 层**：负责数据访问，与数据库交互
4. **Model 层**：定义数据模型和结构体
5. **Config 层**：管理系统配置

## 快速开始

### 环境要求
- Go 1.25.1 或更高版本
- MySQL 5.7 或更高版本
- Redis 6.0 或更高版本

### 安装与配置

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd rbac
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **配置文件**
   修改 `app/config/config.toml` 文件，配置数据库连接和其他参数：
   ```toml
   # 数据库配置
   [database]
   host = "localhost"
   port = 3306
   user = "root"
   password = "password"
   dbname = "rbac"
   
   # Redis 配置
   [redis]
   host = "localhost"
   port = 6379
   password = ""
   db = 0
   
   # 服务器配置
   [server]
   port = 8080
   logs_path = "./logs"
   ```

4. **初始化数据库**
   执行 `api/sql.sql` 文件中的 SQL 语句，创建数据库表结构。

### 构建与运行

1. **构建项目**
   ```bash
   ./build.sh
   ```

2. **运行项目**
   ```bash
   # 直接运行
   go run app/main.go
   
   # 或使用构建产物
   ./output/rbac
   ```

3. **作为系统服务运行**
   ```bash
   # 复制服务配置文件
   sudo cp rbac.service /etc/systemd/system/
   
   # 启动服务
   sudo systemctl start rbac
   
   # 查看服务状态
   sudo systemctl status rbac
   ```

## API 接口

### 用户管理
- `GET /api/rbac/account/list` - 获取用户列表
- `GET /api/rbac/account/:id` - 获取用户详情
- `POST /api/rbac/account/add` - 添加用户
- `PUT /api/rbac/account/update` - 更新用户
- `DELETE /api/rbac/account/:id` - 删除用户
- `POST /api/rbac/account/login` - 用户登录

### 角色管理
- `GET /api/rbac/role/list` - 获取角色列表
- `POST /api/rbac/role/add` - 添加角色
- `PUT /api/rbac/role/update` - 更新角色
- `DELETE /api/rbac/role/:id` - 删除角色

### 页面/菜单管理
- `GET /api/rbac/page/list` - 获取页面列表
- `POST /api/rbac/page/add` - 添加页面
- `PUT /api/rbac/page/update` - 更新页面
- `DELETE /api/rbac/page/:id` - 删除页面

### 权限分配
- `POST /api/rbac/role/page` - 分配角色页面权限
- `POST /api/rbac/account/role` - 分配用户角色

## 核心模块

### 1. JWT 认证

系统使用 JWT 进行用户认证，认证成功后生成 token，后续请求通过 token 验证用户身份。

### 2. 菜单树构建

系统支持多级菜单结构，通过 `buildMenuTreeV2` 方法构建菜单树，方便前端展示。

### 3. 权限控制

基于角色的权限控制，通过角色关联页面/菜单，用户关联角色的方式实现权限管理。

## 安全特性

- 密码加密存储
- JWT token 认证
- 权限校验
- 输入参数验证

## 部署建议

1. **生产环境**：使用 Nginx 作为反向代理，配置 SSL 证书
2. **数据库**：使用主从复制，提高可用性
3. **Redis**：使用集群模式，提高缓存性能和可用性
4. **日志**：配置日志轮转，定期清理日志文件
5. **监控**：部署 Prometheus + Grafana 监控系统运行状态

## 开发规范

1. **代码风格**：遵循 Go 语言规范，使用 `go fmt` 格式化代码
2. **命名规范**：使用驼峰命名法，变量和函数名要清晰表达其用途
3. **注释规范**：关键代码和函数要添加注释，说明其功能和实现逻辑
4. **错误处理**：统一错误处理，返回标准错误格式
5. **日志记录**：关键操作和错误要记录日志，便于问题排查

## 贡献指南

1. **Fork 项目**
2. **创建分支**：`git checkout -b feature/xxx`
3. **提交代码**：`git commit -m "feat: 添加 xxx 功能"`
4. **推送到远程**：`git push origin feature/xxx`
5. **创建 Pull Request**

## 许可证

本项目采用 MIT 许可证，详见 LICENSE 文件。

## 联系方式

如有问题或建议，欢迎联系项目维护者。

---

**备注**：本项目仅供学习和参考，生产环境使用前请根据实际需求进行调整和优化。