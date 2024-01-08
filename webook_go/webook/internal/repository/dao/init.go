package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Article{}) // 想新建表，就往后面加
}
