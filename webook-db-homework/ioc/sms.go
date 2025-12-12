package ioc

import (
	"webook/internal/service/sms"
	"webook/internal/service/sms/memory"
)

func InitSMS() sms.Service {
	/*
		// 创建腾讯云短信客户端
		credential := common.NewCredential(
			"your_secret_id",  // 替换为您的SecretId
			"your_secret_key", // 替换为您的SecretKey
		)

		client, err := smsClient.NewClient(credential, "ap-nanjing", profile.NewClientProfile())
		if err != nil {
			panic(err)
		}
	*/
	smsService := memory.NewService()
	// 创建短信服务
	return smsService
}
