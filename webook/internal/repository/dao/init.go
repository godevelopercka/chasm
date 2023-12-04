package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}) // 多个表就在 &User{} 后面加
}
