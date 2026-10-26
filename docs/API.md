# RBAC 管理端接口文档

> 供前端对接。鉴权：`Authorization: Bearer <token>`  
> 统一响应：

```json
{
  "code": 0,
  "msg": "success",
  "data": {},
  "trace_id": "xxx"
}
```

`code !== 0` 为失败。常见码：`499` 缺参，`500` 内部错误，`600` 账号已存在，`601` 角色不存在，`602` 权限无效，`603` 手机号已存在，`604` 账号密码错误。HTTP 401 表示未登录/token 无效。

---

## 权限模型（前端必读）

```
账号 Account
  └── 角色 Role（当前一人一角）
        └── 页面/按钮 Page（role_page）
```

| page_type | page_level | 含义 |
|-----------|------------|------|
| 1 菜单 | 1 | 一级菜单 |
| 1 菜单 | 2 | 二级菜单（挂页面路由 `url`） |
| 2 按钮 | 3 | 按钮；绑定一个 API：`api_method` + `api_path` |

默认 CRUD 按钮编码：

| button_code | 名称 | 默认 method | 默认 path（相对 api_prefix） |
|-------------|------|-------------|------------------------------|
| query | 查询 | GET | `{api_prefix}` |
| create | 新增 | POST | `{api_prefix}` |
| update | 修改 | PUT | `{api_prefix}/:id` |
| delete | 删除 | DELETE | `{api_prefix}/:id` |

**后端鉴权规则**

1. 登录接口不做鉴权。
2. 已登录后：若该 `METHOD + Gin FullPath` 在 `page` 中存在对应按钮，则角色必须拥有该按钮，否则 `602`。
3. 库中无对应按钮的接口：仅校验登录（逐步挂按钮即可收紧）。
4. 权限全部由角色控制，无账号类型旁路。
5. 白名单（登录后免按钮校验）：`GET /api/admin/me/menu`、`GET /api/admin/me/buttons`；健康检查 `GET /api/admin/health` 无需登录。

**前端建议**

1. 登录后调 `GET /api/admin/me/menu` 渲染侧栏（只用 `page_type=1` 的树）。
2. 调 `GET /api/admin/me/buttons` 或使用菜单里二级节点的 `buttons[]`，按 `button_code` 控制按钮显隐。
3. `api_path` 含 `:id`，与后端路由一致，不要自己拼动态数字去匹配权限。

---

## 1. 登录

`POST /api/admin/login`（无需 Token）

**请求**

```json
{
  "account_number": "admin",
  "password": "明文或与库中一致的密码"
}
```

> 密码入库为 MD5；登录时对明文做 MD5 后再比对。初始账号见 `api/rbac.sql`：`admin` / `admin`，绑定角色「系统管理员」（全量权限）。

**响应 data**

```json
{
  "id": 1,
  "token": "jwt...",
  "user_name": "管理员",
  "mobile": "13800000000",
  "account_number": "admin",
  "status": 1,
  "last_login_time": 0,
  "create_time": 1710000000
}
```

后续请求头：`Authorization: Bearer {token}`

---

## 2. 当前用户（白名单）

### 2.1 我的菜单树

`GET /api/admin/me/menu`

**响应 data**

```json
{
  "account_id": 1,
  "role_id": 2,
  "role_name": "运营",
  "page": [
    {
      "page_id": 1,
      "page_level": 1,
      "page_name": "系统管理",
      "page_path": "/system",
      "page_type": 1,
      "child_page": [
        {
          "page_id": 2,
          "page_level": 2,
          "page_name": "账号管理",
          "page_path": "/system/account",
          "page_type": 1,
          "child_page": [],
          "buttons": [
            {
              "page_id": 10,
              "page_level": 3,
              "page_name": "查询",
              "page_type": 2,
              "button_code": "query",
              "api_method": "GET",
              "api_path": "/api/admin/account",
              "child_page": []
            }
          ]
        }
      ]
    }
  ]
}
```

### 2.2 我的按钮列表（扁平）

`GET /api/admin/me/buttons`

**响应 data**：`PageInfo[]`，仅 `page_type=2` 的按钮，字段同上。

---

## 3. 账号管理

| 方法 | 路径 | 说明 | 对应 button_code |
|------|------|------|------------------|
| GET | `/api/admin/account` | 分页列表 | query |
| POST | `/api/admin/account` | 新增 | create |
| PUT | `/api/admin/account/:id` | 修改 | update |
| DELETE | `/api/admin/account/:id` | 删除 | delete |
| POST | `/api/admin/account/role` | 分配角色 | assign_role |
| DELETE | `/api/admin/account/role` | 移除角色 | remove_role |

