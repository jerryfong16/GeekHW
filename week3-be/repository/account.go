package repository

import (
	"context"
	"log"
	"time"

	"fzy.com/geek-hw-week3/domain"
	"fzy.com/geek-hw-week3/repository/cache"
	"fzy.com/geek-hw-week3/repository/dao"
)

var (
	ErrDuplicateEmail  = dao.ErrDuplicateEmail
	ErrAccountNotFound = dao.ErrRecordNotFound
)

type AccountRepository struct {
	dao   *dao.AccountDAO
	cache *cache.AccountCache
}

func NewAccountRepository(dao *dao.AccountDAO, cache *cache.AccountCache) *AccountRepository {
	return &AccountRepository{
		dao:   dao,
		cache: cache,
	}
}

func (accountRepository *AccountRepository) Create(ctx context.Context, account domain.Account) error {
	return accountRepository.dao.Insert(ctx, dao.Account{
		Email:       account.Email,
		Password:    account.Password,
		CreatedTime: account.CreatedTime.UnixMilli(),
		UpdatedTime: account.UpdatedTime.UnixMilli(),
	})
}

func (accountRepository *AccountRepository) Update(ctx context.Context, account domain.Account) error {
	return accountRepository.dao.Update(ctx, dao.Account{
		Id:          account.Id,
		Email:       account.Email,
		Password:    account.Password,
		Name:        account.Name,
		Birth:       account.Birth,
		About:       account.About,
		CreatedTime: account.CreatedTime.UnixMilli(),
		UpdatedTime: account.UpdatedTime.UnixMilli(),
	})
}

func (accountRepository *AccountRepository) FindById(ctx context.Context, id int64) (domain.Account, error) {
	cacheAccount, err := accountRepository.cache.Get(ctx, id)
	if err == nil {
		return cacheAccount, err
	}

	dbAccount, err := accountRepository.dao.FindById(ctx, id)
	if err != nil {
		return domain.Account{}, err
	}
	account := accountRepository.toDomainAccount(dbAccount)
	if err := accountRepository.cache.Set(ctx, account); err != nil {
		log.Println(err)
	}
	return account, err
}

func (accountRepository *AccountRepository) FindByEmail(ctx context.Context, email string) (domain.Account, error) {
	account, err := accountRepository.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.Account{}, err
	}
	return accountRepository.toDomainAccount(account), err
}

func (accountRepository *AccountRepository) toDomainAccount(account dao.Account) domain.Account {
	return domain.Account{
		Id:          account.Id,
		Email:       account.Email,
		Password:    account.Password,
		Name:        account.Name,
		Birth:       account.Birth,
		About:       account.About,
		CreatedTime: time.UnixMilli(account.CreatedTime),
		UpdatedTime: time.UnixMilli(account.UpdatedTime),
	}
}
