package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/cache"
	"gin-web/repository/dao"
	"time"
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

var userRepository *UserRepository

func NewUserRepository() *UserRepository {
	if userRepository == nil {
		userRepository = &UserRepository{
			dao:   dao.NewUserDAO(),
			cache: cache.NewUserCache(),
		}
	}
	return userRepository
}

func toModelUser(u dao.User) models.User {
	return models.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: &u.Password,
		NickName: u.NickName,
		Gender:   models.UserGender(u.Gender),
		About:    u.About,
		Birthday: time.UnixMilli(u.Birthday).Format(time.DateTime),
	}
}

func (u UserRepository) Create(ctx context.Context, user models.User) error {
	return u.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: *user.Password,
	})
}

func (u UserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := u.dao.FindByEmail(ctx, email)
	return toModelUser(user), err
}

func (u UserRepository) CacheTokenPair(ctx context.Context, email string, pair *models.TokenPair) error {
	return u.cache.HSetTokenPair(ctx, email, pair)
}

func (u UserRepository) TokenPairExist(ctx context.Context, email string) (bool, error) {
	return u.cache.HExistsTokenPair(ctx, email)
}

func (u UserRepository) DelTokenPair(ctx context.Context, email string) error {
	return u.cache.HDelTokenPair(ctx, email)
}
