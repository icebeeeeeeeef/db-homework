package middleware

import (
	"fmt"
	"net/http"
	"strings"
	ijwt "webook/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJwtMiddlewareBuilder struct {
	ijwt.Handler
	Paths []string
}

func NewLoginJwtMiddlewareBuilder(jwtHandler ijwt.Handler) *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{
		Handler: jwtHandler,
	}
}

func (b *LoginJwtMiddlewareBuilder) IgnorePaths(path string) *LoginJwtMiddlewareBuilder {
	b.Paths = append(b.Paths, path)
	return b
}

func (b *LoginJwtMiddlewareBuilder) Build() gin.HandlerFunc {
	//注意部分的请求不需要校验session
	return func(c *gin.Context) {
		// 检查白名单路径
		for _, path := range b.Paths {
			if c.Request.URL.Path == path || strings.HasPrefix(c.Request.URL.Path, path+"/") {
				// 对公开文章：允许未登录访问，但如果带了 token 则解析出 userId 以返回点赞/收藏状态
				if path == "/articles/pub" {
					tokenStr := b.ExtractToken(c)
					if tokenStr == "" {
						c.Next()
						return
					}
					claims := &ijwt.UserClaims{}
					token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("无效的签名算法: %v", token.Header["alg"])
						}
						return []byte("1234567890"), nil
					})
					if err != nil || !token.Valid || claims.UserAgent != c.Request.UserAgent() || b.CheckSSid(c, claims.SSid) != nil {
						// token 无效则按未登录处理
						c.Next()
						return
					}
					c.Set("userClaims", claims)
					c.Set("userId", claims.UserId)
					c.Next()
					return
				}
				c.Next()
				return
			}
		}
		//用jwt来校验
		tokenStr := b.ExtractToken(c)
		claims := &ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效的签名算法: %v", token.Header["alg"])
			}
			return []byte("1234567890"), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "JWT验证失败: " + err.Error(),
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "JWT无效",
			})
			return
		}
		if claims.UserAgent != c.Request.UserAgent() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户代理不匹配",
			})
			return
		}
		//这里可以加上如果redis崩掉就跳过下一个校验，一种降级策略
		err = b.CheckSSid(c, claims.SSid)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} //这里就是在检查是否登出

		c.Set("userClaims", claims)
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
