package repository

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/database"
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

func (r *PermissionRepository) GetById(ctx context.Context, id uint, needRoles bool) (*models.Permission, error) {
	return r.dao.GetById(ctx, id, needRoles)
}

func (r *PermissionRepository) GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Permission, int64, error) {
	return r.dao.GetList(ctx, keyword, limit, offset)
}

func (r *PermissionRepository) DeleteById(ctx context.Context, id uint) error {
	return r.dao.DeleteById(ctx, id)
}

func (r *PermissionRepository) GetRolesCount(ctx context.Context, id uint) int64 {
	return r.dao.GetRolesCount(ctx, id)
}

func (r *PermissionRepository) Update(ctx context.Context, id uint, name, resource string) error {
	return r.dao.Update(ctx, &models.Permission{
		Name:     name,
		Resource: resource,
		BasicModel: database.BasicModel{
			ID: id,
		},
	})
}

func (r *PermissionRepository) AssociateRoles(ctx context.Context, id uint, roles []*models.Role) error {
	return r.dao.AssociateRoles(ctx, id, roles)
}
