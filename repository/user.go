package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/cache"
	"gin-web/repository/dao"
	"sync"
)

var (
	userRepository     *UserRepository
	userRepositoryOnce sync.Once
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository() *UserRepository {
	userRepositoryOnce.Do(func() {
		userRepository = &UserRepository{
			dao:   dao.NewUserDAO(),
			cache: cache.NewUserCache(),
		}
	})
	return userRepository
}

func (u *UserRepository) Create(ctx context.Context, user *models.User) error {
	return u.dao.Create(ctx, user)
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string, needAvatar, needRoles, needPermissions bool) (*models.User, error) {
	return u.dao.GetByEmail(ctx, email, needAvatar, needRoles, needPermissions)
}

func (u *UserRepository) CacheTokenPair(ctx context.Context, email string, pair *models.TokenPair) error {
	return u.cache.CacheTokenPair(ctx, email, pair)
}

func (u *UserRepository) GetTokenPairIsExist(ctx context.Context, email string) (bool, error) {
	return u.cache.GetTokenPairIsExist(ctx, email)
}

func (u *UserRepository) GetById(ctx context.Context, uid uint, needAvatar, needRoles, needPermissions bool) (*models.User, error) {
	return u.dao.GetById(ctx, uid, needAvatar, needRoles, needPermissions)
}

func (u *UserRepository) GetByIds(ctx context.Context, ids []uint, needRoles, needPermissions bool) ([]*models.User, error) {
	return u.dao.GetByIds(ctx, ids, needRoles, needPermissions)
}

func (u *UserRepository) GetTokenPair(ctx context.Context, email string) (*models.TokenPair, error) {
	return u.cache.GetTokenPair(ctx, email)
}
