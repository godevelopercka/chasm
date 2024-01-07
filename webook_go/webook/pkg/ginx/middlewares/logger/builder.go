package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"io"
	"time"
)

// MiddlewareBuilder 注意点：
// 1. 小心日志内容过多。URL 可能很长，请求体，响应体都可能很大，你要考虑是不是完全输出到日志里面
// 2. 考虑 1 的问题，以及用户可能换用不同的日志框架，所以要有足够的灵活性
// 3. 考虑动态开关，结合监听配置文件，要小心并发安全
type MiddlewareBuilder struct {
	allowReqBody  *atomic.Bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc:   fn,
		allowReqBody: atomic.NewBool(false),
	}
}

func (b *MiddlewareBuilder) AllowReqBody(ok bool) *MiddlewareBuilder {
	b.allowReqBody.Store(ok)
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			// URL 本身也可能很长
			Url: url,
		}
		if b.allowReqBody.Load() && ctx.Request.Body != nil {
			// Body 读完就没有了，只能操作一次
			body, _ := ctx.GetRawData()
			reader := io.NopCloser(bytes.NewReader(body))
			ctx.Request.Body = reader
			//ctx.Request.GetBody = func() (io.ReadCloser, error) {
			//
			//}
			if len(body) > 1024 {
				body = body[:1024]
			}
			// 这其实是一个很消耗 CPU 和内存的操作
			// 因为这会引起赋值
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			ctx.Writer = responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()
			b.loggerFunc(ctx, al)
		}()

		// 执行到业务逻辑
		ctx.Next()
		//b.loggerFunc(ctx, al)
	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) WriterHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w responseWriter) Writer(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriterString(data string) (int, error) {
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	// HTTP 请求的方法
	Method string
	// Url 整个请求 URL
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}
