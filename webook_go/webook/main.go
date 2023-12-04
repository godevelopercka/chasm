package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook_go/webook/internal/repository"
	"webook_go/webook/internal/repository/dao"
	"webook_go/webook/internal/service"
	"webook_go/webook/internal/web"
	"webook_go/webook/internal/web/middleware"
)

func main() {
	db := InitDB()
	server := InitWebServer()
	u := InitUser(db)
	u.RegisterRoutes(server)
	server.Run(":8080")
}

func InitWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "ckago.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	// 步骤一：首先把 session 塞进 context 里面，相当于用 context 作为中介存储了 cookie，并将 session 命名成 mysession
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store)) // userId 放在 store

	// 步骤四
	// 登录校验
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())
	return server
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
