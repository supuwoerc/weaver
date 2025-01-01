package service

import (
	"gin-web/models"
	"gin-web/repository"
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
		repository:   repository.NewRoleRepository(),
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
	validIds := lo.Map[*models.Role, uint](roles, func(item *models.Role, _ int) uint {
		return item.ID
	})
	result := lo.Filter(roleIds, func(item uint, _ int) bool {
		return lo.Contains(validIds, item)
	})
	return result, nil
}
