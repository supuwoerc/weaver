package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
	"github.com/samber/lo"
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

func toModelRole(role *dao.Role) *models.Role {
	return &models.Role{
		ID:   role.ID,
		Name: role.Name,
	}
}

func toModelRoles(roles []*dao.Role) []*models.Role {
	return lo.Map(roles, func(item *dao.Role, index int) *models.Role {
		return toModelRole(item)
	})
}

func (r *RoleRepository) Create(ctx context.Context, name string) error {
	return r.dao.Insert(ctx, &dao.Role{
		Name: name,
	})
}

func (r *RoleRepository) GetRolesByIds(ctx context.Context, ids []uint) ([]*models.Role, error) {
	ret, err := r.dao.GetRolesByIds(ctx, ids)
	return toModelRoles(ret), err
}
