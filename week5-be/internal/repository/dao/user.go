package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id       int64          `gorm:"primarykey,autoincrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string

	Nickname string `gorm:"type=varchar(128)"`
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`

	Phone sql.NullString `gorm:"unique"`

	// UTC Timestamp
	// created time
	Ctime int64
	// updated time
	Utime int64
}

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdateById(ctx context.Context, entity User) error
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type GORMUserDAO struct {
	db *gorm.DB
}

func (dao *GORMUserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	var mysqlError *mysql.MySQLError
	if errors.As(err, &mysqlError) {
		if mysqlError.Number == 1062 {
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&user).Error
	return user, err
}

func (dao *GORMUserDAO) UpdateById(ctx context.Context, entity User) error {
	return dao.db.WithContext(ctx).
		Model(&entity).
		Where("id=?", entity.Id).
		Updates(map[string]any{
			"utime":    time.Now().UnixMilli(),
			"nickname": entity.Nickname,
			"birthday": entity.Birthday,
			"about_me": entity.AboutMe,
		}).Error
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&user).Error
	return user, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&user).Error
	return user, err
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}
