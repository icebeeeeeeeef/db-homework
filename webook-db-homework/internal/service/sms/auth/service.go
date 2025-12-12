package auth

import (
	"context"
	"errors"
	"webook/internal/service/sms"

	"github.com/golang-jwt/jwt/v5"
)

type SmsService struct {
	svc sms.Service
	key string
}

func NewSmsService(svc sms.Service, key string) sms.Service {
	return &SmsService{
		svc: svc,
		key: key,
	}
}
func (s *SmsService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {

	var tc Claims

	//这里biz中存的就是jwt token，必须解析后才能正常调用业务，这样我们就可以根据不同的业务来提供不同的模板
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil //传入校验签名的密钥
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token is invalid")
	}

	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}
