package dao

import (
	"context"
	"gin-web/pkg/database"
	"gin-web/pkg/global"
	"gorm.io/gorm"
	"sync"
)

var (
	basicDAO     *BasicDAO
	basicDAOOnce sync.Once
)

type BasicDAO struct {
}

func NewBasicDao() *BasicDAO {
	basicDAOOnce.Do(func() {
		basicDAO = &BasicDAO{}
	})
	return basicDAO
}

func (basic *BasicDAO) Datasource(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return global.DB
	}
	if manager := database.LoadManager(ctx); manager != nil {
		return manager.DB
	}
	return global.DB
}
