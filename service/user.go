package service

import (
	"context"
	"gin-web/models"
	"gin-web/repository"
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
	// TODO:BCrypt加密
	return u.repository.Create(context, user)
}
