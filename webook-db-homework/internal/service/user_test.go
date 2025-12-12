package service

import (
	"context"
	"errors"
	"testing"

	"webook/internal/domain"
	"webook/internal/repository"
	repomocks "webook/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func Test_userLogin(t *testing.T) {
	testCase := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		//输入
		//ctx      context.Context
		email    string
		password string

		//输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(domain.User{
					Email:    "123456@qq.com",
					Password: "$2a$10$HxsHV6T919cua7qSDDRkD.ecElRjt013pu5YWTLq.7BriVdM8c37y!",
				}, nil)
				return repo
			},
			email:    "123456@qq.com",
			password: "Test123!",
			wantUser: domain.User{
				Email:    "123456@qq.com",
				Password: "$2a$10$HxsHV6T919cua7qSDDRkD.ecElRjt013pu5YWTLq.7BriVdM8c37y!",
			},
			wantErr: nil,
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(domain.User{
					Email:    "123456@qq.com",
					Password: "$2a$10$HxsHV6T919cua7qSDDRkD.ecElRjt013pu5YWTLq.7BriVdM8c37y!",
				}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123456@qq.com",
			password: "Test123!",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "其他错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(domain.User{
					Email:    "123456@qq.com",
					Password: "$2a$10$HxsHV6T919cua7qSDDRkD.ecElRjt013pu5YWTLq.7BriVdM8c37y!",
				}, errors.New("其他错误"))
				return repo
			},
			email:    "123456@qq.com",
			password: "Test123!",
			wantUser: domain.User{},
			wantErr:  errors.New("其他错误"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(domain.User{
					Email:    "123456@qq.com",
					Password: "$2a$1cua7qSDDRkD.ecElRjt013pu5YWTLq.7BriVdM8c37y!",
				}, nil)
				return repo
			},
			email:    "123456@qq.com",
			password: "Test123!",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			user, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func Test_Encrypt(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("Test123!"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(res))
}
