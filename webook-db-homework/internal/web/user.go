package web

import (
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"

	"webook/internal/domain"
	"webook/internal/service"
	ijwt "webook/internal/web/jwt"
)

var ErrUserDuplicateEmail = service.ErrUserDuplicateEmail

// UserHandler 精简版：邮箱注册/登录 + 个人资料编辑与查询
type UserHandler struct {
	svc         service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	ijwt.Handler
}

func NewUserHandler(svc service.UserService, jwtHandler ijwt.Handler) *UserHandler {
	const (
		emailRegexPattern    = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
		passWordRegexPattern = "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[!@#$%^&*()_+])[A-Za-z\\d!@#$%^&*()_+]{8,}$"
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passWordExp := regexp.MustCompile(passWordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passWordExp,
		Handler:     jwtHandler,
	}
}

func (u *UserHandler) RegisterRoutes(r *gin.Engine) {
	ug := r.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.GET("/profile", u.Profile)
	ug.POST("/edit", u.Edit)
	ug.POST("/logout", u.Logout)
}

func (u *UserHandler) Signup(c *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "参数错误: " + err.Error()})
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil || !ok {
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "邮箱格式不正确"})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "两次密码不一致"})
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil || !ok {
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "密码需包含数字、字母、特殊字符，长度至少 8"})
		return
	}

	err = u.svc.Signup(c.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		c.JSON(http.StatusOK, Result[string]{Code: 0, Msg: "注册成功"})
	case service.ErrUserDuplicateEmail:
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "邮箱已注册"})
	default:
		c.JSON(http.StatusInternalServerError, Result[string]{Code: 500, Msg: "系统错误"})
	}
}

func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "参数错误: " + err.Error()})
		return
	}

	user, err := u.svc.Login(c.Request.Context(), req.Email, req.Password)
	switch err {
	case nil:
		if er := u.SetLoginToken(c, user.Id); er != nil {
			c.JSON(http.StatusInternalServerError, Result[string]{Code: 500, Msg: "系统错误"})
			return
		}
		c.JSON(http.StatusOK, Result[string]{Code: 0, Msg: "登录成功"})
	case service.ErrInvalidUserOrPassword:
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "邮箱或密码错误"})
	default:
		c.JSON(http.StatusInternalServerError, Result[string]{Code: 500, Msg: "系统错误"})
	}
}

func (u *UserHandler) Logout(c *gin.Context) {
	_ = u.ClearToken(c)
	c.JSON(http.StatusOK, Result[string]{Code: 0, Msg: "已退出登录"})
}

func (u *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req EditReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "参数错误: " + err.Error()})
		return
	}

	if req.Birthday != "" {
		if _, err := time.Parse("2006-01-02", req.Birthday); err != nil {
			c.JSON(http.StatusBadRequest, Result[string]{Code: 400, Msg: "生日格式不正确，使用 YYYY-MM-DD"})
			return
		}
	}

	uid := c.GetInt64("userId")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, Result[string]{Code: 401, Msg: "未登录"})
		return
	}

	err := u.svc.UpdateUserProfile(c.Request.Context(), domain.User{
		Id:       uid,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result[string]{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result[string]{Code: 0, Msg: "编辑成功"})
}

func (u *UserHandler) Profile(c *gin.Context) {
	uid := c.GetInt64("userId")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, Result[string]{Code: 401, Msg: "未登录"})
		return
	}
	user, err := u.svc.GetUserById(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result[string]{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result[domain.User]{Code: 0, Data: user})
}
