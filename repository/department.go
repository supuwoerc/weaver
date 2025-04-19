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
}

type DepartmentRepository struct {
	dao   DepartmentDAO
	cache DepartmentCache
}

func NewDepartmentRepository(dao DepartmentDAO, cache DepartmentCache) *DepartmentRepository {
	return &DepartmentRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *DepartmentRepository) Create(ctx context.Context, dept *models.Department) error {
	return r.dao.Create(ctx, dept)
}

func (r *DepartmentRepository) GetByName(ctx context.Context, name string) (*models.Department, error) {
	return r.dao.GetByName(ctx, name)
}

func (r *DepartmentRepository) GetById(ctx context.Context, id uint) (*models.Department, error) {
	return r.dao.GetById(ctx, id)
}

func (r *DepartmentRepository) GetAll(ctx context.Context) ([]*models.Department, error) {
	return r.dao.GetAll(ctx)
}

func (r *DepartmentRepository) GetAllUserDepartment(ctx context.Context) ([]*models.UserDepartment, error) {
	return r.dao.GetAllUserDepartment(ctx)
}

func (r *DepartmentRepository) GetAllDepartmentLeader(ctx context.Context) ([]*models.DepartmentLeader, error) {
	return r.dao.GetAllDepartmentLeader(ctx)
}

func (r *DepartmentRepository) CacheDepartment(ctx context.Context, key constant.CacheKey, depts []*models.Department) error {
	return r.cache.CacheDepartment(ctx, key, depts)
}

func (r *DepartmentRepository) GetDepartmentCache(ctx context.Context, key constant.CacheKey) ([]*models.Department, error) {
	return r.cache.GetDepartmentCache(ctx, key)
}
