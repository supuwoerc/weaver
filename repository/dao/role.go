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

var roleDAO *RoleDAO

type Role struct {
	gorm.Model
	Name string `gorm:"unique;not null;;comment:角色名"`
}

func NewRoleDAO(ctx *gin.Context) *RoleDAO {
	if roleDAO == nil {
		roleDAO = &RoleDAO{BasicDAO: NewBasicDao(ctx)}
	}
	return roleDAO
}

func (r RoleDAO) Insert(ctx context.Context, role Role) error {
	err := r.db.WithContext(ctx).Create(&role).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return constant.GetError(r.ctx, response.ROLE_CREATE_DUPLICATE_NAME)
	}
	return err
}
