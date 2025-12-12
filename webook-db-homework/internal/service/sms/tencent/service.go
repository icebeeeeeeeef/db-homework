package tencent

import (
	"context"
	"fmt"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId     *string
	signature *string
	client    *sms.Client
}

func NewService(client *sms.Client, appId string, signature string) *Service {
	return &Service{
		client:    client,
		appId:     ekit.ToPtr[string](appId),
		signature: ekit.ToPtr[string](signature),
	}
}

func (s *Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signature
	req.TemplateId = ekit.ToPtr[string](biz)
	req.PhoneNumberSet = slice.Map[string, *string](numbers, func(idx int, num string) *string {
		return &num
	}) //把string 转换成string*指针，因为tencent sdk 中需要的是指针
	req.TemplateParamSet = slice.Map[string, *string](args, func(idx int, arg string) *string {
		return &arg
	})
	resp, err := s.client.SendSms(req) //调用发送的函数，返回响应和错误
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet { //遍历所有的手机号的发送状态，为空或者不为ok都说明发送失败
		if status.Code == nil || *status.Code != "Ok" {
			return fmt.Errorf("send sms failed, code: %s, message: %s", *status.Code, *status.Message)
		}
	}
	return err
}
