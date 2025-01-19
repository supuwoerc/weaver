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

func (r *PermissionService) CreatePermission(ctx context.Context, name, resource string, roleIds []uint) error {
	// 查询有效的角色
	roles, err := r.roleRepository.GetByIds(ctx, roleIds, false, false)
	if err != nil {
		return err
	}
	// 创建权限 & 建立关联关系
	return r.permissionRepository.Create(ctx, name, resource, roles)
}

func (r *PermissionService) GetPermissionList(ctx context.Context, keyword string, limit, offset int) ([]*models.Permission, int64, error) {
	return r.permissionRepository.GetList(ctx, keyword, limit, offset)
}

func (r *PermissionService) GetPermissionDetail(ctx context.Context, id uint) (*models.Permission, error) {
	return r.permissionRepository.GetById(ctx, id, true)
}

func (r *PermissionService) UpdatePermission(ctx context.Context, id uint, name, resource string, roleIds []uint) error {
	// 查询有效的角色
	roles, err := r.roleRepository.GetByIds(ctx, roleIds, false, false)
	if err != nil {
		return err
	}
	// 更新权限
	return r.permissionRepository.Update(ctx, id, name, resource, roles)
}

func (r *PermissionService) DeletePermission(ctx context.Context, id uint) error {
	count := r.permissionRepository.GetRolesCount(ctx, id)
	if count > 0 {
		return response.PermissionExistRoleRef
	}
	return r.permissionRepository.DeleteById(ctx, id)
}
