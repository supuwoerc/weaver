package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
	"sync"
)

var (
	roleRepository     *RoleRepository
	roleRepositoryOnce sync.Once
)

type RoleRepository struct {
	dao *dao.RoleDAO
}

func NewRoleRepository() *RoleRepository {
	roleRepositoryOnce.Do(func() {
		roleRepository = &RoleRepository{
			dao: dao.NewRoleDAO(),
		}
	})
	return roleRepository
}

func (r *RoleRepository) Create(ctx context.Context, name string, users []*models.User, permissions []*models.Permission) error {
	return r.dao.Create(ctx, &models.Role{
		Name:        name,
		Users:       users,
		Permissions: permissions,
	})
}

func (r *RoleRepository) GetByIds(ctx context.Context, ids []uint, needUsers, needPermissions bool) ([]*models.Role, error) {
	return r.dao.GetByIds(ctx, ids, needUsers, needPermissions)
}

func (r *RoleRepository) GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Role, int64, error) {
	return r.dao.GetList(ctx, keyword, limit, offset)
}

func (r *RoleRepository) GetIsExistByName(ctx context.Context, name string) (bool, error) {
	return r.dao.GetIsExistByName(ctx, name)
}
