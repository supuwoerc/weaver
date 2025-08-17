package permission

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/database"
	"github.com/supuwoerc/weaver/pkg/request"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"
	"github.com/supuwoerc/weaver/service"
)

type DAO interface {
	Create(ctx context.Context, permission *models.Permission) error
	GetByIds(ctx context.Context, ids []uint, preload ...string) ([]*models.Permission, error)
	GetById(ctx context.Context, id uint, preload ...string) (*models.Permission, error)
	GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Permission, int64, error)
	DeleteById(ctx context.Context, id, updater uint) error
	GetRolesCount(ctx context.Context, id uint) int64
	Update(ctx context.Context, permission *models.Permission) error
	AssociateRoles(ctx context.Context, id uint, roles []*models.Role) error
	GetByNameOrResource(ctx context.Context, name, resource string) ([]*models.Permission, error)
	CheckUserPermission(ctx context.Context, uid uint, resource string, permissionType constant.PermissionType) (bool, error)
	GetUserPermissions(ctx context.Context, userId uint) ([]*models.Permission, error)
	GetUserPermissionsByType(ctx context.Context, userId uint, permissionType ...constant.PermissionType) ([]*models.Permission, error)
}

type RoleDAO interface {
	GetByIds(ctx context.Context, ids []uint, preload ...string) ([]*models.Role, error)
}

type Service struct {
	*service.BasicService
	permissionDAO DAO
	roleDAO       RoleDAO
}

func NewPermissionService(basic *service.BasicService, permissionDAO DAO, roleDAO RoleDAO) *Service {
	return &Service{
		BasicService:  basic,
		permissionDAO: permissionDAO,
		roleDAO:       roleDAO,
	}
}

func (s *Service) lockPermissionField(ctx context.Context, name, resource string, roleIds []uint) ([]*utils.RedisLock, error) {
	locks := make([]*utils.RedisLock, 0, len(roleIds)+2)
	// 权限名称锁
	permissionNameLock := s.Locksmith.NewLock(constant.PermissionNamePrefix, name)
	if err := permissionNameLock.Lock(ctx, true); err != nil {
		return locks, err
	}
	locks = append(locks, permissionNameLock)
	// 权限资源锁
	permissionResourceLock := s.Locksmith.NewLock(constant.PermissionResourcePrefix, resource)
	if err := permissionResourceLock.Lock(ctx, true); err != nil {
		return locks, err
	}
	locks = append(locks, permissionResourceLock)
	// 角色锁
	for _, roleId := range roleIds {
		roleIdLock := s.Locksmith.NewLock(constant.RoleIdPrefix, strconv.Itoa(int(roleId)))
		if err := roleIdLock.Lock(ctx, true); err != nil {
			return locks, err
		}
		locks = append(locks, roleIdLock)
	}
	return locks, nil
}

func (s *Service) CreatePermission(ctx context.Context, operator uint, params *request.CreatePermissionRequest) error {
	locks, err := s.lockPermissionField(ctx, params.Name, params.Resource, params.Roles)
	defer func() {
		for _, l := range locks {
			if e := l.Unlock(); e != nil {
				s.Logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	return s.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existPermissions, temp := s.permissionDAO.GetByNameOrResource(ctx, params.Name, params.Resource)
		if temp != nil {
			return temp
		}
		if len(existPermissions) > 0 {
			return response.PermissionCreateDuplicate
		}
		// 查询有效的角色
		var roles []*models.Role
		if len(params.Roles) > 0 {
			roles, err = s.roleDAO.GetByIds(ctx, params.Roles)
			if err != nil {
				return err
			}
		}
		// 创建权限 & 建立关联关系
		return s.permissionDAO.Create(ctx, &models.Permission{
			Name:      params.Name,
			Resource:  params.Resource,
			Type:      params.Type,
			Roles:     roles,
			CreatorId: operator,
			UpdaterId: operator,
		})
	})
}

func (s *Service) GetPermissionList(ctx context.Context, keyword string, limit, offset int) ([]*response.PermissionListRowResponse, int64, error) {
	list, total, err := s.permissionDAO.GetList(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return lo.Map(list, func(item *models.Permission, _ int) *response.PermissionListRowResponse {
		return response.ToPermissionListRowResponse(item)
	}), total, nil
}

func (s *Service) GetPermissionDetail(ctx context.Context, id uint) (*response.PermissionDetailResponse, error) {
	permission, err := s.permissionDAO.GetById(ctx, id, "Roles")
	if err != nil {
		return nil, err
	}
	return response.ToPermissionDetailResponse(permission), nil
}

func (s *Service) UpdatePermission(ctx context.Context, operator uint, params *request.UpdatePermissionRequest) error {
	// 对权限自身加锁
	permissionLock := s.Locksmith.NewLock(constant.PermissionIdPrefix, strconv.Itoa(int(params.ID)))
	if err := permissionLock.Lock(ctx, true); err != nil {
		return err
	}
	_, err := s.permissionDAO.GetById(ctx, params.ID)
	if err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := lock.Unlock(); e != nil {
			s.Logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
		}
	}(permissionLock)
	// 对 name & resource & roleIds 加锁
	locks, err := s.lockPermissionField(ctx, params.Name, params.Resource, params.Roles)
	defer func() {
		for _, l := range locks {
			if e := l.Unlock(); e != nil {
				s.Logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	return s.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existPermissions, temp := s.permissionDAO.GetByNameOrResource(ctx, params.Name, params.Resource)
		if temp != nil {
			return temp
		}
		count := len(existPermissions)
		if count > 1 || (count == 1 && existPermissions[0].ID != params.ID) {
			return response.PermissionCreateDuplicate
		}
		// 更新权限
		err = s.permissionDAO.Update(ctx, &models.Permission{
			Name:      params.Name,
			Resource:  params.Resource,
			Type:      params.Type,
			UpdaterId: operator,
			BasicModel: database.BasicModel{
				ID: params.ID,
			},
		})
		if err != nil {
			return err
		}
		// 查询有效的角色
		var roles []*models.Role
		if len(params.Roles) > 0 {
			roles, err = s.roleDAO.GetByIds(ctx, params.Roles)
			if err != nil {
				return err
			}
		}
		// 更新关联关系
		return s.permissionDAO.AssociateRoles(ctx, params.ID, roles)
	})
}

func (s *Service) DeletePermission(ctx context.Context, id, operator uint) error {
	// 对权限自身加锁
	permissionLock := s.Locksmith.NewLock(constant.PermissionIdPrefix, strconv.Itoa(int(id)))
	if err := permissionLock.Lock(ctx, true); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := lock.Unlock(); e != nil {
			s.Logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
		}
	}(permissionLock)
	count := s.permissionDAO.GetRolesCount(ctx, id)
	if count > 0 {
		return response.PermissionExistRoleRef
	}
	return s.permissionDAO.DeleteById(ctx, id, operator)
}
