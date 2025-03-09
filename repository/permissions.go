package repository

import (
	"context"
	"gin-web/models"
)

type PermissionDAO interface {
	Create(ctx context.Context, permission *models.Permission) error
	GetByIds(ctx context.Context, ids []uint, needRoles bool) ([]*models.Permission, error)
	GetById(ctx context.Context, id uint, needRoles bool) (*models.Permission, error)
	GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Permission, int64, error)
	DeleteById(ctx context.Context, id, updater uint) error
	GetRolesCount(ctx context.Context, id uint) int64
	Update(ctx context.Context, permission *models.Permission) error
	AssociateRoles(ctx context.Context, id uint, roles []*models.Role) error
	GetByNameOrResource(ctx context.Context, name, resource string) ([]*models.Permission, error)
}

type PermissionRepository struct {
	dao PermissionDAO
}

func NewPermissionRepository(dao PermissionDAO) *PermissionRepository {
	return &PermissionRepository{
		dao: dao,
	}
}

func (r *PermissionRepository) Create(ctx context.Context, permission *models.Permission) error {
	return r.dao.Create(ctx, permission)
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

func (r *PermissionRepository) DeleteById(ctx context.Context, id, updater uint) error {
	return r.dao.DeleteById(ctx, id, updater)
}

func (r *PermissionRepository) GetRolesCount(ctx context.Context, id uint) int64 {
	return r.dao.GetRolesCount(ctx, id)
}

func (r *PermissionRepository) Update(ctx context.Context, permission *models.Permission) error {
	return r.dao.Update(ctx, permission)
}

func (r *PermissionRepository) AssociateRoles(ctx context.Context, id uint, roles []*models.Role) error {
	return r.dao.AssociateRoles(ctx, id, roles)
}

func (r *PermissionRepository) GetByNameOrResource(ctx context.Context, name, resource string) ([]*models.Permission, error) {
	return r.dao.GetByNameOrResource(ctx, name, resource)
}
