package cache

import (
	"context"
	"errors"
	"fmt"
	redismocks "geek-hw-week5/internal/repository/cache/mocks/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name  string
		build func(ctrl *gomock.Controller) redis.Cmdable

		ctx          context.Context
		businessType string
		phone        string
		code         string

		expectedErr error
	}{
		{
			name: "设置成功",
			build: func(ctrl *gomock.Controller) redis.Cmdable {
				redisCmd := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(0))
				cmd.SetErr(nil)
				redisCmd.EXPECT().
					Eval(
						gomock.Any(),
						luaSetCode,
						[]string{fmt.Sprintf("phone_code:%s:%s", "test", "13166668888")},
						[]any{"123456"},
					).Return(cmd)
				return redisCmd
			},
			ctx:          context.Background(),
			businessType: "test",
			phone:        "13166668888",
			code:         "123456",
			expectedErr:  nil,
		},
		{
			name: "redis异常",
			build: func(ctrl *gomock.Controller) redis.Cmdable {
				redisCmd := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(0))
				cmd.SetErr(errors.New("redis异常"))
				redisCmd.EXPECT().
					Eval(
						gomock.Any(),
						luaSetCode,
						[]string{fmt.Sprintf("phone_code:%s:%s", "test", "13166668888")},
						[]any{"123456"},
					).Return(cmd)
				return redisCmd
			},
			ctx:          context.Background(),
			businessType: "test",
			phone:        "13166668888",
			code:         "123456",
			expectedErr:  errors.New("redis异常"),
		},
		{
			name: "验证码发送过多",
			build: func(ctrl *gomock.Controller) redis.Cmdable {
				redisCmd := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-1))
				cmd.SetErr(nil)
				redisCmd.EXPECT().
					Eval(
						gomock.Any(),
						luaSetCode,
						[]string{fmt.Sprintf("phone_code:%s:%s", "test", "13166668888")},
						[]any{"123456"},
					).Return(cmd)
				return redisCmd
			},
			ctx:          context.Background(),
			businessType: "test",
			phone:        "13166668888",
			code:         "123456",
			expectedErr:  ErrCodeSendTooMany,
		},
		{
			name: "验证码存在缺少过期时间",
			build: func(ctrl *gomock.Controller) redis.Cmdable {
				redisCmd := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-2))
				cmd.SetErr(nil)
				redisCmd.EXPECT().
					Eval(
						gomock.Any(),
						luaSetCode,
						[]string{fmt.Sprintf("phone_code:%s:%s", "test", "13166668888")},
						[]any{"123456"},
					).Return(cmd)
				return redisCmd
			},
			ctx:          context.Background(),
			businessType: "test",
			phone:        "13166668888",
			code:         "123456",
			expectedErr:  errors.New("验证码存在，但是没有过期时间"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			redisCmd := testCase.build(ctrl)
			redisCache := NewCodeCache(redisCmd)

			err := redisCache.Set(testCase.ctx, testCase.businessType, testCase.phone, testCase.code)

			assert.Equal(t, err, testCase.expectedErr)
		})
	}
}
