package service

import (
	"context"
	"gin-web/repository"
	"github.com/gin-gonic/gin"
)

type RoleService struct {
	*BasicService
	repository *repository.RoleRepository
}

var roleService *RoleService

func NewRoleService(ctx *gin.Context) *RoleService {
	if roleService == nil {
		roleService = &RoleService{
			BasicService: NewBasicService(ctx),
			repository:   repository.NewRoleRepository(ctx),
		}
	}
	return roleService
}

func (r *RoleService) CreateRole(context context.Context, name string) error {
	return r.repository.Create(context, name)
}
