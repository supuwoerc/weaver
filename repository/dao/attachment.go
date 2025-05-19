package dao

import (
	"context"

	"github.com/supuwoerc/weaver/models"
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
