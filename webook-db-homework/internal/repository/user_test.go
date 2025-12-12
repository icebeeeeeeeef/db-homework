package repository

import (
	"context"
	"errors"
	"testing"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"

	cachemocks "webook/internal/repository/cache/mocks"
	daomocks "webook/internal/repository/dao/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	testCase := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		//input
		id int64
		//output
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				dao := daomocks.NewMockUserDAO(ctrl)
				cache := cachemocks.NewMockUserCache(ctrl)
				cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, nil)
				return dao, cache
			},
			id:       1,
			wantUser: domain.User{},
			wantErr:  nil,
		},
		{
			name: "数据库查找正常找到数据",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daomock := daomocks.NewMockUserDAO(ctrl)
				cachemock := cachemocks.NewMockUserCache(ctrl)
				cachemock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, cache.ErrUserNotFound)
				daomock.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(dao.User{}, nil)
				cachemock.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)
				return daomock, cachemock
			},
			id:       1,
			wantUser: domain.User{},
			wantErr:  nil,
		},
		{
			name: "数据库中未找到数据，缓存空值",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daomock := daomocks.NewMockUserDAO(ctrl)
				cachemock := cachemocks.NewMockUserCache(ctrl)
				cachemock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, cache.ErrUserNotFound)
				daomock.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(dao.User{}, errors.New("数据库中未找到数据"))
				cachemock.EXPECT().SetEmpty(gomock.Any(), gomock.Any()).Return(nil)
				return daomock, cachemock
			},
			id:       1,
			wantUser: domain.User{},
			wantErr:  errors.New("数据库中未找到数据"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao, cache := tc.mock(ctrl)
			repo := NewUserRepository(dao, cache)
			user, err := repo.FindById(context.Background(), tc.id)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}
