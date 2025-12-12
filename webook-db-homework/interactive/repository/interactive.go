package repository

import (
	"context"
	"fmt"
	"log"
	"webook/interactive/domain"
	"webook/interactive/repository/cache"
	"webook/interactive/repository/dao"

	"github.com/redis/go-redis/v9"
)

type InteractiveRepository interface {
	IncLike(ctx context.Context, biz string, id int64, uid int64) error
	DecLike(ctx context.Context, biz string, id int64, uid int64) error
	IncRead(ctx context.Context, biz string, id int64) error
	BatchIncRead(ctx context.Context, biz []string, ids []int64) error
	IncCollect(ctx context.Context, biz string, id int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error)
}

type InteractiveRepository_ struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func NewInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache) InteractiveRepository {
	return &InteractiveRepository_{dao: dao, cache: cache}
}

func (r *InteractiveRepository_) GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) {
	inters, err := r.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.Interactive)
	for _, inter := range inters {
		res[inter.BizId] = domain.Interactive{
			Biz:        inter.Biz,
			BizId:      inter.BizId,
			Readcnt:    inter.Readcnt,
			Likecnt:    inter.Likecnt,
			Collectcnt: inter.Collectcnt,
		}
	}
	return res, nil
}
func (r *InteractiveRepository_) BatchIncRead(ctx context.Context, biz []string, ids []int64) error {

	/*err := r.dao.BatchIncRead(ctx, biz, ids)
	if err != nil {
		return err
	}
	return r.cache.IncrBatchReadIfPresent(ctx, biz, ids)*/
	return r.dao.BatchIncRead(ctx, biz, ids)
}

func (r *InteractiveRepository_) IncLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := r.dao.IncLike(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return r.cache.IncrLikeIfPresent(ctx, biz, id)
}

func (r *InteractiveRepository_) DecLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := r.dao.DecLike(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return r.cache.DecrLike(ctx, biz, id)
}

func (r *InteractiveRepository_) IncRead(ctx context.Context, biz string, id int64) error {
	err := r.dao.IncRead(ctx, biz, id)
	if err != nil {
		return err
	}
	fmt.Println("IncRead", biz, id)
	return r.cache.IncrReadIfPresent(ctx, biz, id)
}
func (r *InteractiveRepository_) IncCollect(ctx context.Context, biz string, id int64, uid int64) error {
	err := r.dao.IncCollect(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return r.cache.IncrCollectIfPresent(ctx, biz, id)
}

func (r *InteractiveRepository_) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	data, err := r.cache.Get(ctx, biz, id)
	if err == nil {
		return data, nil
	}

	if err != redis.Nil {
		return domain.Interactive{}, err
	}

	// 缓存未命中，查数据库
	info, err := r.dao.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}

	data = domain.Interactive{
		Biz:        biz,
		BizId:      id,
		Readcnt:    info.Readcnt,
		Likecnt:    info.Likecnt,
		Collectcnt: info.Collectcnt,
	}

	// 异步设置缓存
	go func() {
		er := r.cache.Set(context.Background(), biz, id, data)
		if er != nil {
			log.Println("set cache error", er)
		}
	}()

	return data, nil
}
func (r *InteractiveRepository_) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (bool, error) {

	_, err := r.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}

}
func (r *InteractiveRepository_) GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := r.dao.GetCollectInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}
