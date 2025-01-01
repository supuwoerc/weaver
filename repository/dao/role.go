package dao

import (
	"context"
	"errors"
	"gin-web/pkg/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type RoleDAO struct {
	*BasicDAO
}

type Role struct {
	gorm.Model
	Name  string  `gorm:"unique;not null;comment:角色名"`
	Users []*User `gorm:"many2many:user_role;"`
}

func NewRoleDAO() *RoleDAO {
	return &RoleDAO{BasicDAO: NewBasicDao()}
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
