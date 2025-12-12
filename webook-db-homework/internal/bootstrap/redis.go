package bootstrap

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	addr := viper.GetString("db.redis.addr")
	fmt.Println("redis addr:", addr)
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
