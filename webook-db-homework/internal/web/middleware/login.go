package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	Paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (b *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	b.Paths = append(b.Paths, path)
	return b
}

func (b *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	//注意部分的请求不需要校验session
	return func(c *gin.Context) {
		// 检查白名单路径
		for _, path := range b.Paths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// 获取session
		sess := sessions.Default(c)
		if sess == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "session初始化失败",
			})
			return
		}

		// 检查用户ID
		id := sess.Get("userId")
		if id == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户未登录",
			})
			return
		}

		// 检查session超时
		updateTime := sess.Get("update_time")
		now := time.Now().Unix()

		if updateTime == nil {
			// 首次访问，设置更新时间
			sess.Set("update_time", now)
			sess.Options(sessions.Options{
				MaxAge: 60, // 60秒过期
			})
			if err := sess.Save(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "session保存失败",
				})
				return
			}
			c.Next()
			return
		}

		// 验证更新时间类型
		updateTimeVal, ok := updateTime.(int64) //把interface{}转换为int64
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "session数据格式错误",
			})
			return
		}

		// 检查是否超时（60秒）
		if now-updateTimeVal > 60 {
			// session超时，清除session
			sess.Clear()
			if err := sess.Save(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "session清除失败",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "session已过期，请重新登录",
			})
			return
		}

		// 更新访问时间
		sess.Set("update_time", now)
		if err := sess.Save(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "session更新失败",
			})
			return
		}

		// 将用户ID存储到上下文中，供后续处理器使用
		c.Set("userId", id)
		c.Next() //执行后续的中间件或控制器
	}
}
