package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"webook/internal/service"
	"webook/internal/service/oauth2/wechat"
	ijwt "webook/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
	ijwt.Handler
	UserService service.UserService
	stateKey    []byte
	cfg         WechatHandlerConfig
}

type WechatHandlerConfig struct {
	Secure bool
}

func NewOAuth2WechatHandler(svc wechat.Service, userService service.UserService, cfg WechatHandlerConfig, jwtHandler ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:         svc,
		UserService: userService,
		stateKey:    []byte("95osj3fUD7foisd7sdfk9sdf91eru"),
		cfg:         cfg,
		Handler:     jwtHandler,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 5,
			Msg:  err.Error(),
		})
		return
	}
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 5,
			Msg:  "构造扫码登录URL失败",
		})
		return
	}
	err = h.setStateCookie(ctx, state.String())
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 0,
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")

	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 5,
			Msg:  "验证码失败",
		})
		return
	}

	//这里要处理登录了，首先就是设置jwttoken
	uid, err := h.UserService.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}
	err = h.SetLoginToken(ctx, uid.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result[string]{
		Code: 0,
		Msg:  "啦啦啦啦啦",
	})
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("无法获取jwt-state %w", err)
	}
	var stateClaims StateClaims
	token, err := jwt.ParseWithClaims(ck, &stateClaims, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("无法解析jwt-state %w", err)
	}
	if state != stateClaims.State {
		return errors.New("状态码不匹配")
	}
	return nil
}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return fmt.Errorf("无法设置jwt-state %w", err)
	}
	ctx.SetCookie("jwt-state", tokenStr, 600, "/oauth2/wechat/callback", "", h.cfg.Secure, true)
	return nil
}

type StateClaims struct {
	State string `json:"state"`
	jwt.RegisteredClaims
}
