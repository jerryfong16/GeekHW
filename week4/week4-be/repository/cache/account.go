package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"fzy.com/geek-hw-week4/domain"
	"github.com/redis/go-redis/v9"
)

type AccountCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewAccountCache(cmd redis.Cmdable) *AccountCache {
	return &AccountCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (accountCache *AccountCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (accountCache *AccountCache) Get(ctx context.Context, accountId int64) (domain.Account, error) {
	key := accountCache.key(accountId)
	data, err := accountCache.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.Account{}, err
	}
	var account domain.Account
	err = json.Unmarshal([]byte(data), &account)
	return account, err
}

func (accountCache *AccountCache) Set(ctx context.Context, account domain.Account) error {
	key := accountCache.key(account.Id)
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}
	return accountCache.cmd.Set(ctx, key, data, accountCache.expiration).Err()
}
