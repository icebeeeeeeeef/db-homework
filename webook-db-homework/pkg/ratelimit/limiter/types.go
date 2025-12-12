package ratelimit

import "context"

type Limiter interface {
	Limit(ctx context.Context, key string) (bool, error)
	//其中key代表的是限流的对象
	//返回值代表是否有触发限流，如果触发限流，则返回true，否则返回false
	//err代表限流器内部错误
}
