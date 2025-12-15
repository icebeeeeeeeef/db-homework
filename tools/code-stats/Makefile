# Makefile for code statistics tool

CXX = g++
CXXFLAGS = -std=c++17 -Wall -Wextra -O2
TARGET = code_stats
SOURCE = code_stats.cpp

# 默认目标
all: $(TARGET)

# 编译目标
$(TARGET): $(SOURCE)
	$(CXX) $(CXXFLAGS) -o $(TARGET) $(SOURCE)

# 清理
clean:
	rm -f $(TARGET)

# 安装（可选）
install: $(TARGET)
	cp $(TARGET) /usr/local/bin/

# 测试
test: $(TARGET)
	./$(TARGET) .

# 帮助
help:
	@echo "可用的目标:"
	@echo "  all     - 编译程序 (默认)"
	@echo "  clean   - 清理编译文件"
	@echo "  test    - 测试程序（分析当前目录）"
	@echo "  install - 安装到系统"
	@echo "  help    - 显示此帮助信息"

.PHONY: all clean install test help
