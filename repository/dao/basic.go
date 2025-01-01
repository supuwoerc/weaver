package dao

import (
	"context"
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
	return global.DB
}
