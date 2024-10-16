package domain

import "time"

type Account struct {
	Id          int64     `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Name        string    `json:"name"`
	Birth       string    `json:"birth"`
	About       string    `json:"about"`
	CreatedTime time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
}
