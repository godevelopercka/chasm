package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 登录校验模块

// 方法一(最好 )
// LoginMiddlewareBuilder 扩展性
type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

// 判断哪些路径不需要校验，做成切片格式
func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now()) // 刷新 session 通过 time 的话必须先注册这个才能用
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		// 我怎么知道，一分钟已经过去了？
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 重设一遍 session
		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60, // 设置 60 秒过期
		})
		now := time.Now()
		// 说明还没有刷新过，刚登陆，还没刷新过
		if updateTime == nil {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}
		// updateTime 是有的
		// 多久刷新一次
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > time.Second*10 { // 10秒刷新一次
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}
	}
}

// 方法二
var IgnorePaths []string

func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range IgnorePaths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 不需要登录校验的
		//if ctx.Request.URL.Path == "/users/login" ||
		//	ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

// 方法三
func CheckLoginV1(paths []string,
	abc int,
	bac int64,
	asdsd string) gin.HandlerFunc {
	if len(paths) == 0 {
		paths = []string{}
	}
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 不需要登录校验的
		//if ctx.Request.URL.Path == "/users/login" ||
		//	ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
