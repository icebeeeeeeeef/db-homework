package async

import (
	"context"
	"webook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
}

func NewSMSService(svc sms.Service) sms.Service {
	return &SMSService{
		svc: svc,
	}
}

func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	return s.svc.Send(ctx, biz, args, numbers...)
}
