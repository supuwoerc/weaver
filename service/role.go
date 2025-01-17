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

func (r *RoleService) CreateRole(ctx context.Context, name string, userIds, permissionIds []uint) error {
	// TODO:检查用户信息/角色信息/记录信息到用户时间线
	return r.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询有效的用户
		// 查询有效的角色
		// 创建角色 & 建立关联关系
		if err := r.repository.Create(ctx, name, nil, nil); err != nil {
			return err
		}
		return nil
	})
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

func (r *RoleService) GetRoleList(ctx context.Context, name string, limit, offset int) ([]*models.Role, int64, error) {
	return r.repository.GetRoleList(ctx, name, limit, offset)
}
