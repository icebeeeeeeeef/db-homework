package dao

import (
	"context"
	"errors"
	"testing"

	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(t *testing.T) *sql.DB
		u       User
		ctx     context.Context
		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO `users`.*").WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			u: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			ctx:     context.Background(),
			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users`.*").WillReturnError(&mysql.MySQLError{Number: 1062})
				require.NoError(t, err)
				return mockDB
			},
			u: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			ctx:     context.Background(),
			wantErr: ErrUserDuplicateEmail,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users`.*").WillReturnError(errors.New("数据库错误"))
				require.NoError(t, err)
				return mockDB
			},
			u: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			ctx:     context.Background(),
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true, //跳过版本检查
			}), &gorm.Config{
				DisableAutomaticPing:   true, //禁用自动ping，因为gorm执行sql语句时会自动ping
				SkipDefaultTransaction: true, //禁用事务，因为gorm执行sql语句时会自动开启事务
			})
			d := NewUserDAO(db)
			err = d.Insert(tc.ctx, tc.u)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}
