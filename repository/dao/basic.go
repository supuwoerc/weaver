package dao

import (
	"context"

	"github.com/supuwoerc/weaver/pkg/database"

	"gorm.io/gorm"
)

type BasicDAO struct {
	DB         *gorm.DB
	QueryLimit int
}

func NewBasicDao(db *gorm.DB) *BasicDAO {
	return &BasicDAO{
		DB:         db,
		QueryLimit: 500,
	}
}

func (basic *BasicDAO) Datasource(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return basic.DB
	}
	if manager := database.LoadManager(ctx); manager != nil {
		return manager.DB.WithContext(ctx)
	}
	return basic.DB.WithContext(ctx)
}

func (basic *BasicDAO) Preload(field string, args ...any) func(d *gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Preload(field, args...)
	}
}

func queryAll[T any](db *gorm.DB, limit int) ([]T, error) {
	var result []T
	offset := 0
	for {
		var page []T
		err := db.Limit(limit).Offset(offset).Find(&page).Error
		if err != nil {
			return nil, err
		}
		if len(page) > 0 {
			result = append(result, page...)
		}
		if len(page) < limit {
			break
		}
		offset += limit
	}
	return result, nil
}
