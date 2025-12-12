package ioc

import (
	"time"
	"webook/pkg/logger"

	dao2 "webook/interactive/repository/dao"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	dsn := viper.GetString("db.mysql.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: glogger.New(gormLoggerfunc(l.Debug), glogger.Config{
			SlowThreshold:             time.Millisecond * 50, //多少毫秒算慢
			IgnoreRecordNotFoundError: true,                  //忽略记录不存在错误,这个要分情况讨论，查不到是否算作错误
			LogLevel:                  glogger.Info,          //日志级别
		}),
	})
	if err != nil {
		panic(err)
	}
	err = dao2.InitTables(db)
	if err != nil {
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
