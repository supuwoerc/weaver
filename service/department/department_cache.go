package department

import (
	"context"
	"fmt"
	"time"

	"github.com/supuwoerc/weaver/pkg/constant"
)

func (p *Service) CacheKey() string {
	return constant.AutoManageDeptCache
}

func (p *Service) RefreshCache(ctx context.Context) error {
	start := time.Now()
	p.Logger.WithContext(ctx).Infow("refresh department", "begin", start.Format(time.DateTime))
	defer func() {
		p.Logger.WithContext(ctx).Infow("refresh department",
			"end", time.Now().Format(time.DateTime), "cost",
			fmt.Sprintf("%dms", time.Since(start).Milliseconds()),
		)
	}()
	_, err, _ := p.deptTreeSfg.Do(string(constant.DepartmentTreeRefreshSfgKey), func() (interface{}, error) {
		departments, err := p.departmentDAO.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		sfgKey := constant.DepartmentTreeSfgKey
		crewSfgKey := constant.DepartmentTreeWithCrewSfgKey
		if err = p.departmentCache.CacheDepartment(ctx, sfgKey, departments); err != nil {
			return nil, err
		}
		if err = p.processDepartmentWithCrew(ctx, departments); err != nil {
			return nil, err
		}
		return nil, p.departmentCache.CacheDepartment(ctx, crewSfgKey, departments)
	})
	return err
}

func (p *Service) CleanCache(ctx context.Context) error {
	start := time.Now()
	p.Logger.WithContext(ctx).Infow("clean department", "begin", start.Format(time.DateTime))
	defer func() {
		p.Logger.WithContext(ctx).Infow("clean department",
			"end", time.Now().Format(time.DateTime), "cost",
			fmt.Sprintf("%dms", time.Since(start).Milliseconds()),
		)
	}()
	_, err, _ := p.deptTreeSfg.Do(string(constant.DepartmentTreeCleanSfgKey), func() (interface{}, error) {
		keys := []constant.CacheKey{constant.DepartmentTreeSfgKey, constant.DepartmentTreeWithCrewSfgKey}
		return nil, p.departmentCache.RemoveDepartmentCache(ctx, keys...)
	})
	return err
}
