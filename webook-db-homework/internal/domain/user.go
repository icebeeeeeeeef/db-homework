package domain

import "time"

type User struct {
	Id         int64
	Email      string
	Phone      string // 新增手机号字段
	Password   string
	Birthday   string
	Nickname   string
	Ctime      time.Time
	AboutMe    string
	WechatUser WechatUser //这里为什么不组合，因为可能还有其他比如DingDingInfo 可能会有同名的字段
}
