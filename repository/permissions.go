package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
	"sync"
)

var (
	permissionRepository     *PermissionRepository
	permissionRepositoryOnce sync.Once
)

type PermissionRepository struct {
	dao *dao.PermissionDAO
}

func NewPermissionRepository() *PermissionRepository {
	permissionRepositoryOnce.Do(func() {
		permissionRepository = &PermissionRepository{
			dao: dao.NewPermissionDAO(),
		}
	})
	return permissionRepository
}

func (r *PermissionRepository) Create(ctx context.Context, name, resource string, roles []*models.Role) error {
	return r.dao.Create(ctx, &models.Permission{
		Name:     name,
		Resource: resource,
		Roles:    roles,
	})
}
func (r *PermissionRepository) GetByIds(ctx context.Context, ids []uint, needRoles bool) ([]*models.Permission, error) {
	return r.dao.GetByIds(ctx, ids, needRoles)
}
