package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook_go/webook/internal/web"
)

// 步骤三
// JWT 登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		// 我现在用 JWT 来校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.SplitN(tokenHeader, " ", 2) // 因为前端的 Authorization 是 Bear xxxx 的，中间有个空格
		if len(segs) != 2 {
			// 有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		println(claims)
		// ParseWithClaims 会将 claims 里面的 userId 解析出来，所以上一步要用指针传入 ParseWithClaims
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("NDIOaqI8vCUZfWoNVcol0CuqFwHbu4cn"), nil // 这里的加密算法一定要和登录功能中设置的一样，不然无效
		})
		if err != nil {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("claims", claims)
		// err 为 nil, token 不为 nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 你是要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 每十秒钟刷新一次
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("5M0jCqTP8bbDmWKn5wL1UYub7TOGhdP7"))
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
	}
}
