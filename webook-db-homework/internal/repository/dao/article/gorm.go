package dao

import (
	"context"
	"errors"
	"time"

	"webook/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{db: db}
}

func (dao *GORMArticleDAO) GetByID(ctx context.Context, id int64) (Article, error) {
	var article Article
	err := dao.db.WithContext(ctx).Model(&Article{}).Where("id = ?", id).First(&article).Error
	return article, err
}

func (dao *GORMArticleDAO) GetPubByID(ctx context.Context, id int64) (ReaderArticle, error) {
	var article ReaderArticle
	err := dao.db.WithContext(ctx).
		Model(&ReaderArticle{}).
		Where("id = ? AND status = ?", id, domain.ArticleStatusPublished).
		First(&article).Error
	return article, err
}

func (dao *GORMArticleDAO) ListPub(ctx context.Context, offset int, limit int) ([]ReaderArticle, error) {
	var articles []ReaderArticle
	err := dao.db.WithContext(ctx).
		Model(&ReaderArticle{}).
		Where("status = ?", domain.ArticleStatusPublished).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error
	return articles, err
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CreatedAt = now
	article.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&article).Error
	return article.ID, err

}

func (dao *GORMArticleDAO) Update(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.UpdatedAt = now

	//Updates方法会根据结构体的非零值进行更新,这是依赖于gorm忽略零值的特性，会用主键进行更新
	res := dao.db.WithContext(ctx).Model(&article).Where("id = ? AND author_id = ?", article.ID, article.AuthorID).Updates(map[string]any{
		"title":      article.Title,
		"content":    article.Content,
		"updated_at": now,
		"status":     article.Status,
		//在这里显示的把更新的字段都写出来，而不是利用gorm的特性更新，这样代码更易读
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("修改文章失败，可能是文章不存在或不是自己的文章")
	}
	return nil
}

func (dao *GORMArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	var (
		err error
		id  int64
	)
	id = article.ID
	err = dao.db.Transaction(func(tx *gorm.DB) error {

		if err = tx.Error; err != nil {
			return err
		}
		if article.ID > 0 {
			err = dao.Update(ctx, article)
		} else {
			id, err = dao.Insert(ctx, article)
		}
		if err != nil {
			return err
		}
		article.ID = id
		id, err = dao.Upsert(ctx, ReaderArticle(article))
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (dao *GORMArticleDAO) Upsert(ctx context.Context, article ReaderArticle) (int64, error) {
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		// ID 冲突的时候。实际上，在 MYSQL 里面你写不写都可以
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":      article.Title,
			"content":    article.Content,
			"updated_at": now,
			"status":     article.Status,
		}),
	}).Create(&article).Error
	return article.ID, err
}

func (dao *GORMArticleDAO) SyncStatus(ctx context.Context, article Article) error {
	article.UpdatedAt = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? AND author_id = ?", article.ID, article.AuthorID).
			Updates(map[string]any{
				"status":     article.Status,
				"updated_at": article.UpdatedAt,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("修改文章失败，可能是文章不存在或不是自己的文章")
		}

		res = tx.Model(&ReaderArticle{}).
			Where("id = ?", article.ID).
			Updates(map[string]any{
				"status":     article.Status,
				"updated_at": article.UpdatedAt,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("修改文章失败，可能是文章不存在或不是自己的文章")
		}
		return nil
	})
}

func (dao *GORMArticleDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	var articles []Article
	err := dao.db.WithContext(ctx).Model(&Article{}).Where("author_id = ?", uid).Offset(offset).Limit(limit).Order("updated_at DESC").Find(&articles).Error
	return articles, err
}
