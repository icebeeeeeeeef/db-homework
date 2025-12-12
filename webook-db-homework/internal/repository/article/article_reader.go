package repository

import (
	"context"
	"webook/internal/domain"
	articledao "webook/internal/repository/dao/article"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, article domain.Article) error
}

func NewArticleReaderRepository(dao articledao.ArticleDAO) ArticleReaderRepository {
	return &ArticleReaderRepository_{dao: dao}
}

type ArticleReaderRepository_ struct {
	dao articledao.ArticleDAO
}

//这里的Save是upsert语义，有则修改没有则创建

func (r *ArticleReaderRepository_) Save(ctx context.Context, article domain.Article) error {
	return nil
}
