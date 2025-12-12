package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	AtKey = []byte("1234567890")
	RtKey = []byte("1234567890")
)

var _ Handler = (*RedisJWTHandler)(nil)

type RedisJWTHandler struct {
	atKey []byte
	rtKey []byte
	Cmd   redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		Cmd:   cmd,
		atKey: AtKey,
		rtKey: RtKey,
	}
}

func (h *RedisJWTHandler) SetJWTToken(c *gin.Context, userId int64, ssid string) error {
	claims := UserClaims{
		SSid:   ssid,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserAgent: c.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokentstr, err := token.SignedString(h.atKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		return err
	}
	c.Header("x-jwt-token", tokentstr)
	return nil
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	val, exists := ctx.Get("userClaims")
	if !exists {
		return nil
	}
	claims, ok := val.(*UserClaims)
	if !ok {
		return nil
	}
	return h.Cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.SSid), "1", time.Hour*24*7).Err()
}

func (h *RedisJWTHandler) CheckSSid(ctx *gin.Context, ssid string) error {
	val, err := h.Cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if val == 0 {
			return nil
		}
		return errors.New("ssid not invalid")
	default:
		return err
	}
}

func (h *RedisJWTHandler) ExtractToken(c *gin.Context) string {
	tokenHeader := c.GetHeader("Authorization")
	if tokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "缺少Authorization头",
		})
		return ""
	}

	parts := strings.Split(tokenHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func (h *RedisJWTHandler) SetLoginToken(c *gin.Context, userId int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(c, userId, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(c, userId, ssid)
	if err != nil {
		return err
	}
	return nil
}

func (h *RedisJWTHandler) SetRefreshToken(c *gin.Context, userId int64, ssid string) error {
	claims := RefreshClaims{
		SSid:   ssid,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokentstr, err := token.SignedString(h.rtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		return err
	}
	c.Header("x-refresh-token", tokentstr)
	return nil
}
