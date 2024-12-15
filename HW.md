# Homework Week 4

- `internal/repository/cache/code.go`

## Interface

```go
type CodeCache interface {
	Set(ctx context.Context, businessType, phone, code string) error
	Verify(ctx context.Context, businessType, phone, code string) (bool, error)
}
```

## Local Cache

```go
type LocalCodeCacheStore struct {
	Val string
	Cnt int
	Exp int64
}

// 使用 map 做为本地缓存
// 使用 RWLock 保证并发读写安全，读锁非互斥，写锁互斥
// 添加建议 cleanup 算法：在 Set 时，如果 map 中元素个数为 128 的倍数，则进行遍历删除无用键值对
type LocalCodeCache struct {
	lock  sync.RWMutex
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
	cache.lock.RLock()
	defer cache.lock.RUnlock()
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
```

## Initialization

```go
func NewCodeCache(cmd redis.Cmdable) CodeCache {
	//return &RedisCodeCache{cmd: cmd}
	return &LocalCodeCache{
		store: make(map[string]*LocalCodeCacheStore),
	}
}
```

