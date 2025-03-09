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

func NewAttachmentDAO(basicDAO *BasicDAO) *AttachmentDAO {
	attachmentDAOOnce.Do(func() {
		attachmentDAO = &AttachmentDAO{
			BasicDAO: basicDAO,
		}
	})
	return attachmentDAO
}

func (a *AttachmentDAO) Create(ctx context.Context, records []*models.Attachment) error {
	return a.Datasource(ctx).Create(records).Error
}
