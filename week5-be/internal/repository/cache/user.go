package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"geek-hw-week5/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, domainUser domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	data, err := cache.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal([]byte(data), &user)
	return user, err
}

func (cache *RedisUserCache) Set(ctx context.Context, domainUser domain.User) error {
	key := cache.key(domainUser.Id)
	data, err := json.Marshal(domainUser)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, key, data, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
