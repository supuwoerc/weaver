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

func (a *AttachmentDAO) Create(ctx context.Context, records []*models.Attachment) error {
	return a.Datasource(ctx).Create(records).Error
}
