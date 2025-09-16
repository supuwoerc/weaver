package dao

import (
	"context"

	"github.com/supuwoerc/weaver/models"
)

type RolePermissionDAO struct {
	*BasicDAO
}

func NewRolePermissionDAO(basicDAO *BasicDAO) *RolePermissionDAO {
	return &RolePermissionDAO{
		BasicDAO: basicDAO,
	}
}

func (r *RolePermissionDAO) GetRolesByPermissionID(ctx context.Context, permissionID uint, limit int, offset int) ([]*models.Role, error) {
	var result []*models.Role
	query := r.Datasource(ctx).Model(&models.RolePermission{}).
		Table("sys_role_permission as rp")
	err := query.Where("permission_id = ?", permissionID).
		Joins("inner join sys_role r on rp.role_id = r.id").
		Find(&result).Limit(limit).Offset(offset).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
