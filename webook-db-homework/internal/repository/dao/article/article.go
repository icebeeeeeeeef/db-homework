package dao

import (
	"context"
)

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	Update(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, article ReaderArticle) (int64, error)
	SyncStatus(ctx context.Context, article Article) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error)
	GetByID(ctx context.Context, id int64) (Article, error)
	GetPubByID(ctx context.Context, id int64) (ReaderArticle, error)
	ListPub(ctx context.Context, offset int, limit int) ([]ReaderArticle, error)
}
