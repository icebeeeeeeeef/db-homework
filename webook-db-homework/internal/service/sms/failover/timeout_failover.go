package failover

import (
	"context"
	"sync/atomic"
	"webook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	//超时次数
	cnt int32
	//阈值，超过该次数触发切换服务商
	threshold int32
	// 服务商索引
	idx int32
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	cnt := atomic.LoadInt32(&t.cnt)
	idx := atomic.LoadInt32(&t.idx)
	//如果超过阈值
	if cnt > t.threshold {
		newidx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newidx) { //对比内存地址上的值和刚获取的值，相同就把newidx写入
			//这种操作更加轻量化，如果发现idx没有更新，则不进行重置
			atomic.StoreInt32(&t.cnt, 0)
		}
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, biz, args, numbers...)
	switch err {
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:
		//其他异常，不重试
	}
	return err
}

//其实这里不是严格的超过n个就切换，而是近似
//因为有可能第一个客户端刚刚判断完cnt超过阈值，还没来得及切换，第二个客户端就判断了，所以cnt会重置
//这样的话相当于第一个客户端的超时没有被记录，因为刚记录就被第二个重置了，这并不是严格的n个，但是无伤大雅
