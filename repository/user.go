package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

var userRepository *UserRepository

func NewUserRepository() *UserRepository {
	if userRepository == nil {
		userRepository = &UserRepository{
			dao: dao.NewUserDAO(),
		}
	}
	return userRepository
}

func (u UserRepository) Create(ctx context.Context, user models.User) error {
	return u.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}
