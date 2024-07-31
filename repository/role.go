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

var roleRepository *RoleRepository

func NewRoleRepository(ctx *gin.Context) *RoleRepository {
	if roleRepository == nil {
		roleRepository = &RoleRepository{
			dao: dao.NewRoleDAO(ctx),
		}
	}
	return roleRepository
}

func toModelRole(r dao.Role) models.Role {
	return models.Role{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (r RoleRepository) Create(ctx context.Context, name string) error {
	return r.dao.Insert(ctx, dao.Role{
		Name: name,
	})
}
