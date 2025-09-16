package dao

import (
	"context"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/database"
)

type RolePermissionDAO struct {
	*BasicDAO
}

func NewRolePermissionDAO(basicDAO *BasicDAO) *RolePermissionDAO {
	return &RolePermissionDAO{
		BasicDAO: basicDAO,
	}
}

func (r *RolePermissionDAO) GetRolesByPermissionID(ctx context.Context, permissionID uint, keyword string, limit int, offset int) ([]*models.Role, error) {
	var result []*models.Role
	query := r.Datasource(ctx).Model(&models.RolePermission{}).
		Select("r.*").
		Table("sys_role_permission as rp").
		Joins("inner join sys_role r on rp.role_id = r.id")
	if keyword != "" {
		query = query.Where("r.name like ?", database.FuzzKeyword(keyword))
	}
	err := query.
		Where("rp.permission_id = ?", permissionID).
		Order("r.updated_at desc").
		Limit(limit).
		Offset(offset).
		Find(&result).
		Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
