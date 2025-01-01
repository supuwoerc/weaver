package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/cache"
	"gin-web/repository/dao"
	"gin-web/repository/transducer"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		dao:   dao.NewUserDAO(),
		cache: cache.NewUserCache(),
	}
}

func toModelUser(u *dao.User) *models.User {
	user := models.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Nickname: transducer.NullString(u.Nickname),
		Gender:   transducer.NullValue(u.Gender),
		About:    transducer.NullString(u.About),
		Birthday: transducer.NullTime(u.Birthday),
		Roles:    toModelRoles(u.Roles),
	}
	return &user
}

func (u *UserRepository) Create(ctx context.Context, user models.User) error {
	return u.dao.Insert(ctx, &dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.dao.FindByEmail(ctx, email)
	return toModelUser(user), err
}

func (u *UserRepository) CacheTokenPair(ctx context.Context, email string, pair *models.TokenPair) error {
	return u.cache.HSetTokenPair(ctx, email, pair)
}

func (u *UserRepository) TokenPairExist(ctx context.Context, email string) (bool, error) {
	return u.cache.HExistsTokenPair(ctx, email)
}

func (u *UserRepository) AssociateRoles(ctx context.Context, uid uint, roleIds []uint) error {
	var roles []dao.Role
	for _, id := range roleIds {
		roles = append(roles, dao.Role{
			Model: gorm.Model{
				ID: id,
			},
		})
	}
	return u.dao.AssociateRoles(ctx, uid, &roles)
}

func (u *UserRepository) FindByUid(ctx context.Context, uid uint, needRoles bool) (*models.User, error) {
	user, err := u.dao.FindByUid(ctx, uid, needRoles)
	return toModelUser(user), err
}

func (u *UserRepository) FindRolesByUid(ctx context.Context, uid uint) ([]*models.Role, error) {
	roles, err := u.dao.FindRolesByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	result := lo.Map[*dao.Role, *models.Role](roles, func(item *dao.Role, _ int) *models.Role {
		return toModelRole(item)
	})
	return result, nil
}

func (u *UserRepository) GetTokenPair(ctx context.Context, email string) (*models.TokenPair, error) {
	return u.cache.HGetTokenPair(ctx, email)
}
