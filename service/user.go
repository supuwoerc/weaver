package service

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/constant"
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
	password, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var pwd = string(password)
	user.Password = &pwd
	return u.repository.Create(context, user)
}

func (u *UserService) Login(ctx context.Context, email string, password string) (models.User, error) {
	user, err := u.repository.FindByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil {
		return models.User{}, constant.USER_LOGIN_FAIL_ERR
	}
	return user, nil
}
