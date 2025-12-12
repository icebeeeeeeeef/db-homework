//go:build test

package startup

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webook/internal/repository/dao"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserHandlerTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (s *UserHandlerTestSuite) SetupSuite() {
	s.server = gin.Default()
	var err error
	viper.Set("db.mysql.dsn", "root:123123@tcp(localhost:3306)/webook_test?parseTime=true&loc=Local&charset=utf8mb4")
	s.db, err = InitDBNoLogger()
	require.NoError(s.T(), err)
	userHdl := InitUserHandler()
	userHdl.RegisterRoutes(s.server)
}

func (s *UserHandlerTestSuite) TearDownTest() {
	// 清理测试数据
	s.db.Exec("TRUNCATE TABLE users")
}

func (s *UserHandlerTestSuite) TestUserHandler_Signup() {
	t := s.T()
	testCases := []struct {
		name string
		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)
		// 构造请求
		req UserSignupRequest
		// 预期响应
		wantCode   int
		wantResult web.Result[any]
	}{
		{
			name: "正常注册",
			before: func(t *testing.T) {
				// 不需要准备数据
			},
			after: func(t *testing.T) {
				// 验证用户是否创建
				var user dao.User
				err := s.db.Where("email = ?", "test@example.com").First(&user).Error
				assert.NoError(t, err)
				assert.Equal(t, "test@example.com", user.Email.String)
				assert.Equal(t, "测试用户", user.Nickname)
				assert.True(t, user.CTime > 0)
				assert.True(t, user.UTime > 0)
			},
			req: UserSignupRequest{
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				Nickname:        "测试用户",
			},
			wantCode: 200,
			wantResult: web.Result[any]{
				Code: 0,
				Msg:  "注册成功",
			},
		},
		{
			name: "邮箱已存在",
			before: func(t *testing.T) {
				// 创建已存在的用户
				s.db.Create(&dao.User{
					Email:    sql.NullString{String: "existing@example.com", Valid: true},
					Password: "hashed_password",
					Nickname: "已存在用户",
					CTime:    time.Now().Unix(),
					UTime:    time.Now().Unix(),
				})
			},
			after: func(t *testing.T) {
				// 验证用户数量没有增加
				var count int64
				s.db.Model(&dao.User{}).Count(&count)
				assert.Equal(t, int64(1), count)
			},
			req: UserSignupRequest{
				Email:           "existing@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				Nickname:        "新用户",
			},
			wantCode: 200,
			wantResult: web.Result[any]{
				Code: 4,
				Msg:  "邮箱已存在",
			},
		},
		{
			name: "密码不匹配",
			before: func(t *testing.T) {
				// 不需要准备数据
			},
			after: func(t *testing.T) {
				// 验证用户没有创建
				var count int64
				s.db.Model(&dao.User{}).Count(&count)
				assert.Equal(t, int64(0), count)
			},
			req: UserSignupRequest{
				Email:           "test2@example.com",
				Password:        "password123",
				ConfirmPassword: "different_password",
				Nickname:        "测试用户2",
			},
			wantCode: 200,
			wantResult: web.Result[any]{
				Code: 4,
				Msg:  "两次输入的密码不一致",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			var result web.Result[any]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}

func (s *UserHandlerTestSuite) TestUserHandler_Login() {
	t := s.T()
	testCases := []struct {
		name string
		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)
		// 构造请求
		req UserLoginRequest
		// 预期响应
		wantCode   int
		wantResult web.Result[any]
	}{
		{
			name: "正常登录",
			before: func(t *testing.T) {
				// 创建测试用户
				s.db.Create(&dao.User{
					Email:    sql.NullString{String: "login@example.com", Valid: true},
					Password: "$2a$10$example_hashed_password", // 这里应该是实际加密后的密码
					Nickname: "登录用户",
					CTime:    time.Now().Unix(),
					UTime:    time.Now().Unix(),
				})
			},
			after: func(t *testing.T) {
				// 验证登录成功，这里可以验证JWT token等
			},
			req: UserLoginRequest{
				Email:    "login@example.com",
				Password: "password123",
			},
			wantCode: 200,
			wantResult: web.Result[any]{
				Code: 0,
				Msg:  "登录成功",
			},
		},
		{
			name: "用户不存在",
			before: func(t *testing.T) {
				// 不需要准备数据
			},
			after: func(t *testing.T) {
				// 验证没有创建新用户
				var count int64
				s.db.Model(&dao.User{}).Count(&count)
				assert.Equal(t, int64(0), count)
			},
			req: UserLoginRequest{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			wantCode: 200,
			wantResult: web.Result[any]{
				Code: 4,
				Msg:  "用户不存在",
			},
		},
		{
			name: "密码错误",
			before: func(t *testing.T) {
				// 创建测试用户
				s.db.Create(&dao.User{
					Email:    sql.NullString{String: "wrongpass@example.com", Valid: true},
					Password: "$2a$10$example_hashed_password",
					Nickname: "密码错误用户",
					CTime:    time.Now().Unix(),
					UTime:    time.Now().Unix(),
				})
			},
			after: func(t *testing.T) {
				// 验证用户存在但没有登录成功
				var user dao.User
				err := s.db.Where("email = ?", "wrongpass@example.com").First(&user).Error
				assert.NoError(t, err)
			},
			req: UserLoginRequest{
				Email:    "wrongpass@example.com",
				Password: "wrong_password",
			},
			wantCode: 200,
			wantResult: web.Result[any]{
				Code: 4,
				Msg:  "密码错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			var result web.Result[any]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}

func (s *UserHandlerTestSuite) TestUserHandler_Profile() {
	t := s.T()
	testCases := []struct {
		name string
		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)
		// 构造请求
		req UserProfileRequest
		// 预期响应
		wantCode   int
		wantResult web.Result[any]
	}{
		{
			name: "正常更新用户信息",
			before: func(t *testing.T) {
				// 创建测试用户
				s.db.Create(&dao.User{
					Id:       123,
					Email:    sql.NullString{String: "profile@example.com", Valid: true},
					Password: "hashed_password",
					Nickname: "原昵称",
					AboutMe:  "原简介",
					CTime:    time.Now().Unix(),
					UTime:    time.Now().Unix(),
				})
			},
			after: func(t *testing.T) {
				// 验证用户信息已更新
				var user dao.User
				err := s.db.Where("id = ?", 123).First(&user).Error
				assert.NoError(t, err)
				assert.Equal(t, "新昵称", user.Nickname)
				assert.Equal(t, "新简介", user.AboutMe)
				assert.True(t, user.UTime > 0) // 更新时间应该变化
			},
			req: UserProfileRequest{
				Nickname: "新昵称",
				AboutMe:  "新简介",
			},
			wantCode: 200,
			wantResult: web.Result[any]{
				Code: 0,
				Msg:  "更新成功",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			// 设置用户上下文
			s.server.Use(func(c *gin.Context) {
				c.Set("claims", &ijwt.UserClaims{
					UserId: 123,
				})
			})
			data, err := json.Marshal(tc.req)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/users/profile", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			var result web.Result[any]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}

// 请求结构体
type UserSignupRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Nickname        string `json:"nickname"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfileRequest struct {
	Nickname string `json:"nickname"`
	AboutMe  string `json:"aboutMe"`
}
