#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$ROOT_DIR/build/classes"

echo "🦆 Preparing Duck Assistant workspace..."
mkdir -p "$BUILD_DIR"

echo "🦆 Building code statistics helper..."
make -s -C "$ROOT_DIR/tools/code-stats"

echo "🦆 Compiling Java sources..."
rm -rf "$BUILD_DIR"/*
mapfile -t JAVA_SOURCES < <(find "$ROOT_DIR/src" -name "*.java")
javac -d "$BUILD_DIR" "${JAVA_SOURCES[@]}"

echo "🦆 编译完成！请选择启动方式："
echo "1. 图形界面小鸭子助手 (推荐)"
echo "2. 命令行小鸭子助手"
read -r -p "🦆 请选择 (1-2): " choice

echo "🦆 启动中..."
case "$choice" in
    2)
        java -cp "$BUILD_DIR" app.Main --duck
        ;;
    *)
        java -cp "$BUILD_DIR" app.Main --duck-gui
        ;;
esac
