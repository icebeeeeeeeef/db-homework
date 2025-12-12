# Makefile for webook project

.PHONY: test test-e2e test-article test-user clean-db

# 运行所有测试
test:
	go test -v ./...

# 运行 E2E 集成测试
test-e2e:
	go test -v -tags=e2e ./internal/integration/

# 运行文章测试
test-article:
	go test -v -tags=e2e -run TestArticle ./internal/integration/

# 运行用户测试
test-user:
	go test -v -tags=e2e -run TestUser ./internal/integration/

# 清理测试数据库
clean-db:
	mysql -h localhost -u root -p123123 -e "DROP DATABASE IF EXISTS webook_test; CREATE DATABASE webook_test;"

# 运行测试前清理数据库
test-e2e-clean: clean-db test-e2e

# 运行文章测试前清理数据库
test-article-clean: clean-db test-article

# 运行用户测试前清理数据库
test-user-clean: clean-db test-user






