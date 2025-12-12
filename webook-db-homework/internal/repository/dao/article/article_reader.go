package dao

import (
	"context"

	"gorm.io/gorm"
)

type ArticleReaderDAO interface {
	Upsert(ctx context.Context, article ReaderArticle) error
}

type GORMArticleReaderDAO struct {
	db *gorm.DB
}

func NewArticleReaderDAO(db *gorm.DB) ArticleReaderDAO {
	return &GORMArticleReaderDAO{db: db}
}

func (dao *GORMArticleReaderDAO) Upsert(ctx context.Context, article ReaderArticle) error {
	return dao.db.WithContext(ctx).Model(&ReaderArticle{}).Where("id = ?", article.ID).Updates(article).Error
}
