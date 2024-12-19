package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
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
	res, err := cache.cmd.Eval(
		ctx,
		luaSetCode,
		[]string{cache.key(businessType, phone)},
		code,
	).Int()
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

type LocalCodeCacheStore struct {
	Val string
	Cnt int
	Exp int64
}

// 使用 map 做为本地缓存
// 使用 sync.Mutex 保证并发读写安全
// 添加建议 cleanup 算法：在 Set 时，如果 map 中元素个数为 128 的倍数，则进行遍历删除无用键值对
type LocalCodeCache struct {
	lock  sync.Mutex
	store map[string]*LocalCodeCacheStore
}

func (cache *LocalCodeCache) Set(ctx context.Context, businessType, phone, code string) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	key := cache.key(businessType, phone)
	storeValue := LocalCodeCacheStore{
		Val: code,
		Cnt: 3,
		Exp: time.Now().UnixMilli() + 60*1000,
	}
	if v, ok := cache.store[key]; ok {
		if v.Exp <= time.Now().UnixMilli() {
			cache.add(key, storeValue)
			return nil
		}
		return errors.New("请求验证码过于频繁")
	} else {
		cache.add(key, storeValue)
		return nil
	}
}

func (cache *LocalCodeCache) Verify(ctx context.Context, businessType, phone, code string) (bool, error) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	key := cache.key(businessType, phone)
	if v, ok := cache.store[key]; ok {
		if v.Exp <= time.Now().UnixMilli() {
			delete(cache.store, key)
			return false, errors.New("查询不存在")
		}
		if v.Cnt == 0 {
			return false, errors.New("验证码错误次数过多")
		}
		if v.Val != code {
			v.Cnt--
			return false, errors.New("验证码错误")
		}
		delete(cache.store, key)
		return true, nil
	} else {
		return false, errors.New("查询不存在")
	}
}

func (cache *LocalCodeCache) key(businessType, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", businessType, phone)
}

func (cache *LocalCodeCache) add(k string, v LocalCodeCacheStore) {
	cache.store[k] = &v
	if len(cache.store)%128 == 0 {
		cache.cleanup()
	}
}

func (cache *LocalCodeCache) cleanup() {
	for k, v := range cache.store {
		if v.Exp <= time.Now().UnixMilli() {
			delete(cache.store, k)
		}
	}
}

func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{cmd: cmd}
	//return &LocalCodeCache{
	//	store: make(map[string]*LocalCodeCacheStore),
	//}
}