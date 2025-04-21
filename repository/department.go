package repository

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/constant"
)

type DepartmentDAO interface {
	Create(ctx context.Context, dept *models.Department) error
	GetByName(ctx context.Context, name string) (*models.Department, error)
	GetById(ctx context.Context, id uint) (*models.Department, error)
	GetAll(ctx context.Context) ([]*models.Department, error)
	GetAllUserDepartment(ctx context.Context) ([]*models.UserDepartment, error)
	GetAllDepartmentLeader(ctx context.Context) ([]*models.DepartmentLeader, error)
}

type DepartmentCache interface {
	CacheDepartment(ctx context.Context, key constant.CacheKey, depts []*models.Department) error
	GetDepartmentCache(ctx context.Context, key constant.CacheKey) ([]*models.Department, error)
	RemoveDepartmentCache(ctx context.Context, keys ...constant.CacheKey) error
}

type DepartmentRepository struct {
	DepartmentDAO
	DepartmentCache
}

func NewDepartmentRepository(dao DepartmentDAO, cache DepartmentCache) *DepartmentRepository {
	return &DepartmentRepository{
		DepartmentDAO:   dao,
		DepartmentCache: cache,
	}
}
