package repository

import (
	"context"
	"database/sql"
	"errors"
	"geek-hw-week5/internal/domain"
	"geek-hw-week5/internal/repository/cache"
	cachemocks "geek-hw-week5/internal/repository/cache/mocks"
	"geek-hw-week5/internal/repository/dao"
	daomocks "geek-hw-week5/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	currentTimestamp := time.Now().UnixMilli()
	testCases := []struct {
		name string

		build func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx context.Context
		id  int64

		expectedUser domain.User
		expectedErr  error
	}{
		{
			name: "缓存未命中查找成功",
			build: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDao := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				id := int64(123)

				userCache.EXPECT().
					Get(gomock.Any(), id).
					Return(domain.User{}, errors.New("缓存未命中"))

				userDao.EXPECT().
					FindById(gomock.Any(), id).
					Return(
						dao.User{
							Id: id,
							Email: sql.NullString{
								String: "123@qq.com",
								Valid:  true,
							},
							Password: "123456",
							Birthday: 100,
							AboutMe:  "自我介绍",
							Phone: sql.NullString{
								String: "15212345678",
								Valid:  true,
							},
							Ctime: currentTimestamp,
							Utime: 102,
						},
						nil,
					)

				userCache.EXPECT().
					Set(
						gomock.Any(),
						domain.User{
							Id:       123,
							Email:    "123@qq.com",
							Password: "123456",
							Birthday: time.UnixMilli(100),
							AboutMe:  "自我介绍",
							Phone:    "15212345678",
							Ctime:    time.UnixMilli(currentTimestamp),
						},
					).
					Return(nil)

				return userDao, userCache
			},
			ctx: context.Background(),
			id:  123,
			expectedUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: time.UnixMilli(100),
				AboutMe:  "自我介绍",
				Phone:    "15212345678",
				Ctime:    time.UnixMilli(currentTimestamp),
			},
			expectedErr: nil,
		},
		{
			name: "缓存命中",
			build: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDao := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				id := int64(123)

				userCache.EXPECT().
					Get(gomock.Any(), id).
					Return(
						domain.User{
							Id:       123,
							Email:    "123@qq.com",
							Password: "123456",
							Birthday: time.UnixMilli(100),
							AboutMe:  "自我介绍",
							Phone:    "15212345678",
							Ctime:    time.UnixMilli(currentTimestamp),
						},
						nil,
					)

				return userDao, userCache
			},
			ctx: context.Background(),
			id:  123,
			expectedUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: time.UnixMilli(100),
				AboutMe:  "自我介绍",
				Phone:    "15212345678",
				Ctime:    time.UnixMilli(currentTimestamp),
			},
			expectedErr: nil,
		},
		{
			name: "缓存未命中查找失败",
			build: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDao := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				id := int64(123)

				userCache.EXPECT().
					Get(gomock.Any(), id).
					Return(domain.User{}, errors.New("缓存未命中"))

				userDao.EXPECT().
					FindById(gomock.Any(), id).
					Return(dao.User{}, ErrUserNotFound)

				return userDao, userCache
			},
			ctx:          context.Background(),
			id:           123,
			expectedUser: domain.User{},
			expectedErr:  ErrUserNotFound,
		},
		{
			name: "缓存回写失败",
			build: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDao := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				id := int64(123)

				userCache.EXPECT().
					Get(gomock.Any(), id).
					Return(domain.User{}, errors.New("缓存未命中"))

				userDao.EXPECT().
					FindById(gomock.Any(), id).
					Return(
						dao.User{
							Id: id,
							Email: sql.NullString{
								String: "123@qq.com",
								Valid:  true,
							},
							Password: "123456",
							Birthday: 100,
							AboutMe:  "自我介绍",
							Phone: sql.NullString{
								String: "15212345678",
								Valid:  true,
							},
							Ctime: currentTimestamp,
							Utime: 102,
						},
						nil,
					)

				userCache.EXPECT().
					Set(
						gomock.Any(),
						domain.User{
							Id:       123,
							Email:    "123@qq.com",
							Password: "123456",
							Birthday: time.UnixMilli(100),
							AboutMe:  "自我介绍",
							Phone:    "15212345678",
							Ctime:    time.UnixMilli(currentTimestamp),
						},
					).
					Return(errors.New("缓存异常"))

				return userDao, userCache
			},
			ctx: context.Background(),
			id:  123,
			expectedUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: time.UnixMilli(100),
				AboutMe:  "自我介绍",
				Phone:    "15212345678",
				Ctime:    time.UnixMilli(currentTimestamp),
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userDao, userCache := testCase.build(ctrl)
			repo := NewUserRepository(userDao, userCache)

			user, err := repo.FindById(testCase.ctx, testCase.id)
			assert.Equal(t, user, testCase.expectedUser)
			assert.Equal(t, err, testCase.expectedErr)
		})
	}
}
