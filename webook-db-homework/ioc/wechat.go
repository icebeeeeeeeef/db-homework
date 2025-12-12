package ioc

import (
	"os"
	"webook/internal/service/oauth2/wechat"
	"webook/internal/web"
)

func InitOAuthWechatService() wechat.Service {
	appId := os.Getenv("WECHAT_APP_ID")
	if appId == "" {
		appId = "test_app_id" // 默认值，用于开发测试
	}
	appSecret := os.Getenv("WECHAT_APP_SECRET")
	if appSecret == "" {
		appSecret = "test_app_secret" // 默认值，用于开发测试
	}
	return wechat.NewService(appId, appSecret)
}

func NewWechatHandler() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
