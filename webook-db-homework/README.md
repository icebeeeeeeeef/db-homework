# WeBook Homework (Backend + Frontend)

精简版单体项目，保留 MySQL + Redis + 前端展示，用于本地演示与作业：
- **webook-db-homework**：后端，提供用户、文章、互动接口（点赞/收藏/阅读计数）。
- **webook-fe**：前端（Next.js + Ant Design Pro），包含登录/个人信息/文章列表与详情。

## 功能
- 用户：邮箱注册、登录/退出、个人信息查看与修改（昵称/生日/AboutMe），登录态基于 JWT + Redis SSID。
- 文章：草稿保存、发布、撤回、我的文章列表（草稿/撤回/已发布）、公开列表与公开详情。
- 互动：阅读计数、点赞/取消点赞、收藏，用户是否点赞/收藏与总数均持久化 MySQL。
- 数据脚本：`script/mysql/seed_data.sql` 预置 3 个用户和多篇文章、互动数据（密码统一 `Passw0rd!`）。

## 快速启动
### 0. 依赖
- Go 1.23+（后端）
- Node.js 18+（前端）
- MySQL 8.x、Redis 7.x

### 1. 启动数据库与缓存
```bash
# MySQL 8.4
docker run -d --name webook-dbhome-mysql -p 13316:3306 \
  -e MYSQL_ROOT_PASSWORD=123123 \
  -e MYSQL_DATABASE=webook_db \
  mysql:8.4

# Redis 7.4
docker run -d --name webook-dbhome-redis -p 6379:6379 redis:7.4
```
`configs/dev.yaml` 默认指向 `localhost:13316` / `localhost:6379`，如端口修改请同步调整。

### 2. 导入演示数据（可选，幂等）
```bash
mysql -h127.0.0.1 -P13316 -uroot -p123123 < script/mysql/seed_data.sql
# 预置账号：alice/bob/carol@example.com，密码均为 Passw0rd!
```

### 3. 启动后端
```bash
cd webook-db-homework
go run .
# 健康检查 http://localhost:8080/health
```

### 4. 启动前端
```bash
cd ../webook-fe
npm install
npm run dev   # http://localhost:3000
```
登录后进入文章主页，可查看/编辑文章、点赞/收藏、修改个人信息。

## 主要接口（后端）
- 用户：`POST /users/signup`、`POST /users/login`、`POST /users/logout`、`GET /users/profile`、`POST /users/edit`
- 文章：`POST /articles/edit`、`POST /articles/publish`、`POST /articles/withdraw`、`POST /articles/list`、`GET /articles/detail/:id`、`GET /articles/pub/:id`、`POST /articles/pub/list`
- 互动：`POST /articles/pub/like`、`POST /articles/pub/collect`（需登录）；公开详情免登录，但带 token 会返回当前用户的点赞/收藏状态

## 项目结构（精简后）
- `main.go` 手写依赖注入，使用 `configs/dev.yaml` 初始化 MySQL/Redis
- `internal/` 业务分层：repository（DAO+缓存）、service、web（Handler/中间件）
- `interactive/` 互动模块 DAO/缓存/仓储（单体内调用）
- `script/mysql/seed_data.sql` 演示数据脚本
- `webook-fe/` 前端源码

## 注意
- 已移除微服务/Kafka/排行榜等未实现的部分，只保留单体必需代码。
- 撤回文章会同步更新 reader 表并清理公开缓存；公开接口仅返回已发布文章。
- 如果修改数据库端口/账号，请同步更新 `configs/dev.yaml` 和前端 Axios 基础地址（`src/axios/axios.ts`）。***
