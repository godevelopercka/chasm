//go:build k8s

// 使用 k8s 这个编译标签

package config

var Config = WebookConfig{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:3308)/webook",
	},
	Redis: RedisConfig{
		Addr:     "webook-redis:6380",
		Password: "",
		DB:       1,
	},
}
