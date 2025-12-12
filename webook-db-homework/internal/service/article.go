package service

import (
	"context"
	"webook/internal/domain"
	events "webook/internal/events/article"
	repository "webook/internal/repository/article"
	"webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	Detail(ctx context.Context, id int64) (domain.Article, error)
	PubDetail(ctx context.Context, id int64, uid int64) (domain.Article, error)
	ListPub(ctx context.Context, offset int, limit int) ([]domain.Article, error)
}

type ArticleService_ struct {
	repo repository.ArticleRepository
	//authorRepo repository.ArticleAuthorRepository
	//readerRepo repository.ArticleReaderRepository
	l        logger.LoggerV1
	producer events.Producer
}

func NewArticleService(repo repository.ArticleRepository, l logger.LoggerV1, producer events.Producer) ArticleService {
	return &ArticleService_{
		repo:     repo,
		l:        l,
		producer: producer,
	}
}

func (s *ArticleService_) ListPub(ctx context.Context, offset int, limit int) ([]domain.Article, error) {
	return s.repo.ListPub(ctx, offset, limit) //这个和List是一样的
}
func (s *ArticleService_) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusDraft
	if article.ID == 0 {
		return s.repo.Create(ctx, article)
	} else {
		return article.ID, s.update(ctx, article)
	}
}

func (s *ArticleService_) update(ctx context.Context, article domain.Article) error {

	return s.repo.Update(ctx, article)
}

/*
	func (s *ArticleService_) Publish(ctx context.Context, article domain.Article) (int64, error) {
		var (
			err error
			id  = article.ID
		)
		if article.ID > 0 {
			err = s.authorRepo.Update(ctx, article)
		} else {
			id, err = s.authorRepo.Create(ctx, article)
		}
		if err != nil {
			return 0, err
		}
		//到这里是制作库保存完毕
		article.ID = id
		//保持制作库和线上库的id一致

		for i := 0; i < 3; i++ {
			err = s.readerRepo.Save(ctx, article)
			if err == nil {
				break
			}
			s.l.Error("线上库发布的部分重试失败，重试次数", logger.Int("retry", i+1), logger.Error(err))
		}
		if err != nil {
			s.l.Error("线上库发布全部重试失败", logger.Error(err))
			return 0, err
		}
		return article.ID, nil
	}
*/
func (s *ArticleService_) Publish(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusPublished
	return s.repo.Sync2(ctx, article)
}

func (s *ArticleService_) Withdraw(ctx context.Context, article domain.Article) error {
	article.Status = domain.ArticleStatusWithdraw
	return s.repo.SyncStatus(ctx, article)
}

func (s *ArticleService_) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return s.repo.List(ctx, uid, offset, limit)
}

func (s *ArticleService_) Detail(ctx context.Context, id int64) (domain.Article, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ArticleService_) PubDetail(ctx context.Context, id int64, uid int64) (domain.Article, error) {
	//在这里生产一个读事件

	art, err := s.repo.GetPubByID(ctx, id)
	if err == nil && s.producer != nil {
		//异步生产读事件
		go func() {

			//这里不要带上具体的参数，因为可能队列中的信息对应的字段早就被修改了，想要自己去查
			er := s.producer.ProduceReadEvent(ctx, events.ReadEvent{Uid: uid, Aid: id})
			if er != nil {
				s.l.Error("生产读事件失败", logger.Error(er))
			}
		}()
	}

	return art, err
}
