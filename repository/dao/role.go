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

func (r *RoleDAO) Insert(ctx context.Context, role *models.Role) error {
	err := r.Datasource(ctx).Create(role).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.RoleCreateDuplicateName
	}
	return err
}

func (r *RoleDAO) GetRolesByIds(ctx context.Context, ids []uint) ([]*models.Role, error) {
	var roles []*models.Role
	err := r.Datasource(ctx).Where("id in ?", ids).Find(&roles).Error
	return roles, err
}

func (r *RoleDAO) GetRoleList(ctx context.Context, name string, limit, offset int) ([]*models.Role, int64, error) {
	var roles []*models.Role
	var total int64
	query := r.Datasource(ctx)
	if name != "" {
		query = query.Where("name like ?", database.FuzzKeyword(name))
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(limit).Offset(offset).Find(&roles).Error
	return roles, total, err
}
