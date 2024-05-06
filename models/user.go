package models

import "time"

type UserGender int

const (
	GENDER_UNKNOWN UserGender = iota
	GENDER_MALE
	GENDER_FEMALE
)

type User struct {
	Email    string
	Password string
	NickName string
	Gender   UserGender
	About    string
	Birthday time.Time
}
