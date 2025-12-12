//这里我们可以在oauth包下面写钉钉，谷歌等等的其他的第三方账号登录

package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"webook/internal/domain"

	uuid "github.com/lithammer/shortuuid/v4"
)

var redirectURI = ""

type Service interface {
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatUser, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewService(appId string, appSecret string) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	// 微信扫码登录URL模板
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatUser, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	targetURL := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return domain.WechatUser{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatUser{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var result Result
	err = decoder.Decode(&result)
	if err != nil {
		return domain.WechatUser{}, err
	}
	if result.ErrCode != 0 {
		return domain.WechatUser{}, errors.New("获取微信用户信息失败")
	}

	return domain.WechatUser{
		OpenID:  result.OpenID,
		UnionID: result.UnionID,
	}, nil
}

type Result struct {
	ErrCode      int64  `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}
