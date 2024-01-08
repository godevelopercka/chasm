package main

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
	"webook_go/webook/internal/integration/startup"
	ijwt "webook_go/webook/internal/web/jwt"
	"webook_go/webook/internal/web/middleware"
)

func main() {
	//db := InitDB()
	//rdb := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//u := InitUser(db, rdb)
	//server := InitWebServer()
	//u.RegisterRoutes(server)
	initViper()
	initLogger()
	//initViperRemote()
	keys := viper.AllKeys()
	println(keys)
	server := startup.InitWebServer()
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})
	server.Run(":8080")
}

type Demo struct {
	Name string
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("这是 replace 之前")
	// 如果你不 replace, 直接用 zap.L(), 你啥都打不出来
	zap.ReplaceGlobals(logger)
	zap.L().Info("hello, 你搞好了")
	zap.L().Info("实验一波", zap.Error(errors.New("这是一个 error")),
		zap.Int64("id", 123), zap.Any("一个机构体", Demo{Name: "hello"}))
}

func initWebServer(jwtHdl ijwt.Handler) *gin.Engine {
	server := gin.Default()
	// 结合 pkg 中的代码解读
	// 使用 redis 限流
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr, // 将 localhost 改成用 k8s-redis-service 的 name,端口用 port
	//})
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "live.webook.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	// 步骤一：首先把 session 塞进 context 里面，相当于用 context 作为中介存储了 cookie，并将 session 命名成 mysession
	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("NDIOaqI8vCUZfWoNVcol0CuqFwHbu4cn"), []byte("VICKB7WKwidXBpPnzHqeiwTnWLDcuahY"))
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("NDIOaqI8vCUZfWoNVcol0CuqFwHbu4cn"), []byte("VICKB7WKwidXBpPnzHqeiwTnWLDcuahY"))
	//if err != nil {
	//	panic(err)
	//}

	//server.Use(sessions.Sessions("mysession", store)) // userId 放在 store

	// 步骤四
	// 登录校验
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())
	// JWT 登录校验
	server.Use(middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
		IgnorePaths("/users/signup").
		IgnorePaths("/users/refresh_token").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		IgnorePaths("/oauth2/wechat/authurl").
		IgnorePaths("/oauth2/wechat/callback").
		IgnorePaths("/users/login").Build())
	return server
}

func initViperRemote() {
	viper.SetConfigType("yaml")
	// 通过 webook 和其他使用 etcd 的区别出来
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initViper() {
	// 设置默认值
	//viper.SetDefault("db.mysql.dsn", "root:root@tcp(localhost:13316)/webook")
	// 配置文件的名字，但是不包含文件扩展名
	// 不包含 .go, .yaml 之类的后缀
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置用的是 yaml 格式
	// 现实中，有很多格式，JSON, XML, YAML, TOML, ini
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录, 可以有多个
	viper.AddConfigPath("./webook/config")
	//viper.AddConfigPath("/tmp/config")
	//viper.AddConfigPath("/etc/config")
	// 读取配置到 viper 里面，或者你可以理解为加载到内存里面
	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉你文件变了，但不能告诉你哪里变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 启用多个 viper
	//otherViper := viper.New()
	//otherViper.SetConfigName("myjson")
	//otherViper.AddConfigPath("./config")
	//otherViper.SetConfigType("json")
}
