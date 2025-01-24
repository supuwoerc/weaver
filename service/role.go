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

func (r *RoleService) CreateRole(ctx context.Context, name string, userIds, permissionIds []uint) error {
	roleNameLock := utils.NewLock(constant.RoleNamePrefix, name)
	if err := utils.Lock(ctx, roleNameLock); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := utils.Unlock(lock); e != nil {
			global.Logger.Errorf("unlock fail %s", e.Error())
		}
	}(roleNameLock)
	// TODO:记录信息到用户时间线
	return r.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		exist, err := r.roleRepository.GetIsExistByName(ctx, name)
		if err != nil {
			return err
		}
		if exist {
			return response.RoleCreateDuplicateName
		}
		// 查询有效的用户
		users, err := r.userRepository.GetByIds(ctx, userIds, false, false)
		if err != nil {
			return err
		}
		// 查询有效的权限
		permissions, err := r.permissionRepository.GetByIds(ctx, permissionIds, false)
		if err != nil {
			return err
		}
		// 创建角色 & 建立关联关系
		return r.roleRepository.Create(ctx, name, users, permissions)
	})
}

func (r *RoleService) GetRoleList(ctx context.Context, keyword string, limit, offset int) ([]*models.Role, int64, error) {
	return r.roleRepository.GetList(ctx, keyword, limit, offset)
}
