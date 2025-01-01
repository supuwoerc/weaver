package models

import "time"

type User struct {
	ID       uint      `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Nickname *string   `json:"nickname"`
	Gender   *uint8    `json:"gender"`
	About    *string   `json:"about"`
	Birthday time.Time `json:"birthday"`
	Roles    []*Role   `json:"roles"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
