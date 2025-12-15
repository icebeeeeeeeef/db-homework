#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// 声明DLL导出的函数
extern char* toukei_count_path(const char* path);
extern char* toukei_count_with_config(const char* config_json);
extern void toukei_free_string(char* str);

int main() {
    // 示例1: 简单路径统计
    printf("=== 示例1: 简单路径统计 ===\n");
    char* result1 = toukei_count_path(".\\src");
    if (result1 != NULL) {
        printf("统计结果: %s\n", result1);
        toukei_free_string(result1);  // 必须释放内存
    } else {
        printf("统计失败\n");
    }

    // 示例2: 使用JSON配置
    printf("\n=== 示例2: 使用JSON配置 ===\n");
    const char* config_json = "{\
        \"path\": \".\\src\",\
        \"types\": [\"rs\"],\
        \"ignore_blanks\": false,\
        \"ignore_comments\": false,\
        \"show_stats\": true,\
        \"show_function_stats\": true\
    }";
    
    char* result2 = toukei_count_with_config(config_json);
    if (result2 != NULL) {
        printf("统计结果: %s\n", result2);
        toukei_free_string(result2);  // 必须释放内存
    } else {
        printf("统计失败\n");
    }

    return 0;
}