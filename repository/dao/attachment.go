package dao

import (
	"context"
	"gin-web/models"
)

type AttachmentDAO struct {
	*BasicDAO
}

func NewAttachmentDAO(basicDAO *BasicDAO) *AttachmentDAO {
	return &AttachmentDAO{
		BasicDAO: basicDAO,
	}
}

func (a *AttachmentDAO) Create(ctx context.Context, records []*models.Attachment) error {
	return a.Datasource(ctx).Create(records).Error
}
