package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string

	Nickname string
	Birthday time.Time
	AboutMe  string

	Phone string

	Ctime time.Time
}
