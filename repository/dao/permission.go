package dao

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/database"
	"gin-web/pkg/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	permissionDAO     *PermissionDAO
	permissionDAOOnce sync.Once
)

type PermissionDAO struct {
	*BasicDAO
}

func NewPermissionDAO() *PermissionDAO {
	permissionDAOOnce.Do(func() {
		permissionDAO = &PermissionDAO{BasicDAO: NewBasicDao()}
	})
	return permissionDAO
}

func (r *PermissionDAO) Create(ctx context.Context, permission *models.Permission) error {
	err := r.Datasource(ctx).Create(permission).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.PermissionCreateDuplicate
	}
	return err
}

func (r *PermissionDAO) GetByIds(ctx context.Context, ids []uint, needRoles bool) ([]*models.Permission, error) {
	var result []*models.Permission
	query := r.Datasource(ctx).Model(&models.Permission{})
	if needRoles {
		query = query.Preload("Roles")
	}
	err := query.Where("id in (?)", ids).Find(&result).Error
	return result, err
}

func (r *PermissionDAO) GetById(ctx context.Context, id uint, needRoles bool) (*models.Permission, error) {
	var result = models.Permission{
		Roles: make([]*models.Role, 0),
	}
	query := r.Datasource(ctx).Model(&models.Permission{})
	if needRoles {
		query = query.Preload("Roles")
	}
	err := query.Where("id = ?", id).Find(&result).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return &result, response.PermissionCreateDuplicate
	}
	return &result, err
}

func (r *PermissionDAO) GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64
	query := r.Datasource(ctx).Model(&models.Permission{})
	if keyword != "" {
		keyword = database.FuzzKeyword(keyword)
		query = query.Where("name like ? or resource like ?", keyword, keyword)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}
