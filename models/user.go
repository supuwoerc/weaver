package models

type UserGender int

const (
	GENDER_UNKNOWN UserGender = iota
	GENDER_MALE
	GENDER_FEMALE
)

type User struct {
	ID       uint       `json:"id"`
	Email    string     `json:"email"`
	Password *string    `json:"password,omitempty"`
	NickName string     `json:"nick_name"`
	Gender   UserGender `json:"gender"`
	About    string     `json:"about"`
	Birthday string     `json:"birthday"`
	Roles    []*Role    `json:"roles"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
