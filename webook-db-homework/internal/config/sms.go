package config

// TencentSMSConfig 腾讯云短信服务配置
type TencentSMSConfig struct {
	SecretId  string            `json:"secret_id"`  // 腾讯云 SecretId
	SecretKey string            `json:"secret_key"` // 腾讯云 SecretKey
	Region    string            `json:"region"`     // 地域，如 ap-guangzhou
	AppId     string            `json:"app_id"`     // 短信应用ID
	SignName  string            `json:"sign_name"`  // 短信签名
	Templates map[string]string `json:"templates"`  // 模板ID映射
}

// DefaultTencentSMSConfig 默认配置
func DefaultTencentSMSConfig() *TencentSMSConfig {
	return &TencentSMSConfig{
		SecretId:  "", // 需要从环境变量或配置文件读取
		SecretKey: "", // 需要从环境变量或配置文件读取
		Region:    "ap-guangzhou",
		AppId:     "",     // 需要从环境变量或配置文件读取
		SignName:  "您的签名", // 需要从环境变量或配置文件读取
		Templates: map[string]string{
			"login":    "", // 登录模板ID
			"register": "", // 注册模板ID
			"reset":    "", // 重置密码模板ID
		},
	}
}

// GetTemplateID 获取指定类型的模板ID
func (c *TencentSMSConfig) GetTemplateID(smsType string) string {
	if templateID, exists := c.Templates[smsType]; exists {
		return templateID
	}
	return ""
}
