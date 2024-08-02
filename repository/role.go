package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
	"github.com/gin-gonic/gin"
)

type RoleRepository struct {
	dao *dao.RoleDAO
}

func NewRoleRepository(ctx *gin.Context) *RoleRepository {
	return &RoleRepository{
		dao: dao.NewRoleDAO(ctx),
	}
}

func toModelRole(r dao.Role) models.Role {
	return models.Role{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (r *RoleRepository) Create(ctx context.Context, name string) error {
	return r.dao.Insert(ctx, dao.Role{
		Name: name,
	})
}

func (r *RoleRepository) GetRolesByIds(ctx context.Context, ids []uint) ([]dao.Role, error) {
	return r.dao.GetRolesByIds(ctx, ids)
}
