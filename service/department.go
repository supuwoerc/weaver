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
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"golang.org/x/sync/singleflight"
	"strconv"
	"strings"
	"sync"
)

type DepartmentService struct {
	*BasicService
	departmentRepository *repository.DepartmentRepository
	userRepository       *repository.UserRepository
	deptTreeSfg          singleflight.Group
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

func (p *DepartmentService) GetDepartmentTree(ctx context.Context, crew bool) ([]*response.DepartmentTreeResponse, error) {
	key := constant.DepartmentTreeSfgKey
	if crew {
		key = constant.DepartmentTreeCrewSfgKey
	}
	departmentCache, cacheErr := p.processDepartmentCache(ctx, key)
	if cacheErr != nil {
		return nil, cacheErr
	}
	if departmentCache != nil {
		return p.processTree(departmentCache)
	}
	result, err, _ := p.deptTreeSfg.Do(key, func() (interface{}, error) {
		departments, err := p.departmentRepository.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		if crew {
			if err = p.processDepartmentCrew(ctx, departments); err != nil {
				return nil, err
			}
		}
		if err = p.departmentRepository.CacheDepartment(ctx, key, departments); err != nil {
			return nil, err
		}
		return p.processTree(departments), nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*response.DepartmentTreeResponse), nil
}

func (p *DepartmentService) processDepartmentCache(ctx context.Context, key string) ([]*models.Department, error) {
	cache, err := p.departmentRepository.GetDepartmentCache(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return cache, nil
}

func (p *DepartmentService) processDepartmentCrew(ctx context.Context, departments []*models.Department) error {
	var users []*models.User
	var deptLeader []*models.DepartmentLeader
	var userDept []*models.UserDepartment
	var err error
	deptLeader, err = p.departmentRepository.GetAllDepartmentLeader(ctx)
	if err != nil {
		return err
	}
	userDept, err = p.departmentRepository.GetAllUserDepartment(ctx)
	if err != nil {
		return err
	}
	users, err = p.userRepository.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, dept := range departments {
		ud := lo.Filter(userDept, func(item *models.UserDepartment, _ int) bool {
			return item.DepartmentId == dept.ID
		})
		dl := lo.Filter(deptLeader, func(item *models.DepartmentLeader, _ int) bool {
			return item.DepartmentId == dept.ID
		})
		dept.Users = lo.Filter(users, func(item *models.User, _ int) bool {
			return lo.Contains(lo.Map(ud, func(item *models.UserDepartment, _ int) uint {
				return item.UserId
			}), item.ID)
		})
		dept.Leaders = lo.Filter(users, func(item *models.User, _ int) bool {
			return lo.Contains(lo.Map(dl, func(item *models.DepartmentLeader, _ int) uint {
				return item.UserId
			}), item.ID)
		})
		creator, ok := lo.Find(users, func(item *models.User) bool {
			return item.ID == dept.CreatorId
		})
		if ok {
			dept.Creator = *creator
		} else {
			return response.Error
		}
		updater, ok := lo.Find(users, func(item *models.User) bool {
			return item.ID == dept.UpdaterId
		})
		if ok {
			dept.Updater = *updater
		} else {
			return response.Error
		}
	}
	return nil
}

func (p *DepartmentService) processTree(departments []*models.Department) ([]*response.DepartmentTreeResponse, error) {
	var res []*response.DepartmentTreeResponse
	nodeMap := make(map[uint]*response.DepartmentTreeResponse)
	deptMap := make(map[uint]*models.Department)
	for _, dept := range departments {
		deptMap[dept.ID] = dept
	}
	for _, dept := range departments {
		holder, exist := nodeMap[dept.ID]
		var children = make([]*response.DepartmentTreeResponse, 0)
		if exist {
			children = holder.Children
		}
		node, parseErr := response.ToDepartmentTreeResponse(dept, deptMap)
		if parseErr != nil {
			return nil, parseErr
		}
		node.Children = children
		nodeMap[node.ID] = node
		if dept.ParentId == nil {
			res = append(res, node)
		} else {
			_, exist = nodeMap[*dept.ParentId]
			if !exist {
				nodeMap[*dept.ParentId], parseErr = response.ToDepartmentTreeResponse(&models.Department{}, deptMap)
				if parseErr != nil {
					return nil, parseErr
				}
			}
			nodeMap[*dept.ParentId].Children = append(nodeMap[*dept.ParentId].Children, node)
		}
	}
	return res, nil
}
