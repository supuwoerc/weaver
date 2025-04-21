package repository

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/constant"
	"time"
)

type UserDAO interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string, needRoles, needPermissions, needDepts bool) (*models.User, error)
	GetById(ctx context.Context, uid uint, needRoles, needPermissions, needDepts bool) (*models.User, error)
	GetByIds(ctx context.Context, ids []uint, needRoles, needPermissions, needDepts bool) ([]*models.User, error)
	GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.User, int64, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	UpdateAccountStatus(ctx context.Context, id uint, status constant.UserStatus) error
}
type UserCache interface {
	CacheTokenPair(ctx context.Context, email string, pair *models.TokenPair) error
	GetTokenPairIsExist(ctx context.Context, email string) (bool, error)
	GetTokenPair(ctx context.Context, email string) (*models.TokenPair, error)
	CacheActiveAccountCode(ctx context.Context, id uint, code string, duration time.Duration) error
	GetActiveAccountCode(ctx context.Context, id uint) (string, error)
	RemoveActiveAccountCode(ctx context.Context, id uint) error
}
type UserRepository struct {
	UserDAO
	UserCache
}

func NewUserRepository(dao UserDAO, cache UserCache) *UserRepository {
	return &UserRepository{
		UserDAO:   dao,
		UserCache: cache,
	}
}
