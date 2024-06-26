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

var basicDao *BasicDAO

func NewBasicDao(ctx *gin.Context) *BasicDAO {
	if basicDao == nil {
		basicDao = &BasicDAO{db: global.DB, ctx: ctx}
	}
	return basicDao
}
