package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
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
	server := InitWebServer()
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
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
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		IgnorePaths("/users/login").Build())
	return server
}
