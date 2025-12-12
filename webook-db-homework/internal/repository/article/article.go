package repository

import (
	"context"
	"fmt"
	"webook/internal/domain"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	articledao "webook/internal/repository/dao/article"
	"webook/pkg/logger"

	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	//Sync(ctx context.Context, art domain.Article) (int64, error)
	Sync2(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, art domain.Article) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetPubByID(ctx context.Context, id int64) (domain.Article, error)
	ListPub(ctx context.Context, offset int, limit int) ([]domain.Article, error)
}

type ArticleRepository_ struct {
	dao   articledao.ArticleDAO
	db    *gorm.DB
	cache cache.ArticleCache
	//userepo 用于获取作者的名称信息
	userepo repository.UserRepository
	l       logger.LoggerV1
}

func NewCachedArticleRepository(dao articledao.ArticleDAO, db *gorm.DB, l logger.LoggerV1, userepo repository.UserRepository) ArticleRepository {
	return &ArticleRepository_{dao: dao, db: db, l: l, userepo: userepo}
}

// NewArticleRepositoryWithCache 显式设置文章缓存，便于在精简版中直接注入 redis 缓存
func NewArticleRepositoryWithCache(dao articledao.ArticleDAO, db *gorm.DB, l logger.LoggerV1, userepo repository.UserRepository, c cache.ArticleCache) ArticleRepository {
	return &ArticleRepository_{dao: dao, db: db, l: l, userepo: userepo, cache: c}
}

func (c *ArticleRepository_) ListPub(ctx context.Context, offset int, limit int) ([]domain.Article, error) {
	arts, err := c.dao.ListPub(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	res := make([]domain.Article, 0, len(arts))
	for _, a := range arts {
		d := c.toDomain(articledao.Article(a))
		if c.userepo != nil {
			if user, er := c.userepo.FindById(ctx, d.Author.ID); er == nil {
				d.Author.Name = user.Nickname
			}
		}
		res = append(res, d)
	}
	return res, nil
}
func (c *ArticleRepository_) GetPubByID(ctx context.Context, id int64) (domain.Article, error) {
	if c.cache != nil {
		art, err := c.cache.GetPub(ctx, id)
		if err == nil {
			return art, nil
		}
	}

	artdao, err := c.dao.GetPubByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := c.userepo.FindById(ctx, artdao.AuthorID)
	if err != nil {
		return domain.Article{}, err
	}
	//因为这里没有作者的名称信息，所以需要耦合userrepo插入作者姓名，在这里组装后再传上去
	data := domain.Article{
		ID:      artdao.ID,
		Title:   artdao.Title,
		Content: artdao.Content,
		Author: domain.Author{
			ID:   user.Id,
			Name: user.Nickname,
		},
		Status:    artdao.Status,
		CreatedAt: artdao.CreatedAt,
		UpdatedAt: artdao.UpdatedAt,
	}
	if c.cache != nil {
		defer func() {
			c.cache.SetPub(ctx, data)
		}()
	}
	return data, nil
}

func (c *ArticleRepository_) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	if c.cache != nil {
		art, err := c.cache.Get(ctx, id)
		if err == nil {
			return art, nil
		}
	}
	artdao, err := c.dao.GetByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	data := c.toDomain(artdao)
	if c.cache != nil {
		go func() {
			er := c.cache.Set(ctx, data)
			if er != nil {
				c.l.Error("缓存文章失败", logger.Error(er))
			}
		}()
	}
	return data, nil
}

func (c *ArticleRepository_) Create(ctx context.Context, art domain.Article) (int64, error) {
	if c.cache != nil {
		defer func() {
			c.cache.DelFirstPage(ctx, art.Author.ID)
		}()
	}

	return c.dao.Insert(ctx, articledao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorID: art.Author.ID,
		Status:   art.Status,
	})
}

func (c *ArticleRepository_) Update(ctx context.Context, art domain.Article) error {
	if c.cache != nil {
		defer func() {
			c.cache.DelFirstPage(ctx, art.Author.ID)
		}()
	}
	err := c.dao.Update(ctx, articledao.Article{
		ID:       art.ID,
		Title:    art.Title,
		Content:  art.Content,
		AuthorID: art.Author.ID,
		Status:   art.Status,
	})
	return err
}
func (c *ArticleRepository_) toEntity(art domain.Article) articledao.Article {
	return articledao.Article{
		ID:       art.ID,
		Title:    art.Title,
		Content:  art.Content,
		AuthorID: art.Author.ID,
		Status:   art.Status,
	}
}

/*
	func (c *ArticleRepository_) Sync(ctx context.Context, art domain.Article) (int64, error) {
		tx := c.db.WithContext(ctx).Begin() //在repo层面使用事务
		var (
			err error
			id  int64
		)
		if err = tx.Error; err != nil {
			return 0, err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
		entity := c.toEntity(art)
		//先操作制作库上的内容
		if art.ID > 0 {
			err = c.author.Update(ctx, entity)
		} else {
			id, err = c.author.Insert(ctx, entity)
		}
		if err != nil {
			return 0, err
		}
		entity.ID = id
		//再操作线上库上的内容
		err = c.reader.Upsert(ctx, articledao.ReaderArticle{
			Article: entity,
		})
		tx.Commit()
		return id, err
	}
*/
func (c *ArticleRepository_) Sync2(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Sync(ctx, c.toEntity(art))
	if err != nil {
		return 0, err
	}

	// 成功时才更新缓存
	if c.cache != nil {
		go func() {
			c.cache.DelFirstPage(ctx, art.Author.ID)
			c.cache.SetPub(ctx, art)
		}()
	}

	return id, nil
}

func (c *ArticleRepository_) SyncStatus(ctx context.Context, art domain.Article) error {
	if c.cache != nil {
		defer func() {
			c.cache.DelFirstPage(ctx, art.Author.ID)
			c.cache.DelPub(ctx, art.ID)
		}()
	}
	return c.dao.SyncStatus(ctx, c.toEntity(art))
}

func (c *ArticleRepository_) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// 我们可以在repo层组装我们的缓存策略
	if c.cache != nil && offset == 0 && limit <= 100 {
		articles, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			if len(articles) > 0 {
				go func(first domain.Article) {
					c.precache(ctx, first)
				}(articles[0])
			}
			return articles, nil
		}
	}
	// 如果缓存没有命中，则从数据库中查询

	articles, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}

	data := slice.Map[articledao.Article, domain.Article](articles, func(idx int, a articledao.Article) domain.Article {
		return c.toDomain(a)
	})
	if c.cache != nil {
		go func() {
			err := c.cache.SetFirstPage(ctx, uid, data)
			if err != nil {
				c.l.Error("缓存首页文章失败", logger.Error(err))
			}
			if len(data) > 0 {
				c.precache(ctx, data[0])
			}
		}()
	}

	return data, err
}

func (c *ArticleRepository_) toDomain(a articledao.Article) domain.Article {
	return domain.Article{
		ID:      a.ID,
		Title:   a.Title,
		Content: a.Content,
		Author: domain.Author{
			ID: a.AuthorID,
		},
		Status:    a.Status,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (c *ArticleRepository_) precache(ctx context.Context, articles domain.Article) error {
	if c.cache == nil {
		return nil
	}
	if len(articles.Content) > 1024*1024*10 {
		return nil
	}

	err := c.cache.Set(ctx, articles)
	if err != nil {
		return fmt.Errorf("预缓存文章失败: %w", err)
	}
	return nil

}
