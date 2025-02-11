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
	"github.com/samber/lo"
	"strconv"
	"strings"
	"sync"
)

type DepartmentService struct {
	*BasicService
	departmentRepository *repository.DepartmentRepository
	userRepository       *repository.UserRepository
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
			userRepository:       repository.NewUserRepository(),
		}
	})
	return departmentService
}

func lockDepartmentField(ctx context.Context, name string, parentId *uint) ([]*utils.RedisLock, error) {
	locks := make([]*utils.RedisLock, 0)
	// 名称锁
	deptNameLock := utils.NewLock(constant.DepartmentNamePrefix, name)
	if err := utils.Lock(ctx, deptNameLock); err != nil {
		return locks, err
	}
	locks = append(locks, deptNameLock)
	if parentId != nil {
		// 父部门锁
		parentLock := utils.NewLock(constant.DepartmentIdPrefix, *parentId)
		if err := utils.Lock(ctx, parentLock); err != nil {
			return locks, err
		}
		locks = append(locks, parentLock)
	}
	return locks, nil
}

func (p *DepartmentService) CreateDepartment(ctx context.Context, operator uint, name string, parentId *uint, leaderIds, userIds []uint) error {
	locks, err := lockDepartmentField(ctx, name, parentId)
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
		// 查询父部门
		var parentDept *models.Department
		if parentId != nil {
			parentDept, temp = p.departmentRepository.GetById(ctx, *parentId)
			if temp != nil {
				return temp
			}
		}
		dept := &models.Department{
			Name:      name,
			CreatorId: operator,
			UpdaterId: operator,
		}
		// 完善 Parent & Ancestors
		if parentDept != nil {
			parentDeptAncestors := ""
			if parentDept.Ancestors != nil {
				parentDeptAncestors = *parentDept.Ancestors
			}
			t := lo.Filter([]string{parentDeptAncestors, strconv.Itoa(int(parentDept.ID))}, func(item string, _ int) bool {
				return item != ""
			})
			ancestors := strings.Join(t, ",")
			dept.Ancestors = &ancestors
			dept.ParentId = parentId
		}
		// 查询有效的用户
		var users []*models.User
		tempUserIds := lo.Union(userIds, leaderIds)
		if len(tempUserIds) > 0 {
			users, err = p.userRepository.GetByIds(ctx, tempUserIds, false, false, false)
			if err != nil {
				return err
			}
			if len(users) > 0 {
				// leader 也属于部门的成员
				dept.Users = users
				// 设置部门 leader
				if len(leaderIds) > 0 {
					dept.Leaders = lo.Filter(users, func(item *models.User, _ int) bool {
						return lo.SomeBy(leaderIds, func(uid uint) bool {
							return uid == item.ID
						})
					})
				}
			}
		}
		// 创建部门 & 建立关联关系
		return p.departmentRepository.Create(ctx, dept)
	})
}

type departmentBaseData struct {
	departments []*models.Department
	deptLeader  []*models.DepartmentLeader
	userDept    []*models.UserDepartment
}

// TODO:添加缓存优化 single flight
func (p *DepartmentService) GetAllDepartment(ctx context.Context, crew bool) ([]*models.Department, error) {
	departments, err := p.departmentRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if crew {
		var users []*models.User
		var deptLeader []*models.DepartmentLeader
		var userDept []*models.UserDepartment
		deptLeader, err = p.departmentRepository.GetAllDepartmentLeader(ctx)
		if err != nil {
			return nil, err
		}
		userDept, err = p.departmentRepository.GetAllUserDepartment(ctx)
		if err != nil {
			return nil, err
		}
		users, err = p.userRepository.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		for _, dept := range departments {
			tempUserDept := lo.Filter(userDept, func(item *models.UserDepartment, _ int) bool {
				return item.DepartmentId == dept.ID
			})
			userIds := lo.Map(tempUserDept, func(item *models.UserDepartment, _ int) uint {
				return item.UserId
			})
			tempDeptLeader := lo.Filter(deptLeader, func(item *models.DepartmentLeader, _ int) bool {
				return item.DepartmentId == dept.ID
			})
			leaderIds := lo.Map(tempDeptLeader, func(item *models.DepartmentLeader, _ int) uint {
				return item.UserId
			})
			dept.Users = lo.Filter(users, func(item *models.User, _ int) bool {
				return lo.Contains(userIds, item.ID)
			})
			dept.Leaders = lo.Filter(users, func(item *models.User, _ int) bool {
				return lo.Contains(leaderIds, item.ID)
			})
			creator, ok := lo.Find(users, func(item *models.User) bool {
				return item.ID == dept.CreatorId
			})
			if ok {
				dept.Creator = *creator
			} else {
				return nil, response.Error
			}
			updater, ok := lo.Find(users, func(item *models.User) bool {
				return item.ID == dept.UpdaterId
			})
			if ok {
				dept.Updater = *updater
			} else {
				return nil, response.Error
			}
		}
	}
	return departments, nil
}
