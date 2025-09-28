package dao

import (
	"context"
	"errors"
	"time"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/database"
	"github.com/supuwoerc/weaver/pkg/response"
	"gorm.io/gorm/clause"

	"github.com/go-sql-driver/mysql"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type RoleDAO struct {
	*BasicDAO
}

func NewRoleDAO(basicDAO *BasicDAO) *RoleDAO {
	return &RoleDAO{
		BasicDAO: basicDAO,
	}
}

func (r *RoleDAO) Create(ctx context.Context, role *models.Role) error {
	err := r.Datasource(ctx).Create(role).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.RoleCreateDuplicateName
	}
	return err
}

func (r *RoleDAO) GetByIds(ctx context.Context, ids []uint, preload ...string) ([]*models.Role, error) {
	var roles []*models.Role
	query := r.Datasource(ctx).Model(&models.Role{})
	if len(preload) > 0 {
		query = query.Scopes(lo.Map(preload, func(item string, index int) func(d *gorm.DB) *gorm.DB {
			return r.Preload(item)
		})...)
	}
	err := query.Where("id in (?)", ids).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleDAO) GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.Role, int64, error) {
	var roles []*models.Role
	var total int64
	query := r.Datasource(ctx).Model(&models.Role{}).
		Preload("Children").
		Order("updated_at desc,id desc")
	if keyword != "" {
		query = query.Where("name like ?", database.FuzzKeyword(keyword))
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Preload("Creator").Preload("Updater").Limit(limit).Offset(offset).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (r *RoleDAO) GetByName(ctx context.Context, name string) (*models.Role, error) {
	var role *models.Role
	err := r.Datasource(ctx).Model(&models.Role{}).Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.RoleNotExist
		}
		return nil, err
	}
	return role, nil
}

func (r *RoleDAO) GetByID(ctx context.Context, id uint, preload ...string) (*models.Role, error) {
	var result models.Role
	query := r.Datasource(ctx).Model(&models.Role{})
	if len(preload) > 0 {
		query = query.Scopes(lo.Map(preload, func(item string, index int) func(d *gorm.DB) *gorm.DB {
			return r.Preload(item)
		})...)
	}
	err := query.Where("id = ?", id).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.RoleNotExist
		}
		return nil, err
	}
	return &result, nil
}

func (r *RoleDAO) Update(ctx context.Context, role *models.Role) error {
	err := r.Datasource(ctx).Select("*").Omit("created_at", "users", "permissions", "creator_id").
		Updates(role).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.RoleCreateDuplicateName
	}
	return err
}

func (r *RoleDAO) AssociateUsers(ctx context.Context, id uint, users []*models.User) error {
	return r.Datasource(ctx).Model(&models.Role{BasicModel: database.BasicModel{ID: id}}).Association("Users").Replace(users)
}

func (r *RoleDAO) AssociatePermissions(ctx context.Context, id uint, permissions []*models.Permission) error {
	return r.Datasource(ctx).Model(&models.Role{BasicModel: database.BasicModel{ID: id}}).Association("Permissions").Replace(permissions)
}

func (r *RoleDAO) DeleteByID(ctx context.Context, id, updater uint) error {
	return r.Datasource(ctx).Model(&models.Role{}).Where("id = ?", id).
		Select("updater_id", "deleted_at").
		Updates(map[string]any{
			"updater_id": updater,
			"deleted_at": time.Now().UnixMilli(),
		}).Error
}

func (r *RoleDAO) GetUsersCount(ctx context.Context, id uint) int64 {
	return r.Datasource(ctx).Model(&models.Role{
		BasicModel: database.BasicModel{ID: id},
	}).Association("Users").Count()
}

func (r *RoleDAO) GetPermissionsCount(ctx context.Context, id uint) int64 {
	return r.Datasource(ctx).Model(&models.Role{
		BasicModel: database.BasicModel{ID: id},
	}).Association("Permissions").Count()
}

func (r *RoleDAO) GetUserRolesWithoutPosterity(ctx context.Context, uid uint) ([]*models.Role, error) {
	var roles []*models.Role
	err := r.Datasource(ctx).Model(&models.Role{}).
		Table("sys_role as role").
		Joins("inner join sys_user_role as user_role on user_role.role_id = role_permission.role_id").
		Where("user_role.user_id = ?", uid).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleDAO) GetRolesWithPosterity(ctx context.Context, roleIds ...uint) ([]*models.Role, error) {
	if len(roleIds) == 0 {
		return nil, nil
	}
	var roles []*models.Role
	db := r.Datasource(ctx).Model(&models.Role{})
	orConditions := make([]clause.Expression, 0, len(roleIds))
	for _, id := range roleIds {
		orConditions = append(orConditions, gorm.Expr("find_in_set(?,ancestors)", id))
	}
	err := db.Where(clause.Or(orConditions...)).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
