package ratelimit

import (
	"context"
	"fmt"
	"webook_go/webook/internal/service/sms"
	"webook_go/webook/pkg/ratelimit"
)

// 组合式装饰器，当接口有很多方法，但只需要装饰其中一个方法可以用这个
type RatelimitSMSServiceV1 struct {
	sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSServiceV1(limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSServiceV1{
		limiter: limiter,
	}
}

func (s RatelimitSMSServiceV1) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
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
	err = s.Service.Send(ctx, tpl, args, numbers...)
	// 这里也可以加新特性
	return err
}
