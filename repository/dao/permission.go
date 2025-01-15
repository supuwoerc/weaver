package dao

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/response"
	"github.com/go-sql-driver/mysql"
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

func (r *PermissionDAO) Insert(ctx context.Context, permission *models.Permission) error {
	err := r.Datasource(ctx).Create(permission).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.PermissionCreateDuplicate
	}
	return err
}
