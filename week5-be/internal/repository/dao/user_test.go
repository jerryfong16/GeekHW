package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string

		build func(t *testing.T) *sql.DB

		ctx  context.Context
		user User

		expectedErr error
	}{
		{
			name: "插入成功",
			build: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)

				mockRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(mockRes)

				return db
			},
			ctx:         context.Background(),
			user:        User{Nickname: "tom"},
			expectedErr: nil,
		},
		{
			name: "邮箱重复",
			build: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.
					ExpectExec("INSERT INTO .*").
					WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx:         context.Background(),
			user:        User{Nickname: "tom"},
			expectedErr: ErrDuplicateEmail,
		},
		{
			name: "DB异常",
			build: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.
					ExpectExec("INSERT INTO .*").
					WillReturnError(errors.New("DB异常"))
				return db
			},
			ctx:         context.Background(),
			user:        User{Nickname: "tom"},
			expectedErr: errors.New("DB异常"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sqlDB := testCase.build(t)
			db, err := gorm.Open(
				mysql.New(mysql.Config{
					Conn:                      sqlDB,
					SkipInitializeWithVersion: true,
				}),
				&gorm.Config{
					DisableAutomaticPing:   true,
					SkipDefaultTransaction: true,
				},
			)
			assert.NoError(t, err)
			dao := NewUserDAO(db)
			insertErr := dao.Insert(testCase.ctx, testCase.user)
			assert.Equal(t, insertErr, testCase.expectedErr)
		})
	}
}
