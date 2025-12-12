package web

import "webook/internal/domain"

type ArticleVO struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Abstract   string `json:"abstract"`
	Content    string `json:"content"`
	Author     string `json:"author"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"updatedAt"`
	Status     uint8  `json:"status"`
	Liked      bool   `json:"liked"`
	Collected  bool   `json:"collected"`
	ReadCnt    int64  `json:"readCnt"`
	LikeCnt    int64  `json:"likeCnt"`
	CollectCnt int64  `json:"collectCnt"`
}

type ArticleReq struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type LikeReq struct {
	ID   int64 `json:"id"`
	Like bool  `json:"like"`
}

func (r ArticleReq) ToDomain(uid int64) domain.Article {
	return domain.Article{
		ID:      r.ID,
		Title:   r.Title,
		Content: r.Content,
		Author: domain.Author{
			ID: uid,
		},
	}
}
