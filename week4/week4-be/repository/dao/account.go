package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	Email    string `gorm:"unique"`
	Password string
	Name     string
	Birth    string
	About    string

	CreatedTime int64
	UpdatedTime int64
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type AccountDAO struct {
	db *gorm.DB
}

func NewAccountDAO(db *gorm.DB) *AccountDAO {
	return &AccountDAO{
		db: db,
	}
}

func (dao *AccountDAO) Insert(ctx context.Context, account Account) error {
	t := time.Now().UnixMilli()
	account.CreatedTime = t
	account.UpdatedTime = t
	err := dao.db.WithContext(ctx).Create(&account).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == 1062 {
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *AccountDAO) Update(ctx context.Context, account Account) error {
	t := time.Now().UnixMilli()
	account.UpdatedTime = t
	err := dao.db.WithContext(ctx).Save(&account).Error
	return err
}

func (dao *AccountDAO) FindById(ctx context.Context, id int64) (Account, error) {
	var account Account
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&account).Error
	return account, err
}

func (dao *AccountDAO) FindByEmail(ctx context.Context, email string) (Account, error) {
	var account Account
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&account).Error
	return account, err
}
