package repository

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/redis"
	"gin-web/repository/cache"
	"gin-web/repository/dao"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	userRepository     *UserRepository
	userRepositoryOnce sync.Once
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
	dao   UserDAO
	cache UserCache
}

func NewUserRepository(db *gorm.DB, r *redis.CommonRedisClient) *UserRepository {
	userRepositoryOnce.Do(func() {
		userRepository = &UserRepository{
			dao:   dao.NewUserDAO(db),
			cache: cache.NewUserCache(r),
		}
	})
	return userRepository
}

func (u *UserRepository) Create(ctx context.Context, user *models.User) error {
	return u.dao.Create(ctx, user)
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string, needRoles, needPermissions, needDepts bool) (*models.User, error) {
	return u.dao.GetByEmail(ctx, email, needRoles, needPermissions, needDepts)
}

func (u *UserRepository) CacheTokenPair(ctx context.Context, email string, pair *models.TokenPair) error {
	return u.cache.CacheTokenPair(ctx, email, pair)
}

func (u *UserRepository) GetTokenPairIsExist(ctx context.Context, email string) (bool, error) {
	return u.cache.GetTokenPairIsExist(ctx, email)
}

func (u *UserRepository) GetById(ctx context.Context, uid uint, needRoles, needPermissions, needDepts bool) (*models.User, error) {
	return u.dao.GetById(ctx, uid, needRoles, needPermissions, needDepts)
}

func (u *UserRepository) GetByIds(ctx context.Context, ids []uint, needRoles, needPermissions, needDepts bool) ([]*models.User, error) {
	return u.dao.GetByIds(ctx, ids, needRoles, needPermissions, needDepts)
}

func (u *UserRepository) GetTokenPair(ctx context.Context, email string) (*models.TokenPair, error) {
	return u.cache.GetTokenPair(ctx, email)
}

func (u *UserRepository) GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.User, int64, error) {
	return u.dao.GetList(ctx, keyword, limit, offset)
}

func (u *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	return u.dao.GetAll(ctx)
}

func (u *UserRepository) CacheActiveAccountCode(ctx context.Context, id uint, code string, duration time.Duration) error {
	return u.cache.CacheActiveAccountCode(ctx, id, code, duration)
}

func (u *UserRepository) GetActiveAccountCode(ctx context.Context, id uint) (string, error) {
	return u.cache.GetActiveAccountCode(ctx, id)
}

func (u *UserRepository) RemoveActiveAccountCode(ctx context.Context, id uint) error {
	return u.cache.RemoveActiveAccountCode(ctx, id)
}

func (u *UserRepository) UpdateAccountStatus(ctx context.Context, id uint, status constant.UserStatus) error {
	return u.dao.UpdateAccountStatus(ctx, id, status)
}
