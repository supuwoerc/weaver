package dao

import (
	"context"
	"errors"
	"time"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/database"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/go-sql-driver/mysql"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type PermissionDAO struct {
	*BasicDAO
}

func NewPermissionDAO(basicDAO *BasicDAO) *PermissionDAO {
	return &PermissionDAO{
		BasicDAO: basicDAO,
	}
}

func (r *PermissionDAO) Create(ctx context.Context, permission *models.Permission) error {
	err := r.Datasource(ctx).Create(permission).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.PermissionCreateDuplicate
	}
	return err
}

func (r *PermissionDAO) GetByIds(ctx context.Context, ids []uint, preload ...string) ([]*models.Permission, error) {
	var result []*models.Permission
	query := r.Datasource(ctx).Model(&models.Permission{})
	if len(preload) > 0 {
		query = query.Scopes(lo.Map(preload, func(item string, index int) func(d *gorm.DB) *gorm.DB {
			return r.Preload(item)
		})...)
	}
	err := query.Where("id in (?)", ids).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PermissionDAO) GetById(ctx context.Context, id uint, preload ...string) (*models.Permission, error) {
	var result models.Permission
	query := r.Datasource(ctx).Model(&models.Permission{})
	if len(preload) > 0 {
		query = query.Scopes(lo.Map(preload, func(item string, index int) func(d *gorm.DB) *gorm.DB {
			return r.Preload(item)
		})...)
	}
	err := query.Where("id = ?", id).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.PermissionNotExist
		}
		return nil, err
	}
	return &result, nil
}

func (r *PermissionDAO) GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64
	query := r.Datasource(ctx).Model(&models.Permission{}).Order("updated_at desc,id desc")
	if keyword != "" {
		keyword = database.FuzzKeyword(keyword)
		query = query.Where("name like ? or resource like ?", keyword, keyword)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Preload("Creator").Preload("Updater").Limit(limit).Offset(offset).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}
	return permissions, total, nil
}

func (r *PermissionDAO) DeleteById(ctx context.Context, id, updater uint) error {
	return r.Datasource(ctx).Model(&models.Permission{}).Where("id = ?", id).
		Select("updater_id", "deleted_at").
		Updates(map[string]any{
			"updater_id": updater,
			"deleted_at": time.Now().UnixMilli(),
		}).Error
}

func (r *PermissionDAO) GetRolesCount(ctx context.Context, id uint) int64 {
	return r.Datasource(ctx).Model(&models.Permission{
		BasicModel: database.BasicModel{ID: id},
	}).Association("Roles").Count()
}

func (r *PermissionDAO) Update(ctx context.Context, permission *models.Permission) error {
	err := r.Datasource(ctx).Select("*").Omit("created_at", "roles", "creator_id").Updates(permission).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.PermissionCreateDuplicate
	}
	return err
}

func (r *PermissionDAO) AssociateRoles(ctx context.Context, id uint, roles []*models.Role) error {
	return r.Datasource(ctx).Model(&models.Permission{BasicModel: database.BasicModel{ID: id}}).
		Association("Roles").Replace(roles)
}

func (r *PermissionDAO) GetByNameOrResource(ctx context.Context, name, resource string) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.Datasource(ctx).Model(&models.Permission{}).Where("name = ? or resource = ?", name, resource).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *PermissionDAO) CheckUserPermission(ctx context.Context, uid uint, resource string, permissionType constant.PermissionType) (bool, error) {
	var count int64
	err := r.Datasource(ctx).Model(&models.Permission{}).
		Table("sys_permission as permission").
		Joins("inner join sys_role_permission as role_permission on role_permission.permission_id = permission.id").
		Joins("inner join sys_user_role as user_role on user_role.role_id = role_permission.role_id").
		Where("user_role.user_id = ?", uid).
		Where("permission.resource = ?", resource).
		Where("permission.type = ?", permissionType).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserPermissions 获取用户所有权限
func (r *PermissionDAO) GetUserPermissions(ctx context.Context, userId uint) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.Datasource(ctx).Model(&models.Permission{}).
		Table("sys_permission as permission").
		Joins("inner join sys_role_permission role_permission on permission.id = role_permission.permission_id").
		Joins("inner join sys_user_role user_role on role_permission.role_id = user_role.role_id").
		Where("user_role.user_id = ?", userId).
		Group("permission.id").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetUserPermissionsByType 根据类型获取用户权限
func (r *PermissionDAO) GetUserPermissionsByType(ctx context.Context, userId uint, limit int, offset int,
	permissionType ...constant.PermissionType) ([]*models.Permission, error) {
	var permissions []*models.Permission

	query := r.Datasource(ctx).
		Model(&models.Permission{}).
		Distinct("sys_permission.*").
		Joins("INNER JOIN sys_role_permission ON sys_permission.id = sys_role_permission.permission_id").
		Joins("INNER JOIN sys_user_role ON sys_role_permission.role_id = sys_user_role.role_id").
		Where("sys_user_role.user_id = ?", userId).
		Where("sys_permission.type IN (?)", permissionType).
		Order("sys_permission.id DESC").
		Limit(limit).
		Offset(offset)

	err := query.Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
