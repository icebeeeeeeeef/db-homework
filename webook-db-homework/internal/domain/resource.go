package domain

type Resource struct {
	BizId int64  `json:"biz_id"`
	Biz   string `json:"biz"`
}

const (
	ResourceArticle = "article"
	//此后遇到其他的业务可以在这里定义业务名称，方便在请求中传入参数
)
