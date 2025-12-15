#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$ROOT_DIR/build/classes"

echo "🎓 准备编译课堂点名系统..."
mkdir -p "$BUILD_DIR"
rm -rf "$BUILD_DIR"/*

mapfile -t JAVA_SOURCES < <(find "$ROOT_DIR/src" -name "*.java")
javac -d "$BUILD_DIR" "${JAVA_SOURCES[@]}"

CLASSPATH="$BUILD_DIR"
if [[ -n "${MYSQL_JAR:-}" ]]; then
  CLASSPATH="$MYSQL_JAR:$CLASSPATH"
else
  DRIVER_JAR="$ROOT_DIR/libs/mysql-connector-j-8.3.0.jar"
  if [[ ! -f "$DRIVER_JAR" ]]; then
    echo "🎓 未检测到 MySQL 驱动，正在下载 mysql-connector-j-8.3.0.jar ..."
    mkdir -p "$ROOT_DIR/libs"
    curl -L -o "$DRIVER_JAR" "https://repo1.maven.org/maven2/com/mysql/mysql-connector-j/8.3.0/mysql-connector-j-8.3.0.jar"
  fi
  CLASSPATH="$DRIVER_JAR:$CLASSPATH"
fi

echo "🎓 启动课堂点名系统 UI..."
java -cp "$CLASSPATH" app.Main --attendance
