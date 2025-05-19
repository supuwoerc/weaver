package dao

import (
	"context"
	"gin-web/pkg/database"

	"gorm.io/gorm"
)

type BasicDAO struct {
	DB *gorm.DB
}

func NewBasicDao(db *gorm.DB) *BasicDAO {
	return &BasicDAO{
		DB: db,
	}
}

func (basic *BasicDAO) Datasource(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return basic.DB
	}
	if manager := database.LoadManager(ctx); manager != nil {
		return manager.DB
	}
	return basic.DB
}

func (basic *BasicDAO) Preload(field string, args ...any) func(d *gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Preload(field, args...)
	}
}
