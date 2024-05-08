package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
	"time"
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

func toModelUser(u dao.User) models.User {
	return models.User{
		Email:    u.Email,
		Password: u.Password,
		NickName: u.NickName,
		Gender:   models.UserGender(u.Gender),
		About:    u.About,
		Birthday: time.UnixMilli(u.Birthday),
	}
}

func (u UserRepository) Create(ctx context.Context, user models.User) error {
	return u.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (u UserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := u.dao.FindByEmail(ctx, email)
	return toModelUser(user), err
}
