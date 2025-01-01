package dao

import (
	"context"
	"gin-web/pkg/database"
	"gin-web/pkg/global"
	"gorm.io/gorm"
)

type BasicDAO struct {
}

func NewBasicDao() *BasicDAO {
	return &BasicDAO{}
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
