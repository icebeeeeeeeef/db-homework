package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 我们把jwttoken相关的操作都重构到这里，这样代码层次更加清晰
type Handler interface {
	ExtractToken(ctx *gin.Context) string
	SetLoginToken(ctx *gin.Context, userId int64) error
	SetJWTToken(ctx *gin.Context, userId int64, ssid string) error
	ClearToken(ctx *gin.Context) error
	CheckSSid(ctx *gin.Context, ssid string) error
}

type UserClaims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
	UserAgent string `json:"user_agent"`
	SSid      string `json:"ssid"`
}

type RefreshClaims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
	SSid string `json:"ssid"`
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
