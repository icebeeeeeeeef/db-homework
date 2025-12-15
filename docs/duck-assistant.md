# 🦆 Duck Assistant

一个可爱的交互式助手，集成了代码统计、抢红包游戏和 AI 对话功能，现已支持与大模型实时对话（默认人格：唐老鸭 Donald Duck）。

## 🚀 快速启动

### 方法 1：使用启动脚本（推荐）
```bash
./scripts/start_duck_gui.sh   # 直接启动图形化版本
```

### 方法 2：手动启动
```bash
# 构建 code_stats 并编译所有 Java 源码
make -s -C tools/code-stats
mkdir -p build/classes
javac -d build/classes $(find src -name "*.java")

# 启动图形化小鸭子助手
java -cp build/classes app.Main --duck-gui
```

> 如果希望体验命令行版本，可执行：
>
> ```bash
> java Main --duck
> ```

## 🔑 配置 AI 大模型

Duck Assistant 使用 OpenAI Chat Completions 兼容接口。运行前请设置以下环境变量：

| 变量名             | 必填 | 说明                                  | 默认值                              |
| ------------------ | ---- | ------------------------------------- | ----------------------------------- |
| `OPENAI_API_KEY`   | ✅    | 模型服务的 API Key                    | 无                                  |
| `OPENAI_BASE_URL`  | ❌    | 接口 Base URL（末尾不要带 `/`）       | `https://api.openai.com`            |
| `OPENAI_MODEL`     | ❌    | Chat 模型名称                         | `gpt-3.5-turbo`                     |

示例（Linux / macOS）：
```bash
export OPENAI_API_KEY=sk-xxxxxxxx
export OPENAI_BASE_URL=https://api.openai.com
export OPENAI_MODEL=gpt-4o-mini
```

未配置 API Key 时，小鸭子会自动回退到离线的本地回复逻辑，保证功能可用。

## 🎯 主要功能

### 1. 📊 代码统计
- 调用升级版 `code_stats` 工具，可通过 `--dir`、`--languages`、`--functions` 精准分析
- GUI 支持「汇总表 + 柱状图 + 扇形图」多视角切换
- 可选择是否展示空行/注释行列
- 自动统计 C/C++ 与 Python 函数长度（均值 / 中位数 / 极值），用于作业分析

### 2. 🎁 抢红包游戏
- 图形界面 / 命令行双模式
- 新增随机场景（夕阳/SNOW/星空）+ 背景装饰
- HUD 展示剩余时间、场景名、碰撞个数/金额
- 红包形状、大小、统计全面升级，结束后弹窗展示多维统计

### 3. 🤖 AI 对话（唐老鸭人格）
- 默认人格：唐老鸭，语言风格活泼、幽默
- 支持多轮对话，保持上下文
- UI 英文化，避免字体编码问题
- 自动检测中文输入并以中文回复（保持唐老鸭语气）

### 4. 🧢 装扮系统
- 使用装饰器模式管理唐老鸭 + 三只小鸭子的帽子、围巾、眼镜、领带、手杖
- 支持颜色挑选、角色切换、一键同步到全部小鸭子
- 默认造型更贴近经典唐老鸭 + 小鸭子组合

### 5. 🎬 唐老鸭舞台命令
- 右侧舞台展示唐老鸭与三只小鸭子，可点击唐老鸭调出指令面板
- 支持「红包雨」「代码统计」「AI 问答」三大任务（可传入目录、语言过滤、是否启动游戏等参数）
- 舞台可显示红包雨动画、统计状态、AI 思考提示，并与聊天区同步播报

## 💻 界面概览

启动 GUI 后，右侧展示唐老鸭头像和功能按钮，左侧为聊天与控制台。示例按钮（英文 UI）：

```
📊 Code Stats   🎁 Red Packet
🧢 Dress Up     🤖 AI Helper
🚪 Exit
```

## 🧠 AI 对话工作流
1. 玩家输入消息 → 显示在聊天框。
2. 若配置了大模型凭据，Duck Assistant 将携带历史上下文调用 OpenAI 兼容接口。
3. AI 回复以唐老鸭语气展示；若调用失败，则自动回退到离线规则回复。
4. 对话历史保存在内存 List 中，刷新窗口会重置上下文。

## 🔧 目录结构（关键部分）
```
ai/                    # AI 客户端（OpenAI 兼容）
├── AiClient.java

gui/                   # Swing 图形界面
├── DuckAssistantGUI.java
├── DuckAvatarPanel.java
├── ...
```

## 🧪 测试指南
```bash
# 1. 设置环境变量（如有）
export OPENAI_API_KEY=...

# 2. 编译并启动
javac *.java model/*.java geom/*.java gui/*.java ai/*.java
java Main --duck-gui

# 3. 在聊天框输入消息
Hi Duck!  /  你好唐老鸭！
```

## 🐛 常见问题

### 1. 启动提示找不到 `code_stats`
请先编译 C++ 工具：
```bash
make -s -C tools/code-stats
```

### 2. AI 回复报错 `(AI error: ...)`
- 检查网络是否能访问配置的 Base URL
- 核对 API Key 是否有效
- 查看终端输出是否有 HTTP 状态码/错误信息

### 3. 字体显示为方块
界面已切换为英文按钮与提示，若仍有方块，可在系统安装中文字体（例如 Noto Sans CJK）。

## 🎉 Have fun!

唐老鸭已经准备好与玩家互动，祝你玩得开心！🦆
