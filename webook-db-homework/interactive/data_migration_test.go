package main

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"webook/interactive/repository/dao"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

// 测试数据迁移 - 向interactive表插入1000条测试数据
func TestDataMigration(t *testing.T) {
	// 设置数据库配置
	viper.Set("db.mysql.dsn", "root:123123@tcp(172.28.240.236:13316)/webook_db")

	// 初始化数据库连接
	db := initDB()
	require.NotNil(t, db)

	// 确保表结构存在
	err := db.AutoMigrate(&dao.Interactive{})
	require.NoError(t, err)

	// 生成并插入1000条测试数据
	ctx := context.Background()
	err = insertTestData(ctx, db, 1000)
	require.NoError(t, err)

	// 验证数据插入成功
	var count int64
	err = db.Model(&dao.Interactive{}).Count(&count).Error
	require.NoError(t, err)

	fmt.Printf("成功插入 %d 条测试数据到 interactive 表\n", count)
	require.True(t, count >= 1000, "插入的数据条数应该至少为1000条")
}

// 初始化数据库连接
func initDB() *gorm.DB {
	dsn := viper.GetString("db.mysql.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(func(msg string, fields ...interface{}) {
			// 简化日志输出
		}), glogger.Config{
			SlowThreshold: time.Second,
			LogLevel:      glogger.Warn,
		}),
	})
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}
	return db
}

// 插入测试数据
func insertTestData(ctx context.Context, db *gorm.DB, count int) error {
	// 业务类型列表
	bizTypes := []string{"article", "comment", "post", "video", "image", "document"}

	// 批量插入，每批100条
	batchSize := 100
	for i := 0; i < count; i += batchSize {
		var interactives []dao.Interactive

		// 生成当前批次的数据
		currentBatchSize := batchSize
		if i+batchSize > count {
			currentBatchSize = count - i
		}

		for j := 0; j < currentBatchSize; j++ {
			interactive := generateRandomInteractive(bizTypes, int64(i+j+1))
			interactives = append(interactives, interactive)
		}

		// 批量插入
		err := db.WithContext(ctx).CreateInBatches(interactives, batchSize).Error
		if err != nil {
			return fmt.Errorf("批量插入第 %d-%d 条数据失败: %v", i+1, i+currentBatchSize, err)
		}

		fmt.Printf("已插入第 %d-%d 条数据\n", i+1, i+currentBatchSize)
	}

	return nil
}

// 生成随机Interactive数据
func generateRandomInteractive(bizTypes []string, bizId int64) dao.Interactive {
	now := time.Now().UnixMilli()

	// 随机选择业务类型
	biz := bizTypes[rand.Intn(len(bizTypes))]

	// 生成随机计数数据
	readcnt := rand.Int63n(10000) + 1  // 1-10000
	likecnt := rand.Int63n(1000) + 1   // 1-1000
	collectcnt := rand.Int63n(500) + 1 // 1-500

	// 添加一些随机时间偏移
	timeOffset := rand.Int63n(86400*30) * 1000 // 30天内的随机偏移
	createdAt := now - timeOffset
	updatedAt := now - timeOffset + rand.Int63n(3600*1000) // 创建后1小时内的随机更新时间

	return dao.Interactive{
		BizId:      bizId,
		Biz:        biz,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		Readcnt:    readcnt,
		Likecnt:    likecnt,
		Collectcnt: collectcnt,
	}
}

// 简化的日志函数
type gormLoggerFunc func(msg string, fields ...interface{})

func (f gormLoggerFunc) Printf(msg string, args ...interface{}) {
	f(msg)
}
