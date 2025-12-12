package ginx

import (
	"net/http"
	"webook/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// 这个函数就是用来简化一些无需校验身份但是需要绑定body的操作的，比如注册操作
func WrapBody[T any](l logger.LoggerV1, f func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusOK, Result{Code: 400, Msg: err.Error()})
			return
		}
		result, err := f(ctx, req)
		if err != nil {

			//由于考虑到在外部进行错误处理和日志打印，因此无法实现具体错误信息和日志打印
			//这里通过打印请求的path和method来实现日志打印
			//或者其实我们可以通过error的包装来实现打印具体错误信息，也就是在回调函数中包装错误，在wrap中打印错误日志
			l.Error("系统错误", logger.Error(err),
				logger.String("path", ctx.Request.URL.Path),
				logger.String("method", ctx.Request.Method),
			)

			ctx.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		}
		ctx.JSON(http.StatusOK, result)
	}
}

func WrapToken[Claims jwt.Claims](l logger.LoggerV1, f func(ctx *gin.Context, C Claims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get("userId")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, Result{Code: 401, Msg: "用户未登录"})
			return
		}
		claims, ok := val.(Claims)
		if !ok {
			ctx.JSON(http.StatusOK, Result{Code: 401, Msg: "用户未登录"})
			return
		}
		result, err := f(ctx, claims)
		if err != nil {
			l.Error("系统错误", logger.Error(err),
				logger.String("path", ctx.Request.URL.Path),
				logger.String("method", ctx.Request.Method),
			)
		}
		ctx.JSON(http.StatusOK, result)
	}
}

func WrapTokenBody[Claims jwt.Claims, Body any](l logger.LoggerV1, f func(ctx *gin.Context, c Claims, req Body) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Body
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusOK, Result{Code: 400, Msg: err.Error()})
			return
		}
		val, ok := ctx.Get("userId")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, Result{Code: 401, Msg: "用户未登录"})
			return
		}
		claims, ok := val.(Claims)
		if !ok {
			ctx.JSON(http.StatusOK, Result{Code: 401, Msg: "用户未登录"})
			return
		}
		result, err := f(ctx, claims, req)
		if err != nil {
			l.Error("系统错误", logger.Error(err),
				logger.String("path", ctx.Request.URL.Path),
				logger.String("method", ctx.Request.Method),
			)
		}
		ctx.JSON(http.StatusOK, result)
	}
}
