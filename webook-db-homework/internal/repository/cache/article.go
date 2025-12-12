package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"webook/internal/domain"

	"github.com/redis/go-redis/v9"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, articles []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, article domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, article domain.Article) error
	DelPub(ctx context.Context, id int64) error
}
type RedisArticleCache struct {
	client redis.Cmdable
}

func NewRedisArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{client: client}
}

func (c *RedisArticleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.client.Get(ctx, c.authorkey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var article domain.Article
	err = json.Unmarshal(data, &article)
	return article, err
}

func (c *RedisArticleCache) Set(ctx context.Context, article domain.Article) error {
	val, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.authorkey(article.ID), val, time.Second*15).Err()
}

func (c *RedisArticleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.client.Get(ctx, c.readerkey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var article domain.Article
	err = json.Unmarshal(data, &article)
	return article, err
}

func (c *RedisArticleCache) SetPub(ctx context.Context, article domain.Article) error {
	val, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.readerkey(article.ID), val, time.Minute*15).Err()
	//为什么和设置创作者的时间不同，因为这里是读者库，发布后会有很多读者阅读，所以需要更长的时间
}

func (c *RedisArticleCache) DelPub(ctx context.Context, id int64) error {
	return c.client.Del(ctx, c.readerkey(id)).Err()
}
func (c *RedisArticleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	data, err := c.client.Get(ctx, c.firstpagekey(uid)).Bytes()
	if err != nil {
		return nil, err
	}
	var articles []domain.Article
	err = json.Unmarshal(data, &articles)
	return articles, err
}

func (c *RedisArticleCache) SetFirstPage(ctx context.Context, uid int64, articles []domain.Article) error {
	for _, article := range articles {
		//只缓存摘要的部分
		article.Content = article.GenAbstract()
	}
	val := c.firstpageval(articles)
	return c.client.Set(ctx, c.firstpagekey(uid), val, time.Minute*15).Err()
}

func (c *RedisArticleCache) DelFirstPage(ctx context.Context, uid int64) error {
	return c.client.Del(ctx, c.firstpagekey(uid)).Err()
}

func (c *RedisArticleCache) authorkey(id int64) string {
	return fmt.Sprintf("article_author:%d", id)
}

func (c *RedisArticleCache) firstpagekey(uid int64) string {
	return fmt.Sprintf("article:firstpage:%d", uid)
}

func (c *RedisArticleCache) firstpageval(articles []domain.Article) string {
	val, err := json.Marshal(articles)
	if err != nil {
		return ""
	}
	return string(val)

}

func (c *RedisArticleCache) readerkey(id int64) string {
	return fmt.Sprintf("article_reader:%d", id)
}
