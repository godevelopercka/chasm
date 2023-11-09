//go:build k8s

// 使用 k8s 这个编译标签

package config

var Config = WebookConfig{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-live:11313)/webookv1",
		//DSN: "root:root@tcp(localhost:13317)/webookv1",
	},
	Redis: RedisConfig{
		Addr:     "webook-live-redis:6379",
		Password: "",
		DB:       1,
	},
}
