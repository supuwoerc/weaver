package dao

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttachmentDAO struct {
	*BasicDAO
}

type Attachment struct {
	gorm.Model
	Name string // 上传的文件名
	Uid  uint   // 上传的用户
	Type int8   // 文件类型
	Size int64  // 文件大小
	Hash string // 文件摘要
}

func NewAttachmentDAO(ctx *gin.Context) *AttachmentDAO {
	return &AttachmentDAO{BasicDAO: NewBasicDao(ctx)}
}

func (a *AttachmentDAO) Insert(ctx context.Context, records []*Attachment) ([]*Attachment, error) {
	err := a.db.WithContext(ctx).Create(records).Error
	return records, err
}

func (a *AttachmentDAO) IsExistByHash(ctx context.Context, hash string) (bool, error) {
	var count int64
	err := a.db.WithContext(ctx).Where("hash = ?", hash).Count(&count).Error
	return count > 0, err
}
