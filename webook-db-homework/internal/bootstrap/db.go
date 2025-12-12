//go:build !test

package bootstrap

import (
	"fmt"
	"time"
	"webook/internal/repository/dao"
	"webook/pkg/logger"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	dsn := viper.GetString("db.mysql.dsn")
	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: glogger.New(gormLoggerfunc(l.Debug), glogger.Config{
				SlowThreshold:             time.Millisecond * 50,
				IgnoreRecordNotFoundError: true,
				LogLevel:                  glogger.Info,
			}),
		})
		if err == nil {
			// 确认底层连接可用
			sqlDB, er := db.DB()
			if er == nil && sqlDB.Ping() == nil {
				break
			}
			err = er
		}
		l.Warn("数据库连接失败，重试中", logger.Error(err), logger.Int("retry", i+1))
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		panic(fmt.Errorf("connect mysql failed after retries: %w", err))
	}
	if err = dao.InitTables(db); err != nil {
		panic(err)
	}
	return db
}

type gormLoggerfunc func(msg string, fields ...logger.Field)

func (g gormLoggerfunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{
		Key:   "args",
		Value: args,
	})
}
