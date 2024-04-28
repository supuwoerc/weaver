package models

import "gorm.io/gorm"

type UserGender int

const (
	GENDER_UNKNOWN UserGender = iota
	GENDER_MALE
	GENDER_FEMALE
)

type User struct {
	gorm.Model
	Name   string     `gorm:"comment:用户名"`
	Gender UserGender `gorm:"comment:用户名"`
}

func (u User) TableName() string {
	return "user"
}
