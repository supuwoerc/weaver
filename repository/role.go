package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type RoleRepository struct {
	dao *dao.RoleDAO
}

func NewRoleRepository(ctx *gin.Context) *RoleRepository {
	return &RoleRepository{
		dao: dao.NewRoleDAO(ctx),
	}
}

func toModelRole(role *dao.Role) *models.Role {
	return &models.Role{
		ID:   role.ID,
		Name: role.Name,
	}
}

func toModelRoles(roles []*dao.Role) []*models.Role {
	return lo.Map(roles, func(item *dao.Role, index int) *models.Role {
		return toModelRole(item)
	})
}

func (r *RoleRepository) Create(ctx context.Context, name string) error {
	return r.dao.Insert(ctx, &dao.Role{
		Name: name,
	})
}

func (r *RoleRepository) GetRolesByIds(ctx context.Context, ids []uint) ([]*dao.Role, error) {
	return r.dao.GetRolesByIds(ctx, ids)
}
