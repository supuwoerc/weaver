package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/cache"
	"gin-web/repository/dao"
	"gorm.io/gorm"
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
	return u.dao.Insert(ctx, user)
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.dao.FindByEmail(ctx, email)
}

func (u *UserRepository) CacheTokenPair(ctx context.Context, email string, pair *models.TokenPair) error {
	return u.cache.HSetTokenPair(ctx, email, pair)
}

func (u *UserRepository) TokenPairExist(ctx context.Context, email string) (bool, error) {
	return u.cache.HExistsTokenPair(ctx, email)
}

func (u *UserRepository) AssociateRoles(ctx context.Context, uid uint, roleIds []uint) error {
	var roles []*models.Role
	for _, id := range roleIds {
		roles = append(roles, &models.Role{
			Model: gorm.Model{
				ID: id,
			},
		})
	}
	return u.dao.AssociateRoles(ctx, uid, roles)
}

func (u *UserRepository) FindByUid(ctx context.Context, uid uint, needRoles, needPermissions bool) (*models.User, error) {
	return u.dao.FindByUid(ctx, uid, needRoles, needPermissions)
}

func (u *UserRepository) FindRolesByUid(ctx context.Context, uid uint) ([]*models.Role, error) {
	return u.dao.FindRolesByUid(ctx, uid)
}

func (u *UserRepository) GetTokenPair(ctx context.Context, email string) (*models.TokenPair, error) {
	return u.cache.HGetTokenPair(ctx, email)
}
