package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"webook/internal/service/sms"
	ratelimit "webook/pkg/ratelimit/limiter"
)

var errLimited = errors.New("短信发送限流")

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

func (s *RatelimitSMSService) Send(ctx context.Context, biz string, args []string, number ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tecent")
	if err != nil {
		return fmt.Errorf("短信发送限流检查失败: %w", err)
	}
	if limited {
		return errLimited
	}

	//你可以在这里添加代码，新特性

	err = s.svc.Send(ctx, biz, args, number...)

	// 这里也可以添加代码，新特性
	return err
}
