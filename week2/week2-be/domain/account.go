package domain

import "time"

type Account struct {
	Id          int64
	Email       string
	Password    string
	Name        string
	Birth       string
	About       string
	CreatedTime time.Time
	UpdatedTime time.Time
}
