//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook_go/webook/internal/repository"
	"webook_go/webook/internal/repository/article"
	"webook_go/webook/internal/repository/cache"
	"webook_go/webook/internal/repository/dao"
	article2 "webook_go/webook/internal/repository/dao/article"
	"webook_go/webook/internal/service"
	"webook_go/webook/internal/web"
	ijwt "webook_go/webook/internal/web/jwt"
	"webook_go/webook/ioc"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	cache.NewCodeCache,
	repository.NewUserRepository,
	service.NewUserService)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		//articlSvcProvider,
		article2.NewGORMArticleDAO,
		repository.NewCodeRepository,
		article.NewArticleRepository,
		// service 部分
		// 集成测试我们显示指定使用内存实现
		ioc.InitSMSService,

		// 指定啥也不干的 wechat service
		service.NewCodeService,
		service.NewArticleService,
		// 直接基于内存的实现
		ioc.InitWechatService,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		InitWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,
		// 你中间件呢
		// 你注册路由呢
		// 你这个地方没有用到前面的任何东西
		//gin.Default,

		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}

func InitArticleHandler(dao article.ArticleDAO) *web.ArticleHandler {
	wire.Build(thirdProvider,
		article2.NewGORMArticleDAO,
		service.NewArticleService,
		web.NewArticleHandler,
		article.NewArticleRepository)
	return &web.ArticleHandler{}
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil, nil)
}

func InitJwtHdl() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
