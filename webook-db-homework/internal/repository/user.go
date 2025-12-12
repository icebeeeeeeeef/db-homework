package repository

import (
	"context"
	"database/sql"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	UpdateUserProfile(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}
func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, DomainToEntity(u))
}
func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	return EntityToDomain(u), err
}
func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 1. 先查缓存
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 如果缓存的是空值占位，视为未命中
		if u.Id == 0 && u.Email == "" && u.Nickname == "" {
			// fallthrough to DB
		} else {
			return u, nil
		}
	}

	// 2. 缓存未找到，查数据库
	if err == cache.ErrUserNotFound {
		ur, err := r.dao.FindById(ctx, id)
		if err != nil {
			// 数据库也没有，缓存空值防止缓存穿透

			go func() {
				_ = r.cache.SetEmpty(ctx, id)
			}()

			_ = r.cache.SetEmpty(ctx, id)
			return domain.User{}, err
		}

		// 3. 数据库查到数据，转换为domain对象
		u = EntityToDomain(ur)

		// 4. 异步回填缓存，不阻塞主流程

		//_ = r.cache.Set(ctx, u)

		go func() {
			_ = r.cache.Set(ctx, u)
		}()

		return u, nil
	}

	// 5. 其他错误（如Redis连接失败等）
	return domain.User{}, err
}
func (r *CachedUserRepository) UpdateUserProfile(ctx context.Context, u domain.User) error {
	err := r.dao.UpdateUserProfile(ctx, DomainToEntity(u))
	if err != nil {
		return err
	}
	if r.cache != nil {
		go func() {
			_ = r.cache.Delete(context.Background(), u.Id)
		}()
	}
	return nil
}

func DomainToEntity(u domain.User) dao.User {
	return dao.User{
		Id:            u.Id,
		Email:         sql.NullString{String: u.Email, Valid: u.Email != ""},
		Phone:         sql.NullString{String: u.Phone, Valid: u.Phone != ""},
		Password:      u.Password,
		Nickname:      u.Nickname,
		Birthday:      u.Birthday,
		AboutMe:       u.AboutMe,
		WechatUnionID: sql.NullString{String: u.WechatUser.UnionID, Valid: u.WechatUser.UnionID != ""},
		WechatOpenID:  sql.NullString{String: u.WechatUser.OpenID, Valid: u.WechatUser.OpenID != ""},
	}
}

func EntityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		WechatUser: domain.WechatUser{
			UnionID: u.WechatUnionID.String,
			OpenID:  u.WechatOpenID.String,
		},
	}
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	return EntityToDomain(u), err
}

func (r *CachedUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := r.dao.FindByWechat(ctx, openID)
	return EntityToDomain(u), err
}
