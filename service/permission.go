package service

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/response"
	"gin-web/repository"
	"sync"
)

type PermissionService struct {
	*BasicService
	permissionRepository *repository.PermissionRepository
	roleRepository       *repository.RoleRepository
}

var (
	permissionOnce    sync.Once
	permissionService *PermissionService
)

func NewPermissionService() *PermissionService {
	permissionOnce.Do(func() {
		permissionService = &PermissionService{
			BasicService:         NewBasicService(),
			permissionRepository: repository.NewPermissionRepository(),
			roleRepository:       repository.NewRoleRepository(),
		}
	})
	return permissionService
}

func (p *PermissionService) CreatePermission(ctx context.Context, name, resource string, roleIds []uint) error {
	return p.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询有效的角色
		var roles []*models.Role
		var err error
		if len(roleIds) > 0 {
			// TODO:加锁
			roles, err = p.roleRepository.GetByIds(ctx, roleIds, false, false)
			if err != nil {
				return err
			}
		}
		// 创建权限 & 建立关联关系
		return p.permissionRepository.Create(ctx, name, resource, roles)
	})
}

func (p *PermissionService) GetPermissionList(ctx context.Context, keyword string, limit, offset int) ([]*models.Permission, int64, error) {
	return p.permissionRepository.GetList(ctx, keyword, limit, offset)
}

func (p *PermissionService) GetPermissionDetail(ctx context.Context, id uint) (*models.Permission, error) {
	return p.permissionRepository.GetById(ctx, id, true)
}

func (p *PermissionService) UpdatePermission(ctx context.Context, id uint, name, resource string, roleIds []uint) error {
	return p.Transaction(ctx, false, func(ctx context.Context) error {
		// 更新权限
		if err := p.permissionRepository.Update(ctx, id, name, resource); err != nil {
			return err
		}
		// 查询有效的角色
		var roles []*models.Role
		if len(roleIds) > 0 {
			r, err := p.roleRepository.GetByIds(ctx, roleIds, false, false)
			if err != nil {
				return err
			}
			roles = r
		}
		// 更新关联关系
		return p.permissionRepository.AssociateRoles(ctx, id, roles)
	})
}

func (p *PermissionService) DeletePermission(ctx context.Context, id uint) error {
	count := p.permissionRepository.GetRolesCount(ctx, id)
	if count > 0 {
		return response.PermissionExistRoleRef
	}
	return p.permissionRepository.DeleteById(ctx, id)
}
