//go:build test

package startup

import (
	"webook/internal/repository/dao"
	"webook/pkg/logger"

	"errors"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDBNoLogger 测试环境专用：不依赖日志参数的数据库初始化
func InitDBNoLogger() (*gorm.DB, error) {
	dsn := viper.GetString("db.mysql.dsn")
	if dsn == "" {
		dsn = os.Getenv("MYSQL_DSN")
	}
	if dsn == "" {
		return nil, errors.New("empty DSN: set viper key db.mysql.dsn or env MYSQL_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = dao.InitTables(db); err != nil {
		return nil, err
	}
	return db, nil
}

// InitDB 测试环境专用：兼容 wire 的数据库初始化函数
func InitDB(l logger.LoggerV1) *gorm.DB {
	db, err := InitDBNoLogger()
	if err != nil {
		panic(err)
	}
	return db
}