账号 `status`：`1` 启用，`2` 禁用。禁用后无法登录，已登录请求也会被拦（605）。

### 3.1 列表

`GET /api/admin/account?page=1&size=20&user_name=&account_number=&mobile=&status=`

| 参数 | 说明 |
|------|------|
| `user_name` | 姓名模糊 |
| `account_number` | 账号模糊 |
| `mobile` | 手机号模糊 |
| `status` | `1`启用 `2`禁用；不传不限 |

**响应 data**

```json
{
  "list": [
    {
      "id": 1,
      "user_name": "张三",
      "mobile": "13800000000",
      "account_number": "zhangsan",
      "password": "",
      "status": 1,
      "last_login_time": 0,
      "create_time": 1710000000,
      "update_time": 1710000000
    }
  ],
  "total": 1
}
```

### 3.2 新增

```json
{
  "user_name": "张三",
  "mobile": "13800000001",
  "account_number": "zhangsan",
  "password": "123456",
  "status": 1
}
```

`status` 不传默认启用。

### 3.3 修改

`PUT /api/admin/account/1`

```json
{
  "user_name": "张三",
  "mobile": "13800000001",
  "account_number": "zhangsan",
  "password": "",
  "status": 1
}
```

`password` 空则不改密码。

### 3.4 分配 / 移除角色

`POST /api/admin/account/role`

```json
{ "account_id": 2, "role_id": 1 }
```

`DELETE /api/admin/account/role`

```json
{ "account_id": 2 }
```

### 3.5 登录日志

`GET /api/admin/login/log?page=1&size=20&account_number=&account_id=&status=`

可选筛选：`account_id`、`account_number`、`status`（1成功 2失败）。

**响应 data**

```json
{
  "list": [
    {
      "id": 1,
      "account_id": 1,
      "account_number": "admin",
      "login_ip": "127.0.0.1",
      "user_agent": "Mozilla/5.0...",
      "status": 1,
      "message": "登录成功",
      "create_time": 1710000000
    }
  ],
  "total": 1
}
```

对应按钮：`GET /api/admin/login/log`（菜单「登录日志」）。

### 3.6 操作日志

`GET /api/admin/operation/log?page=1&size=20&module=&action=&operator_account=&target_name=`

可选筛选：

| 参数 | 说明 |
|------|------|
| module | 1角色 2角色页面权限 3账号角色 4页面资源 |
| action | 1新增 2修改 3删除 4分配 5移除 |
| operator_account | 操作人账号模糊 |
| target_name | 目标名称模糊 |

自动写入场景：角色增删改、页面增删改、角色页面权限分配、账号被分配/变更/移除角色。

账号角色日志示例：`账号「张三(zhangsan)」被分配角色「运营」`；变更时为 `角色由「A」变更为「B」`。

**响应 data**

```json
{
  "list": [
    {
      "id": 1,
      "module": 1,
      "action": 1,
      "target_type": "role",
      "target_id": 2,
      "target_name": "运营",
      "content": "新增角色：运营",
      "operator_id": 1,
      "operator_name": "管理员",
      "operator_account": "admin",
      "create_time": 1710000000
    }
  ],
  "total": 1
}
```

对应按钮：`GET /api/admin/operation/log`（菜单「操作日志」）。

---

## 4. 角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/role?page=&size=` | 分页列表 |
| GET | `/api/admin/role/all` | 全部（下拉） |
| POST | `/api/admin/role` | 新增 `{ "role_name": "运营" }` |
| PUT | `/api/admin/role/:id` | 修改 |
| DELETE | `/api/admin/role/:id` | 删除 |
| GET | `/api/admin/role/page/:role_id` | 角色已绑菜单树（含 buttons） |
| POST | `/api/admin/role/page` | 覆盖绑定页面/按钮 |

### 4.1 分配角色权限（覆盖式）

`POST /api/admin/role/page`

```json
{
  "role_id": 1,
  "page_ids": "1,2,10,11,12,13"
}
```

`page_ids`：逗号分隔，**建议同时勾选菜单节点与按钮节点**（前端树多选）。先删该角色旧关联再写入。

---

## 5. 页面 / 按钮管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/page?page=&size=&page_name=&api=` | 一级分页树列表（1:N:N） |
| GET | `/api/admin/page/all` | 全部扁平（权限树/父级下拉） |
| POST | `/api/admin/page` | 新增菜单或按钮 |
| PUT | `/api/admin/page/:id` | 修改 |
| DELETE | `/api/admin/page/:id` | 删除（级联子节点与 role_page） |

### 5.0 树形列表（一级→二级→按钮）

`GET /api/admin/page?page=1&size=20&page_name=&api=`

