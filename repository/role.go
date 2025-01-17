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
	return r.dao.Insert(ctx, &models.Role{
		Name:        name,
		Users:       users,
		Permissions: permissions,
	})
}

func (r *RoleRepository) GetRolesByIds(ctx context.Context, ids []uint) ([]*models.Role, error) {
	return r.dao.GetRolesByIds(ctx, ids)
}

func (r *RoleRepository) GetRoleList(ctx context.Context, name string, limit, offset int) ([]*models.Role, int64, error) {
	return r.dao.GetRoleList(ctx, name, limit, offset)
}
