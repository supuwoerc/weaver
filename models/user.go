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
	Name     string     `gorm:"type:varchar(20);comment:用户名"`
	Password string     `gorm:"type:varchar(255);comment:密码"`
	Gender   UserGender `gorm:"type:integer;comment:性别;default:0"`
}

func (u User) TableName() string {
	return "user"
}
