package service

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/repository"
	"github.com/pkg/errors"
	"sync"
)

type DepartmentService struct {
	*BasicService
	departmentRepository *repository.DepartmentRepository
}

var (
	departmentOnce    sync.Once
	departmentService *DepartmentService
)

func NewDepartmentService() *DepartmentService {
	departmentOnce.Do(func() {
		departmentService = &DepartmentService{
			BasicService:         NewBasicService(),
			departmentRepository: repository.NewDepartmentRepository(),
		}
	})
	return departmentService
}

func lockDepartmentField(ctx context.Context, name string) ([]*utils.RedisLock, error) {
	locks := make([]*utils.RedisLock, 0)
	// 名称锁
	deptNameLock := utils.NewLock(constant.DepartmentNamePrefix, name)
	if err := utils.Lock(ctx, deptNameLock); err != nil {
		return locks, err
	}
	locks = append(locks, deptNameLock)
	return locks, nil
}

func (p *DepartmentService) CreateDepartment(ctx context.Context, operator uint, name string, parentId *uint) error {
	locks, err := lockDepartmentField(ctx, name)
	defer func() {
		for _, l := range locks {
			if e := utils.Unlock(l); e != nil {
				global.Logger.Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	return p.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existDept, temp := p.departmentRepository.GetByName(ctx, name)
		if temp != nil && !errors.Is(temp, response.DeptNotExist) {
			return temp
		}
		if existDept != nil {
			return response.DeptCreateDuplicate
		}
		// 创建部门 & 建立关联关系
		// TODO:完善 Parent & Ancestors & Leaders & Users
		dept := &models.Department{
			Name:      name,
			ParentId:  parentId,
			CreatorId: operator,
			UpdaterId: operator,
		}
		return p.departmentRepository.Create(ctx, dept)
	})
}
