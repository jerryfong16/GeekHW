package repository

import (
	"context"
	"database/sql"
	"geek-hw-week5/internal/domain"
	"geek-hw-week5/internal/repository/cache"
	"geek-hw-week5/internal/repository/dao"
	"log"
	"time"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateEmail
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateNonsensitiveFields(ctx context.Context, user domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func (repository *CachedUserRepository) Create(ctx context.Context, user domain.User) error {
	return repository.dao.Insert(ctx, repository.toEntity(user))
}

func (repository *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repository.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repository.toDomain(user), err
}

func (repository *CachedUserRepository) UpdateNonsensitiveFields(ctx context.Context, user domain.User) error {
	return repository.dao.UpdateById(ctx, repository.toEntity(user))
}

func (repository *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := repository.dao.FindByEmail(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repository.toDomain(user), err
}

func (repository *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	domainUser, err := repository.cache.Get(ctx, id)
	if err == nil {
		return domainUser, nil
	}

	user, err := repository.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, ErrUserNotFound
	}
	domainUser = repository.toDomain(user)

	err = repository.cache.Set(ctx, domainUser)
	if err != nil {
		log.Println(err)
	}
	return domainUser, nil
}

func (repository *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

func (repository *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}
