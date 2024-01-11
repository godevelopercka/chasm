package dao

import (
	"gorm.io/gorm"
	"webook_go/webook/internal/repository/dao/article"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &article.Article{}, &article.PublishedArticle{}) // 想新建表，就往后面加
}
