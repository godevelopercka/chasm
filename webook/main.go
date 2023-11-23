package main

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	//db := initDB()
	//rdb := initRedis()
	//server := initWebServer()
	//
	//u := initUser(db, rdb)
	//u.RegisterUserRoutes(server)

	initViper()
	//initViperV1()
	//initViperReader()
	//initViperRemote()
	initLogger()
	println(viper.AllKeys())
	setting := viper.AllSettings()
	fmt.Println(setting)
	server := InitWebServer5() // wire.go 的方法名

	//server := gin.Default() // 关闭 Redis和jwt 测试 k8s 用
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})
	server.Run(":9090")
}

// 解决跨域问题
//func initWebServer() *gin.Engine {
//	server := gin.Default()
//
//	// 结合 pkg 中的代码解读
//	// 使用 redis 限流
//	//redisClient := redis.NewClient(&redis.Options{
//	//	Addr: config.Config.Redis.Addr, // 结合 config
//	//}) // wrk 测试时不需要
//	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build()) // wrk 测试时不需要
//
//	server.Use(cors.New(cors.Config{
//		AllowOrigins: []string{"http://localhost:3000"},
//		//AllowMethods:     []string{"POST", "GET"},
//		AllowHeaders:     []string{"Content-Type", "Authorization"},
//		ExposeHeaders:    []string{"x-jwt-token"}, // 不加这个，前端是拿不到的
//		AllowCredentials: true,
//		AllowOriginFunc: func(origin string) bool {
//			if strings.HasPrefix(origin, "http://localhost") {
//				return true
//			}
//			return strings.Contains(origin, "ckago.com")
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//
//	// 保存 session
//	// 方法二：单实例部署 memstore
//	//store := memstore.NewStore([]byte("mH0fKRn8a5KPK0fSWItFyVrjrkbN9gWM"),
//	//	[]byte("6vMxvPNE90RFJQtcqK8QeGzth5MS1aiz"))
//	// 方法三: 多实例部署 redis
//	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("mH0fKRn8a5KPK0fSWItFyVrjrkbN9gWM"),
//	//	[]byte("6vMxvPNE90RFJQtcqK8QeGzth5MS1aiz"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	//store := memstore.NewStore([]byte("mH0fKRn8a5KPK0fSWItFyVrjrkbN9gWM"), []byte("6vMxvPNE90RFJQtcqK8QeGzth5MS1aiz")) // wrk 测试时不需要
//
//	// 方法一: cookie
//	//store := cookie.NewStore([]byte("secret"))
//	//server.Use(sessions.Sessions("mysession", store))
//	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())
//	// 方法四：jwt
//	//server.Use(sessions.Sessions("mysession", store)) // wrk 测试时不需要
//	server.Use(middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/signup").
//		IgnorePaths("/users/login_sms/code/send").IgnorePaths("/users/login_sms").
//		IgnorePaths("/users/login").Build())
//	return server
//}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("这是 replace 之前")
	// 如果你不 replace，直接用 zap.L(), 你啥都打不出来
	zap.ReplaceGlobals(logger)
	zap.L().Info("hello，你搞好了")

	type Demo struct {
		Name string `json:"name"`
	}

	zap.L().Info("这是实验参数",
		zap.Error(errors.New("这是一个 error")),
		zap.Int64("id", 123),
		zap.Any("一个结构体", Demo{Name: "hello"}))
}

//func initViperReader() {
//	viper.SetConfigType("yaml")
//	cfg :=
//db.mysql:
//	dsn: "root:root@tcp(localhost:13317)/webookv1"
//
//redis:
//	addr: "localhost:6379"
//
//	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
//	if err != nil {
//		panic(err)
//	}
//}

func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379",
		// 通过 webook 和其他使用 etcd 的区别出来
		"/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "指定配置文件路径") // 要在设置里面配置 program arguments : --config==config/dev.yaml
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉你文件变了，不能告诉你，文件的哪些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，它会在 in 里面告诉你变更前的数据，和变更后的数据
		// 更好的设计师，它会直接告诉你差异
		fmt.Println(in.Name, in.Op)
	})
	//viper.SetConfigFile("D:\\workspace\\go\\practice\\webook\\config\\dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViper() {
	//viper.SetDefault("db.mysql.dsn", "root:root@tcp(localhost:3306)/webookv1")
	// 配置文件的名字，但是不包含文件扩展名
	// 不包含 .go，.yaml 之类的后缀
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置用的是 yaml 格式
	// 现实中，有很多格式，JSON，XML，YAML，TOML，ini
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录
	viper.AddConfigPath("./webook/config")
	//viper.AddConfigPath("./tmp/config")
	//viper.AddConfigPath("./etc/webook")
	// 读取配置到 viper 里面，或者你可以理解为加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	otherViper := viper.New()
	otherViper.SetConfigName("myjson")
	otherViper.AddConfigPath("./config")
	otherViper.SetConfigType("json")
}
