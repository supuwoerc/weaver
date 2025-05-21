package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/request"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"golang.org/x/sync/singleflight"
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

type DepartmentService struct {
	*BasicService
	departmentDAO   DepartmentDAO
	departmentCache DepartmentCache
	userDAO         UserDAO
	deptTreeSfg     singleflight.Group
}

func NewDepartmentService(
	basic *BasicService,
	deptDAO DepartmentDAO,
	deptCache DepartmentCache,
	userDAO UserDAO,
) *DepartmentService {
	return &DepartmentService{
		BasicService:    basic,
		departmentDAO:   deptDAO,
		departmentCache: deptCache,
		userDAO:         userDAO,
	}
}

func (p *DepartmentService) lockDepartmentField(ctx context.Context, name string, parentId *uint) ([]*utils.RedisLock, error) {
	locks := make([]*utils.RedisLock, 0)
	// 名称锁
	deptNameLock := p.locksmith.NewLock(constant.DepartmentNamePrefix, name)
	if err := deptNameLock.Lock(ctx, true); err != nil {
		return locks, err
	}
	locks = append(locks, deptNameLock)
	if parentId != nil {
		// 父部门锁
		parentLock := p.locksmith.NewLock(constant.DepartmentIdPrefix, strconv.Itoa(int(*parentId)))
		if err := parentLock.Lock(ctx, true); err != nil {
			return locks, err
		}
		locks = append(locks, parentLock)
	}
	return locks, nil
}

func (p *DepartmentService) CreateDepartment(ctx context.Context, operator uint, params *request.CreateDepartmentRequest) error {
	locks, err := p.lockDepartmentField(ctx, params.Name, params.ParentId)
	defer func() {
		for _, l := range locks {
			if e := l.Unlock(); e != nil {
				p.logger.Errorf("unlock fail %s", e.Error())
			}
		}
	}()
	if err != nil {
		return err
	}
	return p.Transaction(ctx, false, func(ctx context.Context) error {
		// 查询是否重复
		existDept, temp := p.departmentDAO.GetByName(ctx, params.Name)
		if temp != nil && !errors.Is(temp, response.DeptNotExist) {
			return temp
		}
		if existDept != nil {
			return response.DeptCreateDuplicate
		}
		// 查询父部门
		var parentDept *models.Department
		if params.ParentId != nil {
			parentDept, temp = p.departmentDAO.GetById(ctx, *params.ParentId)
			if temp != nil {
				return temp
			}
		}
		dept := &models.Department{
			Name:      params.Name,
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
			dept.ParentId = params.ParentId
		}
		// 查询有效的用户
		var users []*models.User
		tempUserIds := lo.Union(params.Users, params.Leaders)
		if len(tempUserIds) > 0 {
			users, err = p.userDAO.GetByIds(ctx, tempUserIds)
			if err != nil {
				return err
			}
			if len(users) > 0 {
				// leader 也属于部门的成员
				dept.Users = users
				// 设置部门 leader
				if len(params.Leaders) > 0 {
					dept.Leaders = lo.Filter(users, func(item *models.User, _ int) bool {
						return lo.SomeBy(params.Leaders, func(uid uint) bool {
							return uid == item.ID
						})
					})
				}
			}
		}
		// 创建部门 & 建立关联关系
		return p.departmentDAO.Create(ctx, dept)
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
	result, err, _ := p.deptTreeSfg.Do(string(key), func() (interface{}, error) {
		departments, err := p.departmentDAO.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		if crew {
			if err = p.processDepartmentCrew(ctx, departments); err != nil {
				return nil, err
			}
		}
		if err = p.departmentCache.CacheDepartment(ctx, key, departments); err != nil {
			return nil, err
		}
		return p.processTree(departments)
	})
	if err != nil {
		return nil, err
	}
	return result.([]*response.DepartmentTreeResponse), nil
}

func (p *DepartmentService) processDepartmentCache(ctx context.Context, key constant.CacheKey) ([]*models.Department, error) {
	cache, err := p.departmentCache.GetDepartmentCache(ctx, key)
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
	deptLeader, err = p.departmentDAO.GetAllDepartmentLeader(ctx)
	if err != nil {
		return err
	}
	userDept, err = p.departmentDAO.GetAllUserDepartment(ctx)
	if err != nil {
		return err
	}
	users, err = p.userDAO.GetAll(ctx)
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

func (p *DepartmentService) Key() string {
	return constant.AutoManageDeptCache
}

func (p *DepartmentService) Refresh(ctx context.Context) error {
	start := time.Now()
	p.logger.Infow("refresh department", "begin", start.Format(time.DateTime))
	defer func() {
		p.logger.Infow("refresh department",
			"end", time.Now().Format(time.DateTime), "cost",
			fmt.Sprintf("%dms", time.Since(start).Milliseconds()),
		)
	}()
	departments, err := p.departmentDAO.GetAll(ctx)
	if err != nil {
		return err
	}
	sfgKey := constant.DepartmentTreeSfgKey
	crewSfgKey := constant.DepartmentTreeCrewSfgKey
	if err = p.departmentCache.CacheDepartment(ctx, sfgKey, departments); err != nil {
		return err
	}
	if err = p.processDepartmentCrew(ctx, departments); err != nil {
		return err
	}
	return p.departmentCache.CacheDepartment(ctx, crewSfgKey, departments)
}

func (p *DepartmentService) Clean(ctx context.Context) error {
	start := time.Now()
	p.logger.Infow("clean department", "begin", start.Format(time.DateTime))
	defer func() {
		p.logger.Infow("clean department",
			"end", time.Now().Format(time.DateTime), "cost",
			fmt.Sprintf("%dms", time.Since(start).Milliseconds()),
		)
	}()
	return p.departmentCache.RemoveDepartmentCache(ctx, constant.DepartmentTreeSfgKey, constant.DepartmentTreeCrewSfgKey)
}
