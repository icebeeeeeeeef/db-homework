# Red Packet Duck Playground

This monorepo combines the Swing-based Duck Assistant, the console red-packet mini game, and the standalone C++ code statistics helper. The codebase is now grouped by feature area so it is easier to browse and hack on.

## Directory map

- `src/` – all Java sources grouped by package (entrypoint, CLI assistant, GUI, models, AI client, shared utilities)
- `scripts/` – helper launchers for compiling and starting either the GUI or CLI version
- `tools/` – native helpers such as the `code_stats` analyzer and archived binaries
- `docs/` – feature guides and legacy documentation

See the individual README inside every directory for structure and logic details.

## Quick start

```bash
# Build the stats helper + compile all Java sources, then choose GUI/CLI
./scripts/start_duck.sh

# Or launch the GUI directly
./scripts/start_duck_gui.sh
```

Both scripts compile sources into `build/classes` and keep the repository tree clean. The GUI uses the packaged `tools/code-stats/code_stats` binary for live charts; rebuild it with `make -C tools/code-stats` if you change the C++ code.

## Classroom Attendance System (Java + Swing + MySQL + Docker)

1) 启动数据库（第一次会自动建表并导入 10 条学生数据）:
```bash
docker-compose up -d
```
默认配置：数据库 `attendance_db`，用户 `attendance_user` / `attendance_pass`，端口 `3306`。
如果之前启动过容器且密码不一致，请重置数据卷后重启：
```bash
docker-compose down -v
docker-compose up -d
```

若端口被占用，可修改 `docker-compose.yml` 的映射并用环境变量或 JVM 参数告知程序，例如端口改为 3366：
```bash
export DB_PORT=3366   # 或运行时加 -Ddb.port=3366
./scripts/start_duck_gui.sh   # 或 ./scripts/start_attendance.sh
```

2) 运行 Swing 点名 UI（可选设置 MYSQL_JAR 指向 mysql-connector-j 驱动以便本地连接）:
```bash
MYSQL_JAR=/path/to/mysql-connector-j.jar ./scripts/start_attendance.sh
```
或直接 `./scripts/start_attendance.sh` 若驱动已在全局 CLASSPATH。

3) 功能：
- 学生管理：添加学生、查看列表
- 点名：全点 / 抽点，抽点人数可配置，记录到 / 未到，自动写入 attendance_session / attendance_record
- 统计：查看本次点名结果与历史缺勤率（缺勤率 = 缺勤次数 / 总次数）
