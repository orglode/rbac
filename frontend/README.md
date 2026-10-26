# RBAC Admin Frontend

Vue 3 + Vite + Element Plus 管理端，对接仓库内 Go API。

## 开发

```bash
# 先启动后端（监听 :23200）
cd ../app && go run main.go

# 再启动前端
cd ../frontend
npm install
npm run dev
```

访问 http://127.0.0.1:5173 ，默认账号 `admin` / `admin`。

Vite 已将 `/api` 代理到 `http://127.0.0.1:23200`。

## 目录

```
src/
  app/           # 应用壳：路由、布局
  features/      # 按功能：auth / account / role / page / login-log / operation-log
  stores/        # 仅用户会话与菜单（克制）
  shared/        # http、类型、工具
```

## 约定

- 路由懒加载；菜单二级页按 `page_path` 动态注册
- Element Plus 按需自动引入（unplugin）
- 按钮显隐：`user.can(menuPath, buttonCode)`
