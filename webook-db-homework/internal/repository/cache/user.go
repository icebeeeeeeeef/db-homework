package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"webook/internal/domain"

	"github.com/redis/go-redis/v9"
)

var ErrUserNotFound = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id int64) error
	SetEmpty(ctx context.Context, id int64) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := c.key(id)
	val, err := c.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	if len(val) == 0 {
		return domain.User{}, ErrUserNotFound
	}
	var user domain.User
	err = json.Unmarshal(val, &user)
	return user, err
}

func (c *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := c.key(user.Id)
	return c.cmd.Set(ctx, key, val, c.expiration).Err()
}

func (c *RedisUserCache) Delete(ctx context.Context, id int64) error {
	key := c.key(id)
	return c.cmd.Del(ctx, key).Err()
}

// SetEmpty 缓存空值，防止缓存穿透
func (c *RedisUserCache) SetEmpty(ctx context.Context, id int64) error {
	key := c.key(id)
	// 使用较短的过期时间，避免长期占用缓存空间
	return c.cmd.Set(ctx, key, "", time.Minute*5).Err()
}

// IsEmpty 返回缓存是否为空值占位
func (c *RedisUserCache) IsEmpty(val string) bool {
	return val == ""
}

func (c *RedisUserCache) key(id int64) string { //统一的生成键名
	return fmt.Sprintf("webook:user:%d", id)
}
