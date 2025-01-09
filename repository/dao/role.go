package dao

import (
	"context"
	"errors"
	"gin-web/pkg/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	roleDAO     *RoleDAO
	roleDAOOnce sync.Once
)

type RoleDAO struct {
	*BasicDAO
}

type Role struct {
	gorm.Model
	Name        string        `gorm:"unique;not null;comment:角色名"`
	Users       []*User       `gorm:"many2many:user_role;"`
	Permissions []*Permission `gorm:"many2many:role_permission;"`
}

func NewRoleDAO() *RoleDAO {
	roleDAOOnce.Do(func() {
		roleDAO = &RoleDAO{BasicDAO: NewBasicDao()}
	})
	return roleDAO
}

func (r *RoleDAO) Insert(ctx context.Context, role *Role) error {
	err := r.Datasource(ctx).Create(role).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.RoleCreateDuplicateName
	}
	return err
}

func (r *RoleDAO) GetRolesByIds(ctx context.Context, ids []uint) ([]*Role, error) {
	var roles []*Role
	err := r.Datasource(ctx).Where("id in ?", ids).Find(&roles).Error
	return roles, err
}
