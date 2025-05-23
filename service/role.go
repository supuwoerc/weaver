package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/database"
	"github.com/supuwoerc/weaver/pkg/request"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"

	"github.com/samber/lo"
)

type RoleDAO interface {
	Create(ctx context.Context, role *models.Role) error
	GetByIds(ctx context.Context, ids []uint, preload ...string) ([]*models.Role, error)
	GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Role, int64, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	GetById(ctx context.Context, id uint, preload ...string) (*models.Role, error)
	Update(ctx context.Context, role *models.Role) error
	AssociateUsers(ctx context.Context, id uint, users []*models.User) error
	AssociatePermissions(ctx context.Context, id uint, permissions []*models.Permission) error
	DeleteById(ctx context.Context, id, updater uint) error
	GetUsersCount(ctx context.Context, id uint) int64
	GetPermissionsCount(ctx context.Context, id uint) int64
}

type RoleService struct {
	*BasicService
	roleDAO       RoleDAO
	userDAO       UserDAO
	permissionDAO PermissionDAO
}

func NewRoleService(basic *BasicService, roleDAO RoleDAO, userDAO UserDAO, permissionDAO PermissionDAO) *RoleService {
	return &RoleService{
		BasicService:  basic,
		roleDAO:       roleDAO,
		userDAO:       userDAO,
		permissionDAO: permissionDAO,
	}
}

func (r *RoleService) lockRoleField(ctx context.Context, name string, permissionIds []uint) ([]*utils.RedisLock, error) {
	locks := make([]*utils.RedisLock, 0, len(permissionIds)+1)
	// 角色名称锁
	roleNameLock := r.locksmith.NewLock(constant.RoleNamePrefix, name)
	if err := roleNameLock.Lock(ctx, true); err != nil {
		return locks, err
	}
	locks = append(locks, roleNameLock)
	// 角色权限锁
	for _, permissionId := range permissionIds {
		roleIdLock := r.locksmith.NewLock(constant.PermissionIdPrefix, strconv.Itoa(int(permissionId)))
		if err := roleIdLock.Lock(ctx, true); err != nil {
			return locks, err
		}
		locks = append(locks, roleIdLock)
	}
	return locks, nil
}

func (r *RoleService) CreateRole(ctx context.Context, operator uint, params *request.CreateRoleRequest) error {
	locks, err := r.lockRoleField(ctx, params.Name, params.Permissions)
	defer func() {
		for _, l := range locks {
			if e := l.Unlock(); e != nil {
				r.logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	// TODO:记录信息到用户时间线
	return r.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existRole, temp := r.roleDAO.GetByName(ctx, params.Name)
		if temp != nil && !errors.Is(temp, response.RoleNotExist) {
			return temp
		}
		if existRole != nil {
			return response.RoleCreateDuplicateName
		}
		// 查询有效的用户
		var users []*models.User
		if len(params.Users) > 0 {
			users, err = r.userDAO.GetByIds(ctx, params.Users)
			if err != nil {
				return err
			}
		}
		// 查询有效的权限
		var permissions []*models.Permission
		if len(params.Permissions) > 0 {
			permissions, err = r.permissionDAO.GetByIds(ctx, params.Permissions)
			if err != nil {
				return err
			}
		}
		// 创建角色 & 建立关联关系
		return r.roleDAO.Create(ctx, &models.Role{
			Name:        params.Name,
			Users:       users,
			Permissions: permissions,
			CreatorId:   operator,
			UpdaterId:   operator,
		})
	})
}

func (r *RoleService) GetRoleList(ctx context.Context, keyword string, limit, offset int) ([]*response.RoleListRowResponse, int64, error) {
	list, total, err := r.roleDAO.GetList(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return lo.Map(list, func(item *models.Role, _ int) *response.RoleListRowResponse {
		return response.ToRoleListRowResponse(item)
	}), total, nil
}

func (r *RoleService) GetRoleDetail(ctx context.Context, id uint) (*response.RoleDetailResponse, error) {
	role, err := r.roleDAO.GetById(ctx, id, "Users", "Permissions")
	if err != nil {
		return nil, err
	}
	return response.ToRoleDetailResponse(role), nil
}

func (r *RoleService) UpdateRole(ctx context.Context, operator uint, params *request.UpdateRoleRequest) error {
	// 对角色自身加锁
	roleLock := r.locksmith.NewLock(constant.RoleIdPrefix, strconv.Itoa(int(params.ID)))
	if err := roleLock.Lock(ctx, true); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := lock.Unlock(); e != nil {
			r.logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
		}
	}(roleLock)
	_, err := r.roleDAO.GetById(ctx, params.ID)
	if err != nil {
		return err
	}
	// 对 name & permissions 加锁
	locks, err := r.lockRoleField(ctx, params.Name, params.Permissions)
	defer func() {
		for _, l := range locks {
			if e := l.Unlock(); e != nil {
				r.logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	return r.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existRole, temp := r.roleDAO.GetByName(ctx, params.Name)
		if temp != nil && !errors.Is(temp, response.RoleNotExist) {
			return temp
		}
		if existRole != nil && existRole.ID != params.ID {
			return response.RoleCreateDuplicateName
		}
		// 更新角色
		err = r.roleDAO.Update(ctx, &models.Role{
			Name:      params.Name,
			UpdaterId: operator,
			BasicModel: database.BasicModel{
				ID: params.ID,
			},
		})
		if err != nil {
			return err
		}
		// 查询有效的用户
		var users []*models.User
		if len(params.Users) > 0 {
			users, err = r.userDAO.GetByIds(ctx, params.Users)
			if err != nil {
				return err
			}
		}
		// 更新关联关系
		err = r.roleDAO.AssociateUsers(ctx, params.ID, users)
		if err != nil {
			return err
		}
		// 查询有效的权限
		var permissions []*models.Permission
		if len(params.Permissions) > 0 {
			permissions, err = r.permissionDAO.GetByIds(ctx, params.Permissions)
			if err != nil {
				return err
			}
		}
		return r.roleDAO.AssociatePermissions(ctx, params.ID, permissions)
	})
}

func (r *RoleService) DeleteRole(ctx context.Context, id, operator uint) error {
	// 对角色自身加锁
	roleLock := r.locksmith.NewLock(constant.RoleIdPrefix, strconv.Itoa(int(id)))
	if err := roleLock.Lock(ctx, true); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := lock.Unlock(); e != nil {
			r.logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
		}
	}(roleLock)
	permissionsCount := r.roleDAO.GetPermissionsCount(ctx, id)
	if permissionsCount > 0 {
		return response.RoleExistPermissionRef
	}
	usersCount := r.roleDAO.GetUsersCount(ctx, id)
	if usersCount > 0 {
		return response.RoleExistUserRef
	}
	return r.roleDAO.DeleteById(ctx, id, operator)
}
