package dao

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/response"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DepartmentDAO struct {
	*BasicDAO
}

func NewDepartmentDAO(basicDAO *BasicDAO) *DepartmentDAO {
	return &DepartmentDAO{
		BasicDAO: basicDAO,
	}
}

func (r *DepartmentDAO) Create(ctx context.Context, dept *models.Department) error {
	err := r.Datasource(ctx).Create(dept).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.RoleCreateDuplicateName
	}
	return err
}

func (r *DepartmentDAO) GetByName(ctx context.Context, name string) (*models.Department, error) {
	var dept models.Department
	err := r.Datasource(ctx).Model(&models.Department{}).Where("name = ?", name).First(&dept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.DeptNotExist
		}
		return nil, err
	}
	return &dept, nil
}

func (r *DepartmentDAO) GetById(ctx context.Context, id uint) (*models.Department, error) {
	var dept models.Department
	err := r.Datasource(ctx).Model(&models.Department{}).Where("id = ?", id).First(&dept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.DeptNotExist
		}
		return nil, err
	}
	return &dept, nil
}

func (r *DepartmentDAO) GetAll(ctx context.Context) ([]*models.Department, error) {
	var depts []*models.Department
	err := r.Datasource(ctx).Model(&models.Department{}).Find(&depts).Error
	if err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *DepartmentDAO) GetAllUserDepartment(ctx context.Context) ([]*models.UserDepartment, error) {
	var res []*models.UserDepartment
	err := r.Datasource(ctx).Model(&models.UserDepartment{}).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *DepartmentDAO) GetAllDepartmentLeader(ctx context.Context) ([]*models.DepartmentLeader, error) {
	var res []*models.DepartmentLeader
	err := r.Datasource(ctx).Model(&models.DepartmentLeader{}).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
