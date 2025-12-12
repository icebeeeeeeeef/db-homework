package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"webook/internal/repository/cache"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicateEmail = errors.New("Email already exists")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	UpdateUserProfile(ctx context.Context, u User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWechat(ctx context.Context, openID string) (User, error)
}

type GORMUserDAO struct {
	db    *gorm.DB
	cache cache.UserCache
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

// 这里的User结构体是数据库中的user表的结构体
type User struct {
	Id            int64          `gorm:"primaryKey,autoIncrement"`
	Email         sql.NullString `gorm:"unique"`
	Phone         sql.NullString `gorm:"unique"`
	Password      string
	Birthday      string         `gorm:"column:birthday"`
	AboutMe       string         `gorm:"column:about_me"`
	Nickname      string         `gorm:"column:nickname"`
	CTime         int64          `gorm:"column:c_time"`
	UTime         int64          `gorm:"column:u_time"`
	WechatUnionID sql.NullString `gorm:"unique"`
	WechatOpenID  sql.NullString `gorm:"unique"`
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.UTime = now
	u.CTime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok { //用来输出unique字段冲突导致的错误
		const uniqueIndexErrNo uint16 = 1062
		if mysqlErr.Number == uniqueIndexErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}
func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}
func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error
	return u, err
}
func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error
	return u, err
}
func (dao *GORMUserDAO) UpdateUserProfile(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()

	// 只更新指定的字段，包括更新时间
	updates := map[string]interface{}{
		"u_time": now, // 使用数据库字段名
	}

	if u.Nickname != "" {
		updates["nickname"] = u.Nickname
	}
	if u.Birthday != "" {
		updates["birthday"] = u.Birthday
	}
	if u.AboutMe != "" {
		updates["about_me"] = u.AboutMe
	}
	if u.Phone.Valid {
		updates["phone"] = u.Phone.String
	}
	if u.Email.Valid {
		updates["email"] = u.Email.String
	}
	fmt.Printf("DAO层更新用户信息: ID=%v, updates=%+v\n", u.Id, updates)

	err := dao.db.WithContext(ctx).Model(&User{}).Where("id=?", u.Id).Updates(updates).Error
	if err != nil {
		fmt.Printf("DAO层更新用户信息失败: %v\n", err)
	}
	return err
}

func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openID string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id=?", openID).First(&u).Error
	return u, err
}
