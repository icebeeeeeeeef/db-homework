//go:build e2e

package startup

import (
	"fmt"
	"os"
	"testing"
	"time"
	"webook/internal/repository/dao"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitTestDB 初始化测试数据库
// 这个函数专门用于集成测试，会创建一个独立的测试数据库
func InitTestDB() *gorm.DB {
	// 优先使用环境变量
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		// 然后尝试从viper获取
		dsn = viper.GetString("db.mysql.dsn")
	}
	if dsn == "" {
		// 默认的测试数据库配置
		dsn = "root:123123@tcp(localhost:3306)/webook_test?parseTime=true&loc=Local&charset=utf8mb4"
	}

	// 创建测试数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 测试时关闭SQL日志
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		panic(fmt.Sprintf("连接测试数据库失败: %v", err))
	}

	// 初始化表结构
	if err = dao.InitTables(db); err != nil {
		panic(fmt.Sprintf("初始化测试表失败: %v", err))
	}

	return db
}

// InitTestDBWithLogger 带日志的测试数据库初始化
func InitTestDBWithLogger() *gorm.DB {
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		dsn = viper.GetString("db.mysql.dsn")
	}
	if dsn == "" {
		dsn = "root:123123@tcp(localhost:3306)/webook_test?parseTime=true&loc=Local&charset=utf8mb4"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 测试时显示SQL日志
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		panic(fmt.Sprintf("连接测试数据库失败: %v", err))
	}

	if err = dao.InitTables(db); err != nil {
		panic(fmt.Sprintf("初始化测试表失败: %v", err))
	}

	return db
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	// 获取所有表名
	var tables []string
	if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return err
	}

	// 禁用外键检查
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	// 清空所有表
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table)).Error; err != nil {
			return err
		}
	}

	// 重新启用外键检查
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return err
	}

	return nil
}

// ResetAutoIncrement 重置自增ID
func ResetAutoIncrement(db *gorm.DB, table string) error {
	return db.Exec(fmt.Sprintf("ALTER TABLE `%s` AUTO_INCREMENT = 1", table)).Error
}

// SetupTestDB 设置测试数据库，返回清理函数
func SetupTestDB(t *testing.T) *gorm.DB {
	db := InitTestDB()

	// 注册清理函数
	t.Cleanup(func() {
		if err := CleanupTestDB(db); err != nil {
			t.Logf("清理测试数据库失败: %v", err)
		}
	})

	return db
}
