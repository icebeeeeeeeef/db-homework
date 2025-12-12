package startup

import (
	"webook/internal/events"
	"webook/internal/repository"
	articlerepo "webook/internal/repository/article"
	articledao "webook/internal/repository/dao/article"
	"webook/internal/service"
	"webook/pkg/logger"

	"gorm.io/gorm"
)

// NewArticleServiceWithDeps 用于在本地装配链中显式消费 thirdProvider 的输出
// 以确保 *gorm.DB 和 logger.LoggerV1 被初始化（避免 provider 未使用）
func NewArticleServiceWithDeps(db *gorm.DB, l logger.LoggerV1, userRepo repository.UserRepository, producer events.Producer) service.ArticleService {
	// 这里先简单打个日志，确保依赖被真正使用
	l.Info("Init ArticleService with db and logger", logger.Field{Key: "db", Value: db})
	// 创建 ArticleRepository
	articleDAO := articledao.NewArticleDAO(db)
	repo := articlerepo.NewCachedArticleRepository(articleDAO, db, l, userRepo)
	// 构造 Service
	return service.NewArticleService(repo, l, producer)
}
