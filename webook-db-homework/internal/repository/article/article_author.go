package repository

import (
	"context"
	"webook/internal/domain"
	articledao "webook/internal/repository/dao/article"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, author domain.Article) (int64, error)
	Update(ctx context.Context, author domain.Article) error
}

type ArticleAuthorRepository_ struct {
	dao articledao.ArticleDAO
}

func NewCachedArticleAuthorRepository(dao articledao.ArticleDAO) ArticleAuthorRepository {
	return &ArticleAuthorRepository_{dao: dao}
}

func (r *ArticleAuthorRepository_) Create(ctx context.Context, author domain.Article) (int64, error) {
	return 1, nil
}

func (r *ArticleAuthorRepository_) Update(ctx context.Context, author domain.Article) error {
	return nil
}
