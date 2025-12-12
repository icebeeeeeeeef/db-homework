package vcode

import (
	"context"
	"errors"
)

var (
	ErrSetCodeBusy        = errors.New("发送验证码太频繁")
	ErrSetCodeSystemError = errors.New("系统错误")
	ErrVarifyCodeTooMany  = errors.New("验证码错误次数过多")
	ErrVarifyCodeInvalid  = errors.New("验证码无效")
)

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) error
}
