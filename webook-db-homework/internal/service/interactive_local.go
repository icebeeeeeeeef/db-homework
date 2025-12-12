package service

import (
	"context"
	"fmt"
	"webook/internal/domain"

	"github.com/redis/go-redis/v9"
)

type InteractiveService interface {
	Like(ctx context.Context, biz string, id int64, uid int64, like bool) (domain.Interactive, error)
	Collect(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	IncrRead(ctx context.Context, biz string, id int64) error
}

// RedisInteractiveService 使用 Redis 存储互动数据，避免数据库迁移
type RedisInteractiveService struct {
	cmd redis.Cmdable
}

func NewInteractiveService(cmd redis.Cmdable) InteractiveService {
	return &RedisInteractiveService{cmd: cmd}
}

func (s *RedisInteractiveService) likeSetKey(biz string, id int64) string {
	return fmt.Sprintf("inter:%s:like:%d", biz, id)
}
func (s *RedisInteractiveService) collectSetKey(biz string, id int64) string {
	return fmt.Sprintf("inter:%s:collect:%d", biz, id)
}
func (s *RedisInteractiveService) readKey(biz string, id int64) string {
	return fmt.Sprintf("inter:%s:read:%d", biz, id)
}
func (s *RedisInteractiveService) likeCntKey(biz string, id int64) string {
	return fmt.Sprintf("inter:%s:likecnt:%d", biz, id)
}
func (s *RedisInteractiveService) collectCntKey(biz string, id int64) string {
	return fmt.Sprintf("inter:%s:collectcnt:%d", biz, id)
}

func (s *RedisInteractiveService) Like(ctx context.Context, biz string, id int64, uid int64, like bool) (domain.Interactive, error) {
	likeSet := s.likeSetKey(biz, id)
	likeCntKey := s.likeCntKey(biz, id)
	if like {
		added, err := s.cmd.SAdd(ctx, likeSet, uid).Result()
		if err != nil {
			return domain.Interactive{}, err
		}
		if added > 0 {
			_ = s.cmd.Incr(ctx, likeCntKey).Err()
		}
	} else {
		removed, err := s.cmd.SRem(ctx, likeSet, uid).Result()
		if err != nil {
			return domain.Interactive{}, err
		}
		if removed > 0 {
			_ = s.cmd.Decr(ctx, likeCntKey).Err()
		}
	}
	return s.Get(ctx, biz, id, uid)
}

func (s *RedisInteractiveService) Collect(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	collectSet := s.collectSetKey(biz, id)
	collectCntKey := s.collectCntKey(biz, id)
	added, err := s.cmd.SAdd(ctx, collectSet, uid).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if added > 0 {
		_ = s.cmd.Incr(ctx, collectCntKey).Err()
	}
	return s.Get(ctx, biz, id, uid)
}

func (s *RedisInteractiveService) Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	liked, _ := s.cmd.SIsMember(ctx, s.likeSetKey(biz, id), uid).Result()
	collected, _ := s.cmd.SIsMember(ctx, s.collectSetKey(biz, id), uid).Result()
	likeCnt, _ := s.cmd.Get(ctx, s.likeCntKey(biz, id)).Int64()
	collectCnt, _ := s.cmd.Get(ctx, s.collectCntKey(biz, id)).Int64()
	readCnt, _ := s.cmd.Get(ctx, s.readKey(biz, id)).Int64()
	return domain.Interactive{
		ReadCnt:    readCnt,
		LikeCnt:    likeCnt,
		CollectCnt: collectCnt,
		Liked:      liked,
		Collected:  collected,
	}, nil
}

func (s *RedisInteractiveService) IncrRead(ctx context.Context, biz string, id int64) error {
	return s.cmd.Incr(ctx, s.readKey(biz, id)).Err()
}
