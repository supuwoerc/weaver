package dao

import (
	"context"
	"gin-web/models"
	"sync"
)

var (
	attachmentDAO     *AttachmentDAO
	attachmentDAOOnce sync.Once
)

type AttachmentDAO struct {
	*BasicDAO
}

func NewAttachmentDAO() *AttachmentDAO {
	attachmentDAOOnce.Do(func() {
		attachmentDAO = &AttachmentDAO{
			BasicDAO: NewBasicDao(),
		}
	})
	return attachmentDAO
}

func (a *AttachmentDAO) Insert(ctx context.Context, records []*models.Attachment) error {
	err := a.Datasource(ctx).Create(records).Error
	return err
}

func (a *AttachmentDAO) GetIsExistByHash(ctx context.Context, hash string) (bool, error) {
	var count int64
	err := a.Datasource(ctx).Where("hash = ?", hash).Count(&count).Error
	return count > 0, err
}
