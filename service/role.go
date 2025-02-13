package service

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/database"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/repository"
	"github.com/samber/lo"
	"sync"
)

type RoleService struct {
	*BasicService
	roleRepository       *repository.RoleRepository
	userRepository       *repository.UserRepository
	permissionRepository *repository.PermissionRepository
}

var (
	roleOnce    sync.Once
	roleService *RoleService
)

func NewRoleService() *RoleService {
	roleOnce.Do(func() {
		roleService = &RoleService{
			BasicService:         NewBasicService(),
			roleRepository:       repository.NewRoleRepository(),
			userRepository:       repository.NewUserRepository(),
			permissionRepository: repository.NewPermissionRepository(),
		}
	})
	return roleService
}

func lockRoleField(ctx context.Context, name string, permissionIds []uint) ([]*utils.RedisLock, error) {
	locks := make([]*utils.RedisLock, 0, len(permissionIds)+1)
	// 角色名称锁
	roleNameLock := utils.NewLock(constant.RoleNamePrefix, name)
	if err := utils.Lock(ctx, roleNameLock); err != nil {
		return locks, err
	}
	locks = append(locks, roleNameLock)
	// 角色权限锁
	for _, permissionId := range permissionIds {
		roleIdLock := utils.NewLock(constant.PermissionIdPrefix, permissionId)
		if err := utils.Lock(ctx, roleIdLock); err != nil {
			return locks, err
		}
		locks = append(locks, roleIdLock)
	}
	return locks, nil
}

func (r *RoleService) CreateRole(ctx context.Context, operator uint, name string, userIds, permissionIds []uint) error {
	locks, err := lockRoleField(ctx, name, permissionIds)
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
	// TODO:记录信息到用户时间线
	return r.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existRole, temp := r.roleRepository.GetByName(ctx, name)
		if temp != nil && !errors.Is(temp, response.RoleNotExist) {
			return temp
		}
		if existRole != nil {
			return response.RoleCreateDuplicateName
		}
		// 查询有效的用户
		var users []*models.User
		if len(userIds) > 0 {
			users, err = r.userRepository.GetByIds(ctx, userIds, false, false, false)
			if err != nil {
				return err
			}
		}
		// 查询有效的权限
		var permissions []*models.Permission
		if len(permissionIds) > 0 {
			permissions, err = r.permissionRepository.GetByIds(ctx, permissionIds, false)
			if err != nil {
				return err
			}
		}
		// 创建角色 & 建立关联关系
		return r.roleRepository.Create(ctx, &models.Role{
			Name:        name,
			Users:       users,
			Permissions: permissions,
			CreatorId:   operator,
			UpdaterId:   operator,
		})
	})
}

func (r *RoleService) GetRoleList(ctx context.Context, keyword string, limit, offset int) ([]*response.RoleListRowResponse, int64, error) {
	list, total, err := r.roleRepository.GetList(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return lo.Map(list, func(item *models.Role, _ int) *response.RoleListRowResponse {
		return response.ToRoleListRowResponse(item)
	}), total, nil
}

func (r *RoleService) GetRoleDetail(ctx context.Context, id uint) (*response.RoleDetailResponse, error) {
	role, err := r.roleRepository.GetById(ctx, id, true, true)
	if err != nil {
		return nil, err
	}
	return response.ToRoleDetailResponse(role), nil
}

func (r *RoleService) UpdateRole(ctx context.Context, operator uint, id uint, name string, userIds, permissionIds []uint) error {
	// 对角色自身加锁
	roleLock := utils.NewLock(constant.RoleIdPrefix, id)
	if err := utils.Lock(ctx, roleLock); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := utils.Unlock(lock); e != nil {
			global.Logger.Errorf("unlock fail %s", e.Error())
		}
	}(roleLock)
	// 对 name & permissions 加锁
	locks, err := lockRoleField(ctx, name, permissionIds)
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
	return r.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existRole, temp := r.roleRepository.GetByName(ctx, name)
		if temp != nil && !errors.Is(temp, response.RoleNotExist) {
			return temp
		}
		if existRole != nil && existRole.ID != id {
			return response.RoleCreateDuplicateName
		}
		// 更新角色
		err = r.roleRepository.Update(ctx, &models.Role{
			Name:      name,
			UpdaterId: operator,
			BasicModel: database.BasicModel{
				ID: id,
			},
		})
		if err != nil {
			return err
		}
		// 查询有效的用户
		var users []*models.User
		if len(userIds) > 0 {
			users, err = r.userRepository.GetByIds(ctx, userIds, false, false, false)
			if err != nil {
				return err
			}
		}
		// 更新关联关系
		err = r.roleRepository.AssociateUsers(ctx, id, users)
		if err != nil {
			return err
		}
		// 查询有效的权限
		var permissions []*models.Permission
		if len(permissionIds) > 0 {
			permissions, err = r.permissionRepository.GetByIds(ctx, permissionIds, false)
			if err != nil {
				return err
			}
		}
		return r.roleRepository.AssociatePermissions(ctx, id, permissions)
	})
}

func (r *RoleService) DeleteRole(ctx context.Context, id, operator uint) error {
	// 对角色自身加锁
	roleLock := utils.NewLock(constant.RoleIdPrefix, id)
	if err := utils.Lock(ctx, roleLock); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := utils.Unlock(lock); e != nil {
			global.Logger.Errorf("unlock fail %s", e.Error())
		}
	}(roleLock)
	permissionsCount := r.roleRepository.GetPermissionsCount(ctx, id)
	if permissionsCount > 0 {
		return response.RoleExistPermissionRef
	}
	usersCount := r.roleRepository.GetUsersCount(ctx, id)
	if usersCount > 0 {
		return response.RoleExistUserRef
	}
	return r.roleRepository.DeleteById(ctx, id, operator)
}
