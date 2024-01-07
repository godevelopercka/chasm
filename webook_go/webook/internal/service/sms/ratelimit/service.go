package ratelimit

import (
	"context"
	"fmt"
	"webook_go/webook/internal/service/sms"
	"webook_go/webook/pkg/ratelimit"
)

// 字段式装饰器，当不想别人绕开你的接口时，使用这种方法
var errLimited = fmt.Errorf("触发了限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		// 系统错误
		// 可以限流：保守策略，你的下游很坑的时候
		// 可以不限：你的下游很强，业务可用性要求很高，尽量容错策略
		// 包一下这个错误
		return fmt.Errorf("短信服务判断是否限流出现问题,%w", err)
	}
	if limited {
		// 这个需要的时候再改成大写，公开错误
		return errLimited
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	// 这里也可以加新特性
	return err
}
