package dao

import (
	"context"
	"errors"
	"gin-web/pkg/constant"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type RoleDAO struct {
	*BasicDAO
}

type Role struct {
	gorm.Model
	Name  string `gorm:"unique;not null;;comment:角色名"`
	Users []User `gorm:"many2many:user_role"`
}

type PureRole struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"unique;not null;;comment:角色名"`
}

func NewRoleDAO(ctx *gin.Context) *RoleDAO {
	return &RoleDAO{BasicDAO: NewBasicDao(ctx)}
}

func (r *RoleDAO) Insert(ctx context.Context, role *Role) error {
	err := r.db.WithContext(ctx).Create(role).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return constant.GetError(r.ctx, response.ROLE_CREATE_DUPLICATE_NAME)
	}
	return err
}

func (r *RoleDAO) GetRolesByIds(ctx context.Context, ids []uint) ([]*Role, error) {
	var roles []*Role
	err := r.db.WithContext(ctx).Where("id in ?", ids).Find(&roles).Error
	return roles, err
}
