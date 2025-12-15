# 代码统计工具 (Code Statistics Tool)

一个用C++编写的代码统计工具，可以分析指定目录下的代码文件，按编程语言统计文件数量、代码行数、注释行数和空行数。

## 功能特性

- 支持多种编程语言：C/C++, Java, Python, JavaScript, TypeScript, Go, Rust, C#, PHP, Ruby, Swift, Kotlin, HTML, CSS, Shell, SQL, YAML, JSON, XML
- 递归遍历目录
- 智能识别注释和空行
- 详细的统计报告
- 支持多行注释识别
- `--dir=<path>` 指定统计根目录
- `--languages=java,python` 过滤语言列表（可同时匹配名称或缩写）
- `--functions` 额外统计 C/C++、Python 函数长度（均值/最大值/最小值/中位数）
- `--tsv` 模式可供 GUI 直接绘制柱状图 / 扇形图，并在尾部附加 `FUNC` 行

## 编译

### 使用 Makefile
```bash
cd tools/code-stats
make
```

### 手动编译
```bash
cd tools/code-stats
g++ -std=c++17 -Wall -Wextra -O2 -o code_stats code_stats.cpp
```

## 使用方法

```bash
./code_stats <目录路径>

# TSV 输出 + 指定目录
./code_stats --tsv --dir=src

# 仅统计 Java/Python，并输出函数长度
./code_stats --tsv --languages=java,python --functions
```

### 示例

```bash
# 分析当前目录
./code_stats .

# 分析指定目录
./code_stats /path/to/your/project

# 分析 Java 项目
./code_stats --languages=java /path/to/java/project
```

## 输出示例

```
正在分析目录: /path/to/project
请稍候...

=== 代码统计报告 ===

语言            文件数   总行数    代码行    注释行    空行
---------------------------------------------------------------------------
Java            25      1250     980       150       120
C/C++           15      800      650       100       50
Python          10      400      320       50        30
JavaScript      8       200      180       15        5
总计            58      2650     2130      315       205

FUNC    C/C++   24      14.6     3         42        13.0
FUNC    Python  18      8.1      2         19        7.0

=== 详细统计 ===

Java:
  文件数: 25
  总行数: 1250
  代码行: 980 (78.4%)
  注释行: 150 (12.0%)
  空行: 120 (9.6%)
```

## 支持的编程语言

| 语言 | 文件扩展名 | 单行注释 | 多行注释 |
|------|------------|----------|----------|
| C/C++ | .cpp, .c, .cc, .cxx, .c++, .h, .hpp, .hxx | // | /* */ |
| Java | .java | // | /* */ |
| Python | .py, .pyw | # | """ """ |
| JavaScript | .js, .jsx, .mjs | // | /* */ |
| TypeScript | .ts, .tsx | // | /* */ |
| Go | .go | // | /* */ |
| Rust | .rs | // | /* */ |
| C# | .cs | // | /* */ |
| PHP | .php, .phtml | //, # | /* */ |
| Ruby | .rb | # | =begin =end |
| Swift | .swift | // | /* */ |
| Kotlin | .kt, .kts | // | /* */ |
| HTML | .html, .htm | - | <!-- --> |
| CSS | .css, .scss, .sass, .less | - | /* */ |
| Shell | .sh, .bash, .zsh, .fish | # | - |
| SQL | .sql | -- | /* */ |
| YAML | .yml, .yaml | # | - |
| JSON | .json | - | - |
| XML | .xml | - | <!-- --> |

## 系统要求

- C++17 或更高版本
- 支持 std::filesystem 的编译器（GCC 8+, Clang 7+, MSVC 2017+）

## 注意事项

- 工具会递归遍历指定目录下的所有子目录
- 只统计支持的文件类型
- 空行定义为只包含空白字符的行
- 注释行包括单行注释和多行注释
- 多行注释会正确处理开始和结束标记

## 故障排除

如果编译时出现 `std::filesystem` 相关的错误，请确保：

1. 使用 C++17 或更高版本
2. 对于较老的 GCC 版本，可能需要链接 `-lstdc++fs`
3. 对于较老的 Clang 版本，可能需要链接 `-lc++fs`

## 许可证

此工具为开源软件，可自由使用和修改。
