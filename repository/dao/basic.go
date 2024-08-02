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

func NewBasicDao(ctx *gin.Context) *BasicDAO {
	return &BasicDAO{db: global.DB, ctx: ctx}
}
