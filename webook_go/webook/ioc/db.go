package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook_go/webook/config"
	"webook_go/webook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN)) // 将 localhost 改成用 k8s-mysql-service 的 name,端口用 port
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
