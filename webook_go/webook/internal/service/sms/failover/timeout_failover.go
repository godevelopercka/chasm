package failover

import (
	"context"
	"sync/atomic"
	"webook_go/webook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	// 你的服务商
	svcs []sms.Service
	idx  int32
	// 连续超时的个数
	cnt int32

	// 阈值
	// 连续超时超过这个次数，就要切换
	threshold int32
}

func (t TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		// 这里就要切换服务商，新的下标，往后挪了一个
		newIdx := (idx + 1) % int32(len(t.svcs)) // 取余是为了防止溢出
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 我成功往后挪了一位
			atomic.StoreInt32(&t.cnt, 0)
		}
		// else 就是出现并发，别人换成功了
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	case nil:
		// 你的连续状态被打断了
		atomic.StoreInt32(&t.cnt, 0)
	default:
		// 不知道什么错误，你可以考虑换一个：
		// - 超时，可能是偶发的，我尽量再试试
		// - 非超时，我直接下一个
		return err
	}
	return err
}

func NewTimeoutFailoverSMSService() sms.Service {
	return &TimeoutFailoverSMSService{}
}
