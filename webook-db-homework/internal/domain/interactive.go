package domain

// Interactive 简化版互动数据，存储在 Redis 中
type Interactive struct {
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Liked      bool
	Collected  bool
}
