package retry

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"webook/internal/service/sms"
)

type RetrySMSService struct {
	svc         sms.Service
	maxRetries  int           // 最大重试次数（不含首发）
	baseDelay   time.Duration // 初始退避
	maxDelay    time.Duration // 最大退避
	jitterRatio float64       // 抖动比例 0~1
	retryableFn func(error) bool
}

func NewRetrySMSService(
	svc sms.Service,
	maxRetries int,
	baseDelay, maxDelay time.Duration,
	jitterRatio float64,
	retryableFn func(error) bool,
) *RetrySMSService {
	if retryableFn == nil {
		retryableFn = defaultRetryable
	}
	return &RetrySMSService{
		svc:         svc,
		maxRetries:  maxRetries,
		baseDelay:   baseDelay,
		maxDelay:    maxDelay,
		jitterRatio: jitterRatio,
		retryableFn: retryableFn,
	}
}

func (r *RetrySMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	var err error
	// 第一次直接发
	err = r.svc.Send(ctx, tplId, args, numbers...)
	if err == nil || !r.retryableFn(err) {
		return err
	}

	// 重试
	delay := r.baseDelay
	for i := 0; i < r.maxRetries; i++ {
		// 等待退避（尊重 ctx）
		if e := sleepCtx(ctx, addJitter(delay, r.jitterRatio)); e != nil {
			return err // 上一次的错误
		}
		// 再发
		err = r.svc.Send(ctx, tplId, args, numbers...)
		if err == nil || !r.retryableFn(err) {
			return err
		}
		// 指数退避
		delay *= 2
		if delay > r.maxDelay {
			delay = r.maxDelay
		}
	}
	return err
}

func defaultRetryable(err error) bool {
	// 例：对超时/临时网络错误/429/5xx 可重试；业务错误不重试
	var ne interface{ Temporary() bool }
	if errors.As(err, &ne) && ne.Temporary() {
		return true
	}
	// 可在此按错误码判断，比如:
	// if errors.Is(err, ErrTooManyRequests) || errors.Is(err, ErrTimeout) { return true }
	return false
}

func addJitter(d time.Duration, ratio float64) time.Duration {
	if ratio <= 0 {
		return d
	}
	j := int64(float64(d) * ratio)
	return d + time.Duration(rand.Int63n(j+1)) - time.Duration(j/2) // ±j/2 抖动
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

//最简单的重试策略就是连续重试固定的次数
