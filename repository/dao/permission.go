package dao

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/database"
	"gin-web/pkg/response"
	"github.com/go-sql-driver/mysql"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"time"
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
