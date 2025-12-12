package startup

import (
	"strings"
	"time"
	"webook/internal/web"

	ijwt "webook/internal/web/jwt"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middlewares/ratelimit"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandler, oauth2WechatHdl *web.OAuth2WechatHandler, articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable, jwtHandler ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "http://127.0.0.1:8080", "http://localhost:5500", "http://127.0.0.1:5500"},
			AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization", "Origin", "Accept"},
			ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
			AllowCredentials: true, // 是否允许带cookie等凭证
			AllowOriginFunc: func(origin string) bool {
				if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
					return true
				}
				return false
			},

			MaxAge: 12 * time.Hour,
		}),
		middleware.NewLoginJwtMiddlewareBuilder(jwtHandler).IgnorePaths("/users/login").IgnorePaths("/users/signup").IgnorePaths("/users/login_sms/code/send").IgnorePaths("/users/login_sms").IgnorePaths("oauth2/wechat/authurl").IgnorePaths("oauth2/wechat/callback").IgnorePaths("/users/login_jwt").IgnorePaths("/users/refresh_token").Build(),
	}
}
