package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
	"webook_go/webook/internal/repository/dao"
	"webook_go/webook/pkg/logger"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:13316)/webook",
	}
	//dsn := viper.GetString("db.mysql.dsn")
	// remote 看起来不支持 key 的切割
	err := viper.UnmarshalKey("db", &cfg)
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		// 缺了一个 writer
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值，只有执行时间超过这个阈值，才会使用
			// 50ms, 100ms
			// SQL 查询必然要求命中索引，最好就是走一次磁盘 IO
			// 一次磁盘 IO 是不到 10ms
			SlowThreshold:             time.Millisecond * 10,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  glogger.Info,
		}),
	}) // 将 localhost 改成用 k8s-mysql-service 的 name,端口用 port
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}
