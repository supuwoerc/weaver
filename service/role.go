package service

import (
	"context"
	"gin-web/models"
	"gin-web/repository"
	"github.com/samber/lo"
	"sync"
)

type RoleService struct {
	*BasicService
	repository *repository.RoleRepository
}

var (
	roleOnce    sync.Once
	roleService *RoleService
)

func NewRoleService() *RoleService {
	roleOnce.Do(func() {
		roleService = &RoleService{
			BasicService: NewBasicService(),
			repository:   repository.NewRoleRepository(),
		}
	})
	return roleService
}

func (r *RoleService) CreateRole(ctx context.Context, name string) error {
	return r.repository.Create(ctx, name)
}

func (r *RoleService) FilterValidRoles(ctx context.Context, roleIds []uint) ([]uint, error) {
	roles, err := r.repository.GetRolesByIds(ctx, roleIds)
	if err != nil {
		return []uint{}, err
	}
	validIds := lo.Map[*models.Role, uint](roles, func(item *models.Role, _ int) uint {
		return item.ID
	})
	result := lo.Filter(roleIds, func(item uint, _ int) bool {
		return lo.Contains(validIds, item)
	})
	return result, nil
}

func (r *RoleService) GetRoleList(ctx context.Context, name string, limit, offset int) ([]*models.Role, error) {
	return r.repository.GetRoleList(ctx, name, limit, offset)
}
