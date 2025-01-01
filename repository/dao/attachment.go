package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type AttachmentDAO struct {
	*BasicDAO
}

type Attachment struct {
	gorm.Model
	Name string         // 上传的文件名
	Uid  sql.Null[uint] // 上传的用户
	Type int8           // 文件类型
	Size int64          // 文件大小
	Hash string         // 文件摘要
	Path string         // 文件路径
}

func NewAttachmentDAO() *AttachmentDAO {
	return &AttachmentDAO{BasicDAO: NewBasicDao()}
}

func (a *AttachmentDAO) Insert(ctx context.Context, records []*Attachment) ([]*Attachment, error) {
	err := a.Datasource(ctx).Create(records).Error
	return records, err
}

func (a *AttachmentDAO) IsExistByHash(ctx context.Context, hash string) (bool, error) {
	var count int64
	err := a.Datasource(ctx).Where("hash = ?", hash).Count(&count).Error
	return count > 0, err
}
