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
	dao RoleDAO
}

func NewRoleRepository() *RoleRepository {
	roleRepositoryOnce.Do(func() {
		roleRepository = &RoleRepository{
			dao: dao.NewRoleDAO(),
		}
	})
	return roleRepository
}

func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	return r.dao.Create(ctx, role)
}

func (r *RoleRepository) GetByIds(ctx context.Context, ids []uint, needUsers, needPermissions bool) ([]*models.Role, error) {
	return r.dao.GetByIds(ctx, ids, needUsers, needPermissions)
}

func (r *RoleRepository) GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Role, int64, error) {
	return r.dao.GetList(ctx, keyword, limit, offset)
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	return r.dao.GetByName(ctx, name)
}

func (r *RoleRepository) GetById(ctx context.Context, id uint, needUsers, needPermissions bool) (*models.Role, error) {
	return r.dao.GetById(ctx, id, needUsers, needPermissions)
}

func (r *RoleRepository) Update(ctx context.Context, role *models.Role) error {
	return r.dao.Update(ctx, role)
}

func (r *RoleRepository) AssociateUsers(ctx context.Context, id uint, users []*models.User) error {
	return r.dao.AssociateUsers(ctx, id, users)
}

func (r *RoleRepository) AssociatePermissions(ctx context.Context, id uint, permissions []*models.Permission) error {
	return r.dao.AssociatePermissions(ctx, id, permissions)
}

func (r *RoleRepository) DeleteById(ctx context.Context, id, updater uint) error {
	return r.dao.DeleteById(ctx, id, updater)
}

func (r *RoleRepository) GetUsersCount(ctx context.Context, id uint) int64 {
	return r.dao.GetUsersCount(ctx, id)
}

func (r *RoleRepository) GetPermissionsCount(ctx context.Context, id uint) int64 {
	return r.dao.GetPermissionsCount(ctx, id)
}
