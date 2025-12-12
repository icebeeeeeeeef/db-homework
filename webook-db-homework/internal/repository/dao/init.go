package dao

import (
	intrdao "webook/interactive/repository/dao"
	article "webook/internal/repository/dao/article"

	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&article.Article{},
		&article.ReaderArticle{},
		&intrdao.Interactive{},
		&intrdao.UserLikeSomething{},
		&intrdao.UserCollectSomething{},
	)
}
