package service

import (
	"context"
	"errors"
	"fmt"
	"geek-hw-week5/internal/domain"
	"geek-hw-week5/internal/repository"
	repomocks "geek-hw-week5/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncryptedPasswordGeneration(t *testing.T) {
	plaintext := "123456"
	hash, _ := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	fmt.Println(string(hash))
}

func TestUserServiceImpl_Login(t *testing.T) {
	testCases := []struct {
		name string

		buildRepo func(ctrl *gomock.Controller) repository.UserRepository

		ctx      context.Context
		email    string
		password string

		expectedUser domain.User
		expectedErr  error
	}{
		{
			name: "登录成功",
			buildRepo: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(
						domain.User{
							Email:    "123@qq.com",
							Password: "$2a$10$uIkwx5EScHw97vbyVP5xheQpnH7U0ILNPLpz8zbHqWVUA/s8pVz6i",
							Phone:    "13166668888",
						},
						nil,
					)
				return repo
			},
			email:    "123@qq.com",
			password: "123456",
			expectedUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$uIkwx5EScHw97vbyVP5xheQpnH7U0ILNPLpz8zbHqWVUA/s8pVz6i",
				Phone:    "13166668888",
			},
			expectedErr: nil,
		},
		{
			name: "用户不存在",
			buildRepo: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(
						domain.User{},
						repository.ErrUserNotFound,
					)
				return repo
			},
			email:       "123@qq.com",
			password:    "123456",
			expectedErr: ErrInvalidUserOrPassword,
		},
		{
			name: "DB异常",
			buildRepo: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(
						domain.User{},
						errors.New("DB 异常"),
					)
				return repo
			},
			email:       "123@qq.com",
			password:    "123456",
			expectedErr: errors.New("DB 异常"),
		},
		{
			name: "密码错误",
			buildRepo: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(
						domain.User{
							Email:    "123@qq.com",
							Password: "$2a$10$uIkwx5EScHw97vbyVP5xheQpnH7U0ILNPLpz8zbHqWVUA/s8pVz6i",
							Phone:    "13166668888",
						},
						nil,
					)
				return repo
			},
			email:       "123@qq.com",
			password:    "12345",
			expectedErr: ErrInvalidUserOrPassword,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := testCase.buildRepo(ctrl)
			svc := NewUserService(repo)

			user, err := svc.Login(testCase.ctx, testCase.email, testCase.password)
			assert.Equal(t, user, testCase.expectedUser)
			assert.Equal(t, err, testCase.expectedErr)
		})
	}
}
