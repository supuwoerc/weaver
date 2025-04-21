package repository

import (
	"context"
	"gin-web/models"
)

type AttachmentDAO interface {
	Create(ctx context.Context, records []*models.Attachment) error
}

type AttachmentRepository struct {
	AttachmentDAO
}

func NewAttachmentRepository(dao AttachmentDAO) *AttachmentRepository {
	return &AttachmentRepository{
		AttachmentDAO: dao,
	}
}
