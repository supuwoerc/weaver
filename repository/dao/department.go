package dao

import (
	"context"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/response"

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
	err := r.Datasource(ctx).Omit("Leaders", "Users").Create(dept).Error
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

func (r *DepartmentDAO) GetByID(ctx context.Context, id uint) (*models.Department, error) {
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

// GetAll 查询全部部门数据
func (r *DepartmentDAO) GetAll(ctx context.Context) ([]*models.Department, error) {
	depts, err := queryAll[*models.Department](r.Datasource(ctx).Model(&models.Department{}), r.QueryLimit)
	if err != nil {
		return nil, err
	}
	return depts, nil
}

// GetAllUserDepartment 查询全部部门-人员关联数据
func (r *DepartmentDAO) GetAllUserDepartment(ctx context.Context) ([]*models.UserDepartment, error) {
	res, err := queryAll[*models.UserDepartment](r.Datasource(ctx).Model(&models.UserDepartment{}), r.QueryLimit)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetAllDepartmentLeader 查询全部部门-Leader关联数据
func (r *DepartmentDAO) GetAllDepartmentLeader(ctx context.Context) ([]*models.DepartmentLeader, error) {
	res, err := queryAll[*models.DepartmentLeader](r.Datasource(ctx).Model(&models.DepartmentLeader{}), r.QueryLimit)
	if err != nil {
		return nil, err
	}
	return res, nil
}
