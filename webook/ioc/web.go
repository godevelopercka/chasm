package ioc

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"practice/webook/internal/web"
	ijwt "practice/webook/internal/web/jwt"
	"practice/webook/internal/web/middleware"
	"practice/webook/pkg/ginx/middlewares/logger"
	logger2 "practice/webook/pkg/logger"
	"strings"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, userHdl *web.UserHandler, oauth2WechatHdl web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterUserRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable, jwtHdl ijwt.Handler, l logger2.LoggerV1) []gin.HandlerFunc {
	bd := logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
		l.Debug("HTTP请求", logger2.Field{Key: "al", Value: al})
	}).AllowReqBody(true).AllowRespBody()
	viper.OnConfigChange(func(in fsnotify.Event) {
		ok := viper.GetBool("web.logreq")
		bd.AllowReqBody(ok)
	})
	return []gin.HandlerFunc{
		corsHdl(),
		bd.Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/login").Build(),
		//ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"}, // 不加这个，前端是拿不到的
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "ckago.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
func InitWebServer5(funcs []gin.HandlerFunc,
	userHdl *web.UserHandler,
	//artHdl *web.ArticleHandler,
	oauth2Hdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(funcs...)
	// 注册路由
	userHdl.RegisterUserRoutes(server)
	//artHdl.RegisterUserRoutes(server)
	oauth2Hdl.RegisterRoutes(server)
	return server
}
