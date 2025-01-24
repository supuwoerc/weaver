package dao

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/database"
	"gin-web/pkg/response"
	"github.com/go-sql-driver/mysql"
	"sync"
)

var (
	roleDAO     *RoleDAO
	roleDAOOnce sync.Once
)

type RoleDAO struct {
	*BasicDAO
}

func NewRoleDAO() *RoleDAO {
	roleDAOOnce.Do(func() {
		roleDAO = &RoleDAO{BasicDAO: NewBasicDao()}
	})
	return roleDAO
}

func (r *RoleDAO) Create(ctx context.Context, role *models.Role) error {
	err := r.Datasource(ctx).Create(role).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.RoleCreateDuplicateName
	}
	return err
}

func (r *RoleDAO) GetByIds(ctx context.Context, ids []uint, needUsers, needPermissions bool) ([]*models.Role, error) {
	var roles []*models.Role
	query := r.Datasource(ctx).Model(&models.Role{})
	if needUsers {
		query = query.Preload("Users")
	}
	if needPermissions {
		query = query.Preload("Permissions")
	}
	err := query.Where("id in (?)", ids).Find(&roles).Error
	return roles, err
}

func (r *RoleDAO) GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Role, int64, error) {
	var roles []*models.Role
	var total int64
	query := r.Datasource(ctx).Model(&models.Role{})
	if keyword != "" {
		query = query.Where("name like ?", database.FuzzKeyword(keyword))
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(limit).Offset(offset).Find(&roles).Error
	return roles, total, err
}

func (r *RoleDAO) GetIsExistByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.Datasource(ctx).Model(&models.Role{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
