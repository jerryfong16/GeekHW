package service

import (
	"context"
	"errors"
	"fzy.com/geek-hw-week2/domain"
	"fzy.com/geek-hw-week2/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail         = repository.ErrDuplicateEmail
	ErrInvalidEmailOrPassword = errors.New("account not exists")
)

type AccountService struct {
	accountRepository *repository.AccountRepository
}

func NewAccountService(accountRepository *repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
	}
}

func (accountService *AccountService) Signup(ctx context.Context, account domain.Account) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Password = string(encryptedPassword)
	return accountService.accountRepository.Create(ctx, account)
}

func (accountService *AccountService) Login(ctx context.Context, email string, password string) (domain.Account, error) {
	account, err := accountService.accountRepository.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrAccountNotFound) {
		return domain.Account{}, ErrInvalidEmailOrPassword
	}
	if err != nil {
		return domain.Account{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return domain.Account{}, ErrInvalidEmailOrPassword
	}
	return account, nil
}

func (accountService *AccountService) EditProfile(ctx context.Context, account domain.Account, name string, birth string, about string) error {
	account.Name = name
	account.Birth = birth
	account.About = about
	return accountService.accountRepository.Update(ctx, account)
}

func (accountService *AccountService) GetProfileById(ctx context.Context, id int64) (domain.Account, error) {
	account, err := accountService.accountRepository.FindById(ctx, id)
	if err != nil {
		return domain.Account{}, err
	}
	return account, nil
}
