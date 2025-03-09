package service

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/database"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/samber/lo"
	"sync"
)

type PermissionRepository interface {
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

type PermissionService struct {
	*BasicService
	permissionRepository PermissionRepository
	roleRepository       RoleRepository
}

var (
	permissionOnce    sync.Once
	permissionService *PermissionService
)

// TODO:包装一个基础的参数结构体,wire注入这个Value
func NewPermissionService(basic *BasicService, permissionRepo PermissionRepository, roleRepository RoleRepository) *PermissionService {
	permissionOnce.Do(func() {
		permissionService = &PermissionService{
			BasicService:         basic,
			permissionRepository: permissionRepo,
			roleRepository:       roleRepository,
		}
	})
	return permissionService
}

func (p *PermissionService) lockPermissionField(ctx context.Context, name, resource string, roleIds []uint) ([]*utils.RedisLock, error) {
	locks := make([]*utils.RedisLock, 0, len(roleIds)+2)
	// 权限名称锁
	permissionNameLock := p.locksmith.NewLock(constant.PermissionNamePrefix, name)
	if err := permissionNameLock.Lock(ctx, true); err != nil {
		return locks, err
	}
	locks = append(locks, permissionNameLock)
	// 权限资源锁
	permissionResourceLock := p.locksmith.NewLock(constant.PermissionResourcePrefix, resource)
	if err := permissionResourceLock.Lock(ctx, true); err != nil {
		return locks, err
	}
	locks = append(locks, permissionResourceLock)
	// 角色锁
	for _, roleId := range roleIds {
		roleIdLock := p.locksmith.NewLock(constant.RoleIdPrefix, roleId)
		if err := roleIdLock.Lock(ctx, true); err != nil {
			return locks, err
		}
		locks = append(locks, roleIdLock)
	}
	return locks, nil
}

func (p *PermissionService) CreatePermission(ctx context.Context, operator uint, name, resource string, roleIds []uint) error {
	locks, err := p.lockPermissionField(ctx, name, resource, roleIds)
	defer func() {
		for _, l := range locks {
			if e := l.Unlock(); e != nil {
				p.logger.Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	return p.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existPermissions, temp := p.permissionRepository.GetByNameOrResource(ctx, name, resource)
		if temp != nil {
			return temp
		}
		if len(existPermissions) > 0 {
			return response.PermissionCreateDuplicate
		}
		// 查询有效的角色
		var roles []*models.Role
		if len(roleIds) > 0 {
			roles, err = p.roleRepository.GetByIds(ctx, roleIds, false, false)
			if err != nil {
				return err
			}
		}
		// 创建权限 & 建立关联关系
		return p.permissionRepository.Create(ctx, &models.Permission{
			Name:      name,
			Resource:  resource,
			Roles:     roles,
			CreatorId: operator,
			UpdaterId: operator,
		})
	})
}

func (p *PermissionService) GetPermissionList(ctx context.Context, keyword string, limit, offset int) ([]*response.PermissionListRowResponse, int64, error) {
	list, total, err := p.permissionRepository.GetList(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return lo.Map(list, func(item *models.Permission, _ int) *response.PermissionListRowResponse {
		return response.ToPermissionListRowResponse(item)
	}), total, nil
}

func (p *PermissionService) GetPermissionDetail(ctx context.Context, id uint) (*response.PermissionDetailResponse, error) {
	permission, err := p.permissionRepository.GetById(ctx, id, true)
	if err != nil {
		return nil, err
	}
	return response.ToPermissionDetailResponse(permission), nil
}

func (p *PermissionService) UpdatePermission(ctx context.Context, operator uint, id uint, name, resource string, roleIds []uint) error {
	// 对权限自身加锁
	permissionLock := p.locksmith.NewLock(constant.PermissionIdPrefix, id)
	if err := permissionLock.Lock(ctx, true); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := lock.Unlock(); e != nil {
			p.logger.Errorf("unlock fail %s", e.Error())
		}
	}(permissionLock)
	// 对 name & resource & roleIds 加锁
	locks, err := p.lockPermissionField(ctx, name, resource, roleIds)
	defer func() {
		for _, l := range locks {
			if e := l.Unlock(); e != nil {
				p.logger.Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	return p.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existPermissions, temp := p.permissionRepository.GetByNameOrResource(ctx, name, resource)
		if temp != nil {
			return temp
		}
		count := len(existPermissions)
		if count > 1 || (count == 1 && existPermissions[0].ID != id) {
			return response.PermissionCreateDuplicate
		}
		// 更新权限
		err = p.permissionRepository.Update(ctx, &models.Permission{
			Name:      name,
			Resource:  resource,
			UpdaterId: operator,
			BasicModel: database.BasicModel{
				ID: id,
			},
		})
		if err != nil {
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

func (p *PermissionService) DeletePermission(ctx context.Context, id, operator uint) error {
	// 对权限自身加锁
	permissionLock := p.locksmith.NewLock(constant.PermissionIdPrefix, id)
	if err := permissionLock.Lock(ctx, true); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := lock.Unlock(); e != nil {
			p.logger.Errorf("unlock fail %s", e.Error())
		}
	}(permissionLock)
	count := p.permissionRepository.GetRolesCount(ctx, id)
	if count > 0 {
		return response.PermissionExistRoleRef
	}
	return p.permissionRepository.DeleteById(ctx, id, operator)
}
