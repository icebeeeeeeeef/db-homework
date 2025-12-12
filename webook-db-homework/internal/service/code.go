package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

var (
	ErrCodeSendTooMany     = repository.ErrCodeSendTooMany
	ErrCodeSendSystemError = repository.ErrCodeSendSystemError
	ErrCodeVerifyTooMany   = repository.ErrCodeVerifyTooMany
	ErrCodeVerifyInvalid   = repository.ErrCodeVerifyInvalid
	codetplId              = "1234567890"
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) error
}

type CodeService_ struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &CodeService_{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService_) Send(ctx context.Context, biz string, phone string) error {
	//先生成验证码
	num := svc.generateCode()
	//存到redis中
	err := svc.repo.Set(ctx, "login", phone, num)
	if err != nil {
		return err
	}

	// 发送验证码
	err = svc.smsSvc.Send(ctx, codetplId, []string{num}, phone)
	if err != nil {
		//这里要不要直接返回错误，需要把redis中的key删除掉吗
		//由于err可能是超时连接的错误，因此实际上不需要删除key
		//这里可以进行重试，重试指定次数，可以写一个retryService类来异步重试设定好的次数
		//可以写一个retryService类来异步重试设定好的次数，在函数的retry函数中回调send函数
		return err
	}

	return nil
}
func (svc *CodeService_) Verify(ctx context.Context, biz string, phone string, inputCode string) error {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService_) generateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
