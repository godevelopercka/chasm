package ioc

import (
	"github.com/redis/go-redis/v9"
	"practice/webook/internal/service/sms"
	"practice/webook/internal/service/sms/memory"
	"practice/webook/internal/service/sms/ratelimit"
	"practice/webook/internal/service/sms/retryable"
	limiter "practice/webook/pkg/ratelimit"
	"time"
)

func InitSMSService() sms.Service {
	// 换内存，还是换别的
	return memory.NewService()
}

func InitSMSServiceV1(cmd redis.Cmdable) sms.Service {
	// 换内存，还是换别的
	svc := ratelimit.NewRatelimitSMSService(memory.NewService(), limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100))
	return retryable.NewService(svc, 3)
}
