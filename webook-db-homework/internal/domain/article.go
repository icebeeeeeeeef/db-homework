package domain

const (
	ArticleStatusUnknown = iota
	ArticleStatusDraft
	ArticleStatusWithdraw
	ArticleStatusPublished
	ArticleStatusDeleted
)

type Article struct {
	ID        int64
	Title     string
	Content   string
	Author    Author
	Status    uint8
	CreatedAt int64
	UpdatedAt int64
}

func (a *Article) IsUnknown() bool {
	return a.Status == ArticleStatusUnknown
}
func (a *Article) IsValid() bool {
	return a.Status == ArticleStatusDraft || a.Status == ArticleStatusWithdraw || a.Status == ArticleStatusPublished
}

func (a *Article) IsDraft() bool {
	return a.Status == ArticleStatusDraft
}

func (a *Article) IsWithdraw() bool {
	return a.Status == ArticleStatusWithdraw
}

func (a *Article) IsPublished() bool {
	return a.Status == ArticleStatusPublished
}

func (a *Article) IsDeleted() bool {
	return a.Status == ArticleStatusDeleted
}

type Author struct {
	ID   int64
	Name string
}

func (a *Article) GenAbstract() string {
	//这里用rune来处理，因为汉字占两个字符，用byte会出问题
	s := []rune(a.Content)
	if len(s) > 100 {
		return string(s[:100])
	}
	return string(s)

}
