package failover

import (
	"context"
	"webook/internal/service/sms"
)

type FailoverSMSService struct {
	svcs []sms.Service
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

func (f *FailoverSMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {

}
