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
	err := r.Datasource(ctx).Model(&models.RolePermission{}).
		Select("r.*").
		Table("sys_role_permission as rp").
		Where("permission_id = ?", permissionID).
		Joins("inner join sys_role r on rp.role_id = r.id").
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
