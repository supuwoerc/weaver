package repository

import (
	"context"
	"gin-web/models"
)

type RoleDAO interface {
	Create(ctx context.Context, role *models.Role) error
	GetByIds(ctx context.Context, ids []uint, needUsers, needPermissions bool) ([]*models.Role, error)
	GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Role, int64, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	GetById(ctx context.Context, id uint, needUsers, needPermissions bool) (*models.Role, error)
	Update(ctx context.Context, role *models.Role) error
	AssociateUsers(ctx context.Context, id uint, users []*models.User) error
	AssociatePermissions(ctx context.Context, id uint, permissions []*models.Permission) error
	DeleteById(ctx context.Context, id, updater uint) error
	GetUsersCount(ctx context.Context, id uint) int64
	GetPermissionsCount(ctx context.Context, id uint) int64
}
type RoleCache interface{}

type RoleRepository struct {
	RoleDAO
}

func NewRoleRepository(dao RoleDAO) *RoleRepository {
	return &RoleRepository{
		RoleDAO: dao,
	}
}
