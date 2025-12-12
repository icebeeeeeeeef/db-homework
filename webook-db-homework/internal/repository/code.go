package repository

import (
	"context"
	"webook/internal/repository/cache/vcode"
)

var (
	ErrCodeSendTooMany     = vcode.ErrSetCodeBusy
	ErrCodeSendSystemError = vcode.ErrSetCodeSystemError
	ErrCodeVerifyTooMany   = vcode.ErrVarifyCodeTooMany
	ErrCodeVerifyInvalid   = vcode.ErrVarifyCodeInvalid
)

type CodeRepository interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) error
}

type CachedCodeRepository struct {
	cache vcode.CodeCache
}

func NewCodeRepository(cache vcode.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: cache,
	}
}

func (r *CachedCodeRepository) Set(ctx context.Context, biz string, phone string, code string) error {
	return r.cache.Set(ctx, biz, phone, code)
}

func (r *CachedCodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return r.cache.Set(ctx, biz, phone, code)
}

func (r *CachedCodeRepository) Verify(ctx context.Context, biz string, phone string, code string) error {
	return r.cache.Verify(ctx, biz, phone, code)
}
