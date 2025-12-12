package ioc

import (
	"fmt"

	rlock "github.com/gotomicro/redis-lock"
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

func InitRedisLock(client redis.Cmdable) *rlock.Client {
	return rlock.NewClient(client)
}
