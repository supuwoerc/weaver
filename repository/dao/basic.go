package dao

import (
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BasicDAO struct {
	db  *gorm.DB
	ctx *gin.Context
}

func NewBasicDao(ctx *gin.Context, db *gorm.DB) *BasicDAO {
	if db == nil {
		db = global.DB
	}
	return &BasicDAO{db: db, ctx: ctx}
}

func (b *BasicDAO) Transaction(fn func(dao *BasicDAO, tx *gorm.DB) error) error {
	return b.db.Transaction(func(tx *gorm.DB) error {
		basicDAO := NewBasicDao(b.ctx, tx)
		return fn(basicDAO, tx)
	})
}
