package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
	"sync"
)

var (
	departmentRepository     *DepartmentRepository
	departmentRepositoryOnce sync.Once
)

type DepartmentRepository struct {
	dao *dao.DepartmentDAO
}

func NewDepartmentRepository() *DepartmentRepository {
	departmentRepositoryOnce.Do(func() {
		departmentRepository = &DepartmentRepository{
			dao: dao.NewDepartmentDAO(),
		}
	})
	return departmentRepository
}

func (r *DepartmentRepository) Create(ctx context.Context, dept *models.Department) error {
	return r.dao.Create(ctx, dept)
}

func (r *DepartmentRepository) GetByName(ctx context.Context, name string) (*models.Department, error) {
	return r.dao.GetByName(ctx, name)
}
