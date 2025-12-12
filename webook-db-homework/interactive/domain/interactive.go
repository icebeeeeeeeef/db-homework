package domain

type Interactive struct {
	Biz        string `json:"biz"`
	BizId      int64  `json:"biz_id"`
	Readcnt    int64  `json:"readcnt"`
	Likecnt    int64  `json:"likecnt"`
	Collectcnt int64  `json:"collectcnt"`
	Liked      bool   `json:"liked"`
	Collected  bool   `json:"collected"`
}

/*
type Collection struct {
	Name  string     `json:"name"`
	UID   int64      `json:"uid"`
	Items []Resource `json:"items"`
}
*/
