package ratelimit

import "context"

type Limiter interface {
	// Limited 有没有触发限流。 key 就是限流对象
	// bool 代表是否限流，true 就是要限流
	// err 代表限流器本身是否错误
	Limit(ctx context.Context, key string) (bool, error)
}
