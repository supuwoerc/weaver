package repository

import (
	"context"
	"gin-web/models"
)

type AttachmentDAO interface {
	Create(ctx context.Context, records []*models.Attachment) error
}

type AttachmentRepository struct {
	dao AttachmentDAO
}

func NewAttachmentRepository(dao AttachmentDAO) *AttachmentRepository {
	return &AttachmentRepository{
		dao: dao,
	}
}

func (r *AttachmentRepository) Create(ctx context.Context, records []*models.Attachment) error {
	return r.dao.Create(ctx, records)
}
