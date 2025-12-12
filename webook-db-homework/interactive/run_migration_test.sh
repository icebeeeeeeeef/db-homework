#!/bin/bash

# 数据迁移测试运行脚本
echo "开始运行数据迁移测试..."

# 进入interactive目录
cd "$(dirname "$0")"

# 运行测试
go test -v -run TestDataMigration

echo "测试完成！"
