package service

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
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

func lockPermissionField(ctx context.Context, name, resource string, roleIds []uint) ([]*utils.RedisLock, error) {
	locks := make([]*utils.RedisLock, 0, len(roleIds)+2)
	// 权限名称锁
	permissionNameLock := utils.NewLock(constant.PermissionNamePrefix, name)
	if err := utils.Lock(ctx, permissionNameLock); err != nil {
		return locks, err
	}
	locks = append(locks, permissionNameLock)
	// 权限资源锁
	permissionResourceLock := utils.NewLock(constant.PermissionResourcePrefix, resource)
	if err := utils.Lock(ctx, permissionResourceLock); err != nil {
		return locks, err
	}
	locks = append(locks, permissionResourceLock)
	// 角色锁
	for _, roleId := range roleIds {
		roleIdLock := utils.NewLock(constant.RoleIdPrefix, roleId)
		if err := utils.Lock(ctx, roleIdLock); err != nil {
			return locks, err
		}
		locks = append(locks, roleIdLock)
	}
	return locks, nil
}

func (p *PermissionService) CreatePermission(ctx context.Context, name, resource string, roleIds []uint) error {
	locks, err := lockPermissionField(ctx, name, resource, roleIds)
	defer func() {
		for _, l := range locks {
			if e := utils.Unlock(l); e != nil {
				global.Logger.Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	return p.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询有效的角色
		var roles []*models.Role
		if len(roleIds) > 0 {
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
	// 对权限自身加锁
	permissionLock := utils.NewLock(constant.PermissionIdPrefix, id)
	if err := utils.Lock(ctx, permissionLock); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := utils.Unlock(lock); e != nil {
			global.Logger.Errorf("unlock fail %s", e.Error())
		}
	}(permissionLock)
	// 对 name & resource & roleIds 加锁
	locks, err := lockPermissionField(ctx, name, resource, roleIds)
	defer func() {
		for _, l := range locks {
			if e := utils.Unlock(l); e != nil {
				global.Logger.Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	return p.Transaction(ctx, false, func(ctx context.Context) error {
		// 更新权限
		if err = p.permissionRepository.Update(ctx, id, name, resource); err != nil {
			return err
		}
		// 查询有效的角色
		var roles []*models.Role
		if len(roleIds) > 0 {
			roles, err = p.roleRepository.GetByIds(ctx, roleIds, false, false)
			if err != nil {
				return err
			}
		}
		// 更新关联关系
		return p.permissionRepository.AssociateRoles(ctx, id, roles)
	})
}

func (p *PermissionService) DeletePermission(ctx context.Context, id uint) error {
	// 对权限自身加锁
	permissionLock := utils.NewLock(constant.PermissionIdPrefix, id)
	if err := utils.Lock(ctx, permissionLock); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := utils.Unlock(lock); e != nil {
			global.Logger.Errorf("unlock fail %s", e.Error())
		}
	}(permissionLock)
	count := p.permissionRepository.GetRolesCount(ctx, id)
	if count > 0 {
		return response.PermissionExistRoleRef
	}
	return p.permissionRepository.DeleteById(ctx, id)
}
