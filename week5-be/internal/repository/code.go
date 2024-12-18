package repository

import (
	"context"
	"geek-hw-week5/internal/repository/cache"
)

var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany

type CodeRepository interface {
	Set(ctx context.Context, businessType, phone, code string) error
	Verify(ctx context.Context, businessType, phone, code string) (bool, error)
}

type CachedCodeRepository struct {
	cache cache.CodeCache
}

func (repository *CachedCodeRepository) Set(ctx context.Context, businessType, phone, code string) error {
	return repository.cache.Set(ctx, businessType, phone, code)
}

func (repository *CachedCodeRepository) Verify(ctx context.Context, businessType, phone, code string) (bool, error) {
	return repository.cache.Verify(ctx, businessType, phone, code)
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{cache: cache}
}
