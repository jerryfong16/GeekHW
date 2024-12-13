package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string

	ErrCodeSendTooMany   = errors.New("发送太频繁")
	ErrCodeVerifyTooMany = errors.New("验证太频繁")
)

type CodeCache interface {
	Set(ctx context.Context, businessType, phone, code string) error
	Verify(ctx context.Context, businessType, phone, code string) (bool, error)
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

func (cache *RedisCodeCache) Set(ctx context.Context, businessType, phone, code string) error {
	res, err := cache.cmd.Eval(ctx, luaSetCode, []string{cache.key(businessType, phone)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return err
	}
	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (cache *RedisCodeCache) Verify(ctx context.Context, businessType, phone, code string) (bool, error) {
	res, err := cache.cmd.Eval(ctx, luaVerifyCode, []string{cache.key(businessType, phone)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return false, err
	}
	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}

func (cache *RedisCodeCache) key(businessType, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", businessType, phone)
}

func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{cmd: cmd}
}
