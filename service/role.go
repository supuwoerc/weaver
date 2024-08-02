package service

import (
	"gin-web/repository"
	"gin-web/repository/dao"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type RoleService struct {
	*BasicService
	repository *repository.RoleRepository
}

func NewRoleService(ctx *gin.Context) *RoleService {
	return &RoleService{
		BasicService: NewBasicService(ctx),
		repository:   repository.NewRoleRepository(ctx),
	}
}

func (r *RoleService) CreateRole(name string) error {
	return r.repository.Create(r.ctx.Request.Context(), name)
}

func (r *RoleService) FilterValidRoles(roleIds []uint) ([]uint, error) {
	roles, err := r.repository.GetRolesByIds(r.ctx.Request.Context(), roleIds)
	if err != nil {
		return []uint{}, err
	}
	validIds := lo.Map[dao.Role, uint](roles, func(item dao.Role, _ int) uint {
		return item.ID
	})
	result := lo.Filter(roleIds, func(item uint, _ int) bool {
		return lo.Contains(validIds, item)
	})
	return result, nil
}
