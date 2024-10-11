package repository

import (
	"context"
	"fzy.com/geek-hw-week2/domain"
	"fzy.com/geek-hw-week2/repository/dao"
	"time"
)

var (
	ErrDuplicateEmail  = dao.ErrDuplicateEmail
	ErrAccountNotFound = dao.ErrRecordNotFound
)

type AccountRepository struct {
	dao *dao.AccountDAO
}

func NewAccountRepository(dao *dao.AccountDAO) *AccountRepository {
	return &AccountRepository{
		dao: dao,
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
	account, err := accountRepository.dao.FindById(ctx, id)
	if err != nil {
		return domain.Account{}, err
	}
	return accountRepository.toDomainAccount(account), err
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
