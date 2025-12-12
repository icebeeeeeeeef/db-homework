package ioc

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	addr := viper.GetString("db.redis.addr")
	fmt.Println(addr)
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return redisClient
}
