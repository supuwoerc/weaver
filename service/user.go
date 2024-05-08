package service

import (
	"context"
	"gin-web/models"
	"gin-web/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*BasicService
	repository *repository.UserRepository
}

var userService *UserService

func NewUserService() *UserService {
	if userService == nil {
		userService = &UserService{
			BasicService: NewBasicService(),
			repository:   repository.NewUserRepository(),
		}
	}
	return userService
}

func (u *UserService) SignUp(context context.Context, user models.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)
	return u.repository.Create(context, user)
}
