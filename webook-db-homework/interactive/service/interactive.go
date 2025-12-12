package service

import (
	"context"
	"webook/interactive/domain"
	"webook/interactive/repository"

	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	IncrReadIfPresent(ctx context.Context, biz string, id int64) error
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error)
}

type InteractiveService_ struct {
	repo repository.InteractiveRepository
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &InteractiveService_{
		repo: repo,
	}
}

func (svc *InteractiveService_) GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) {
	return svc.repo.GetByIds(ctx, biz, ids)
}
func (svc *InteractiveService_) IncrReadIfPresent(ctx context.Context, biz string, id int64) error {
	return svc.repo.IncRead(ctx, biz, id)
}

func (svc *InteractiveService_) Like(ctx context.Context, biz string, id int64, uid int64) error {
	return svc.repo.IncLike(ctx, biz, id, uid)
}

func (svc *InteractiveService_) CancelLike(ctx context.Context, biz string, id int64, uid int64) error {
	return svc.repo.DecLike(ctx, biz, id, uid)
}

func (svc *InteractiveService_) Collect(ctx context.Context, biz string, id int64, uid int64) error {
	return svc.repo.IncCollect(ctx, biz, id, uid)
}

func (svc *InteractiveService_) Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	var (
		eg          errgroup.Group
		liked       bool
		collected   bool
		interactive domain.Interactive
		err         error
	)
	eg.Go(func() error {
		liked, err = svc.repo.GetLikeInfo(ctx, biz, id, uid)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		collected, err = svc.repo.GetCollectInfo(ctx, biz, id, uid)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		interactive, err = svc.repo.Get(ctx, biz, id)
		if err != nil {
			return err
		}
		return nil
	})
	err = eg.Wait() //在这里等所有协程都执行完
	if err != nil {
		return domain.Interactive{}, err
	}
	interactive.Liked = liked
	interactive.Collected = collected
	return interactive, nil
}
