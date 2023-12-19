//go:build !k8s

// 没有 k8s 这个编译标签
// 开发环境 go:build dev
// 测试环境 go:build test
// 生产环境 go:build e2e

package config

var Config = WebookConfig{
	DB: DBConfig{
		// 本地连接
		DSN: "root:root@tcp(localhost:13316)/webook",
	},
	Redis: RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	},
}
