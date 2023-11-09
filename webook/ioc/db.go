package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"practice/webook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	//dsn := viper.GetString("db.mysql.dsn")
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(localhost:13317)/webookv1", // 设置 DSN 的默认值，如果 yaml 文件中有则会覆盖默认值
	}
	// 看起来，remote 不支持 key 的切割
	err := viper.UnmarshalKey("db", &cfg)
	//dsn := viper.GetString("db.mysql.dsn")
	//if err != nil {
	//	panic(err)
	//}
	db, err := gorm.Open(mysql.Open(cfg.DSN)) // 结合 config
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
