package dao

import (
	"context"

	"gorm.io/gorm"
)

type ArticleAuthorDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	Update(ctx context.Context, article Article) error
}

type GORMArticleAuthorDAO struct {
	db *gorm.DB
}

func NewArticleAuthorDAO(db *gorm.DB) ArticleAuthorDAO {
	return &GORMArticleAuthorDAO{db: db}
}

func (dao *GORMArticleAuthorDAO) Insert(ctx context.Context, article Article) (int64, error) {
	return dao.db.WithContext(ctx).Create(&article).RowsAffected, dao.db.WithContext(ctx).Create(&article).Error
}

func (dao *GORMArticleAuthorDAO) Update(ctx context.Context, article Article) error {
	return dao.db.WithContext(ctx).Model(&Article{}).Where("id = ?", article.ID).Updates(article).Error
}
