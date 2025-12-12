package service

import (
	"context"

	intrdao "webook/interactive/repository/dao"
	"webook/internal/domain"

	"gorm.io/gorm"
)

// DBInteractiveService 持久化互动数据到 MySQL，并在每次操作后返回最新互动信息
type DBInteractiveService struct {
	dao intrdao.InteractiveDAO
}

func NewDBInteractiveService(db *gorm.DB) InteractiveService {
	return &DBInteractiveService{
		dao: intrdao.NewInteractiveDAO(db),
	}
}

func (s *DBInteractiveService) Like(ctx context.Context, biz string, id int64, uid int64, like bool) (domain.Interactive, error) {
	var err error
	if like {
		err = s.dao.IncLike(ctx, biz, id, uid)
	} else {
		err = s.dao.DecLike(ctx, biz, id, uid)
	}
	if err != nil {
		return domain.Interactive{}, err
	}
	return s.Get(ctx, biz, id, uid)
}

func (s *DBInteractiveService) Collect(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	if err := s.dao.IncCollect(ctx, biz, id, uid); err != nil {
		return domain.Interactive{}, err
	}
	return s.Get(ctx, biz, id, uid)
}

func (s *DBInteractiveService) Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	info, err := s.dao.Get(ctx, biz, id)
	if err == intrdao.ErrRecordNotFound {
		// 没有互动记录也要返回默认值
		return domain.Interactive{}, nil
	}
	if err != nil {
		return domain.Interactive{}, err
	}
	liked, errLike := s.dao.GetLikeInfo(ctx, biz, id, uid)
	if errLike != nil && errLike != intrdao.ErrRecordNotFound {
		return domain.Interactive{}, errLike
	}
	collected, errCollect := s.dao.GetCollectInfo(ctx, biz, id, uid)
	if errCollect != nil && errCollect != intrdao.ErrRecordNotFound {
		return domain.Interactive{}, errCollect
	}
	isLiked := errLike == nil && liked.Status
	isCollected := errCollect == nil && collected.ID > 0
	return domain.Interactive{
		ReadCnt:    info.Readcnt,
		LikeCnt:    info.Likecnt,
		CollectCnt: info.Collectcnt,
		Liked:      isLiked,
		Collected:  isCollected,
	}, nil
}

func (s *DBInteractiveService) IncrRead(ctx context.Context, biz string, id int64) error {
	return s.dao.IncRead(ctx, biz, id)
}