| 参数 | 说明 |
|------|------|
| `page` / `size` | 按**一级菜单**分页 |
| `page_name` | 名称模糊；命中后带回关联上下级 |
| `api` | 对 `api_method` / `api_path` 模糊匹配 |

**响应 data**：`list` 为 `PageInfo[]` 树（与 `me/menu` 同结构：`child_page` + `buttons`），`total` 为一级菜单总数。

```json
{
  "list": [
    {
      "page_id": 1,
      "page_level": 1,
      "page_name": "系统管理",
      "page_path": "/system",
      "page_type": 1,
      "child_page": [
        {
          "page_id": 2,
          "page_level": 2,
          "page_name": "账号管理",
          "page_path": "/system/account",
          "page_type": 1,
          "child_page": [],
          "buttons": [
            {
              "page_id": 10,
              "page_level": 3,
              "page_name": "查询",
              "page_type": 2,
              "button_code": "query",
              "api_method": "GET",
              "api_path": "/api/admin/account",
              "child_page": []
            }
          ]
        }
      ]
    }
  ],
  "total": 1
}
```

### 5.1 新增一级/二级菜单

```json
{
  "page_name": "账号管理",
  "icon": "User",
  "url": "/system/account",
  "parent_id": 1,
  "page_level": 2,
  "page_type": 1,
  "sort": 10,
  "api_prefix": "/api/admin/account"
}
```

- `page_type` 默认 `1`（菜单）。
- `icon`：菜单图标，存 `@element-plus/icons-vue` 组件名（如 `Setting`、`User`），一级/二级菜单可用；按钮可空。
- `sort`：展示排序，**数字越大越靠前**；相同时按 `create_time` 倒序（默认 0）。
- 二级菜单且传入 `api_prefix` 时，**自动创建** 查询/新增/修改/删除 四个按钮（level=3）。
- `api_prefix` 不入库，仅创建时使用。

### 5.2 单独新增按钮

```json
{
  "page_name": "导出",
  "parent_id": 2,
  "page_level": 3,
  "page_type": 2,
  "button_code": "export",
  "api_method": "POST",
  "api_path": "/api/admin/account/export"
}
```

### 5.3 全量扁平字段（`/page/all`）

含：`id, page_name, icon, url, parent_id, page_level, page_type, button_code, api_method, api_path, sort, create_time`。

列表/菜单树统一按 `sort DESC, create_time DESC` 排序。

---

## 6. 站内信

消息类型：`1` 系统广播（全员） / `2` 业务通知（可指定人或默认全员，`biz_type` 可选） / `3` 指定用户（`account_ids` 必填）。  
收件箱行含已读状态；发送时为每位接收人写入 `site_message_inbox`。

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| GET | `/api/admin/message/unread-count` | 未读数 `{ count }` | 登录白名单 |
| GET | `/api/admin/message/inbox/recent` | 最近未读（顶栏） | 登录白名单 |
| GET | `/api/admin/message/inbox?page=&size=&msg_type=&read_status=` | 我的收件箱；`read_status` 0不限/1未读/2已读 | 登录白名单 |
| POST | `/api/admin/message/read` | 标已读 `{ inbox_ids:[], all:false }` | 登录白名单 |
| GET | `/api/admin/message` | 发送记录（含已读/接收数） | 按钮 query |
| POST | `/api/admin/message` | 发送 `{ title, content, msg_type, biz_type?, account_ids? }` | 按钮 create |

菜单：`/system/message`（站内信）。

---

## 7. 健康检查

`GET /api/admin/health`（公开，无需登录）

**响应 data**

```json
{ "time": 1710000000, "ip": "127.0.0.1" }
```

- `time`：服务端当前 Unix 时间戳
- `ip`：请求客户端 IP

---

## 8. 错误码一览

| code | 说明 |
|------|------|
| 0 | 成功 |
| 499 | 缺少参数 |
| 500 | 内部错误 |
| 600 | 账号已存在 |
| 601 | 角色不存在 |
| 602 | 权限无效 |
| 603 | 手机号已存在 |
| 604 | 账号密码错误 |
| 605 | 账号已禁用 |
| 401 | HTTP：未授权 |

---

## 9. 初始化数据库

全部初始化 SQL 只维护 `api/rbac.sql`（建库 `rbac` + 表 + 种子数据）：

```bash
mysql -uroot < api/rbac.sql
```

| 种子 | 值 |
|------|-----|
| 数据库 | `rbac` |
| `role` id=1 | 系统管理员，绑定全部 `page` |
| `account` admin | 密码 `admin`（MD5） |
| `account_role` | admin → 系统管理员 |

权限只认角色，无 `account_type`。
