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
	PermissionDAO
}

func NewPermissionRepository(dao PermissionDAO) *PermissionRepository {
	return &PermissionRepository{
		PermissionDAO: dao,
	}
}
