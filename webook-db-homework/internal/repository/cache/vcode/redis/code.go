package redis

import (
	"context"
	_ "embed"
	"fmt"
	"webook/internal/repository/cache/vcode"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/varify_code.lua
var luaVarifyCode string

type RedisCodeCache struct {
	client redis.Cmdable
	cache  vcode.CodeCache
}

func NewRedisCodeCache(client redis.Cmdable) vcode.CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code, 600).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return vcode.ErrSetCodeBusy
	default:
		return vcode.ErrSetCodeSystemError
	}
}

func (c *RedisCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz string, phone string, code string) error {
	res, err := c.client.Eval(ctx, luaVarifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return vcode.ErrVarifyCodeTooMany
	case -2:
		return vcode.ErrVarifyCodeInvalid
	default:
		return vcode.ErrVarifyCodeInvalid
	}
}
