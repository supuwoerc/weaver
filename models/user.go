package models

type UserGender int8

const (
	GenderMale UserGender = iota + 1
	GenderFemale
)

type User struct {
	ID       uint       `json:"id"`
	Email    string     `json:"email"`
	Password *string    `json:"password,omitempty"`
	Nickname string     `json:"nickname"`
	Gender   UserGender `json:"gender"`
	About    string     `json:"about"`
	Birthday string     `json:"birthday"`
	Roles    []*Role    `json:"roles"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
