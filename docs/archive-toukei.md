# Toukei - 代码统计工具

Toukei 是一个高效的代码统计工具，可以分析项目中的代码行数、空行数、注释行数，以及函数统计信息。

## 功能特性

- 支持多种编程语言的代码统计
- 统计总行数、空行数、注释行数
- 函数行数统计（支持识别各种语言的函数定义）
- 统计数据分析（均值、中位数、最大值、最小值）
- 提供 DLL 接口，支持被其他语言调用
- 使用 JSON 格式进行配置和数据传输

## 支持的编程语言

- Rust
- C++
- Python
- JavaScript
- Go
- C
- C#
- Java
- HTML
- CSS
- XML
- SCSS
- JSON
- TypeScript

## 命令行使用

```bash
# 基本使用 - 统计指定目录
cargo run -- --path your_project_path

# 显示详细统计信息
cargo run -- --path your_project_path --show-stats

# 显示函数统计信息
cargo run -- --path your_project_path --show-function-stats

# 忽略空行
cargo run -- --path your_project_path --ignore-blanks

# 忽略注释
cargo run -- --path your_project_path --ignore-comments
```

## DLL 接口使用

### 编译 DLL

```bash
cargo build --release
```

编译后，在 `target/release/` 目录下会生成 `toukei.dll` 文件。

### C 接口说明

```c
// 简单路径统计
extern char* toukei_count_path(const char* path);

// 使用 JSON 配置进行统计
extern char* toukei_count_with_config(const char* config_json);

// 释放返回的字符串内存
extern void toukei_free_string(char* str);
```

### JSON 配置格式

```json
{
    "path": "项目路径",
    "types": ["rs", "cpp", "py"],  // 可选，指定文件类型
    "ignore_blanks": false,          // 可选，是否忽略空行
    "ignore_comments": false,        // 可选，是否忽略注释
    "ignore_files": ["test.rs"],     // 可选，忽略的文件
    "show_stats": true,              // 可选，显示行数统计信息
    "show_function_stats": true      // 可选，显示函数统计信息
}
```

### JSON 返回格式

```json
[
    {
        "lang": "rs",              // 语言
        "files": 26,               // 文件数
        "lines": 1834,             // 总行数
        "blanks": 204,             // 空行数
        "comments": 8,             // 注释行数
        "functions": 38,           // 函数数
        "function_lines": [3, 5, 7], // 每个函数的行数
        "stats": {                 // 行数统计信息（如果启用）
            "min": 70,
            "max": 1834,
            "median": 70,
            "mean": 70.54
        },
        "function_stats": {        // 函数统计信息（如果启用）
            "min": 3,
            "max": 69,
            "median": 11,
            "mean": 17.47
        }
    }
]
```

## 使用示例

### C 语言调用示例

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// 声明DLL导出的函数
extern char* toukei_count_path(const char* path);
extern char* toukei_count_with_config(const char* config_json);
extern void toukei_free_string(char* str);

int main() {
    // 简单路径统计
    char* result = toukei_count_path("./src");
    if (result != NULL) {
        printf("统计结果: %s\n", result);
        toukei_free_string(result);  // 必须释放内存
    }
    
    return 0;
}
```

## 开发说明

### 项目结构

- `src/core/` - 核心语言定义
- `src/config.rs` - 配置处理
- `src/file_reader.rs` - 文件读取和统计逻辑
- `src/ffi.rs` - 外部接口
- `examples/` - 使用示例

### 添加新语言支持

在 `src/core/lang_def.rs` 中定义语言特性，然后在 `src/consts.rs` 中添加相应的语言定义。

## 许可证

MIT