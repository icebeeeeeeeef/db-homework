package cache

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"webook/interactive/domain"

	"github.com/redis/go-redis/v9"
)

type InteractiveCache interface {
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, id int64, interactive domain.Interactive) error
	IncrReadIfPresent(ctx context.Context, biz string, id int64) error
	IncrLikeIfPresent(ctx context.Context, biz string, id int64) error
	IncrCollectIfPresent(ctx context.Context, biz string, id int64) error
	DecrLike(ctx context.Context, biz string, id int64) error
	DecrCollect(ctx context.Context, biz string, id int64) error
	//IncrBatchReadIfPresent(ctx context.Context, biz []string, ids []int64) error
}

const (
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"
	fieldLikeCnt    = "like_cnt"
)

var (
	//go:embed lua/incr_interactive.lua
	luaIncrCnt string
)

type InteractiveCache_ struct {
	client redis.Cmdable
}

func NewInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &InteractiveCache_{client: client}
}

/*
func (c *InteractiveCache_) IncrBatchReadIfPresent(ctx context.Context, biz []string, ids []int64) error {
	return c.client.Eval(ctx, luaIncrCnt, []string{c.key(biz, id)}, fieldReadCnt, 1).Err()
}*/

func (c *InteractiveCache_) key(biz string, id int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, id)
}

func (c *InteractiveCache_) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	data, err := c.client.HGetAll(ctx, c.key(biz, id)).Result() //hgetall 返回的是一个map
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(data) == 0 {
		return domain.Interactive{}, redis.Nil
	}

	readcnt, err := strconv.ParseInt(data[fieldReadCnt], 10, 64)
	if err != nil {
		return domain.Interactive{}, fmt.Errorf("解析 readcnt 失败: %w", err)
	}

	likecnt, err := strconv.ParseInt(data[fieldLikeCnt], 10, 64)
	if err != nil {
		return domain.Interactive{}, fmt.Errorf("解析 likecnt 失败: %w", err)
	}

	collectcnt, err := strconv.ParseInt(data[fieldCollectCnt], 10, 64)
	if err != nil {
		return domain.Interactive{}, fmt.Errorf("解析 collectcnt 失败: %w", err)
	}

	return domain.Interactive{
		Readcnt:    readcnt,
		Likecnt:    likecnt,
		Collectcnt: collectcnt,
	}, nil
}

func (c *InteractiveCache_) Set(ctx context.Context, biz string, id int64, interactive domain.Interactive) error {
	key := c.key(biz, id)
	return c.client.HSet(ctx, key, map[string]interface{}{
		fieldReadCnt:    interactive.Readcnt,
		fieldLikeCnt:    interactive.Likecnt,
		fieldCollectCnt: interactive.Collectcnt,
	}).Err()
}

func (c *InteractiveCache_) IncrReadIfPresent(ctx context.Context, biz string, id int64) error {
	return c.client.Eval(ctx, luaIncrCnt, []string{c.key(biz, id)}, fieldReadCnt, 1).Err()
}

func (c *InteractiveCache_) IncrLikeIfPresent(ctx context.Context, biz string, id int64) error {
	return c.client.Eval(ctx, luaIncrCnt, []string{c.key(biz, id)}, fieldLikeCnt, 1).Err()
}

func (c *InteractiveCache_) IncrCollectIfPresent(ctx context.Context, biz string, id int64) error {
	return c.client.Eval(ctx, luaIncrCnt, []string{c.key(biz, id)}, fieldCollectCnt, 1).Err()
}

func (c *InteractiveCache_) DecrLike(ctx context.Context, biz string, id int64) error {
	return c.client.Eval(ctx, luaIncrCnt, []string{c.key(biz, id)}, fieldLikeCnt, -1).Err()
}

func (c *InteractiveCache_) DecrCollect(ctx context.Context, biz string, id int64) error {
	return c.client.Eval(ctx, luaIncrCnt, []string{c.key(biz, id)}, fieldCollectCnt, -1).Err()
}
