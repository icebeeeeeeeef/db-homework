package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("邮箱/密码错误")

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	UpdateUserProfile(ctx context.Context, u domain.User) error
	GetUserById(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatUser domain.WechatUser) (domain.User, error)
}

type UserService_ struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserService_{
		repo: repo,
	}
}

func (svc *UserService_) Signup(ctx context.Context, u domain.User) error {
	//service层中要考虑把密码加密以及存储起来的问题，而这些则属于repo层的业务
	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encrypted)
	return svc.repo.Create(ctx, u)
}
func (svc *UserService_) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
func (svc *UserService_) UpdateUserProfile(ctx context.Context, u domain.User) error {
	return svc.repo.UpdateUserProfile(ctx, u)
}
func (svc *UserService_) GetUserById(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *UserService_) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, nil
	}
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return domain.User{}, err
	}
	//这里的u只有一个phone，因此我们要再查一次
	u, err = svc.repo.FindByPhone(ctx, phone)
	//但是这里会遇到主从延迟的问题
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (svc *UserService_) FindOrCreateByWechat(ctx context.Context, wechatUser domain.WechatUser) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechatUser.OpenID)
	if err != repository.ErrUserNotFound {
		return u, nil
	}
	u = domain.User{
		WechatUser: domain.WechatUser{
			UnionID: wechatUser.UnionID,
			OpenID:  wechatUser.OpenID,
		},
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return domain.User{}, err
	}
	//这里的u只有一个phone，因此我们要再查一次
	u, err = svc.repo.FindByWechat(ctx, wechatUser.OpenID)
	//但是这里会遇到主从延迟的问题
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
