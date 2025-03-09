package repository

import (
	"context"
	"gin-web/models"
	"sync"
)

var (
	attachmentRepository     *AttachmentRepository
	attachmentRepositoryOnce sync.Once
)

type AttachmentDAO interface {
	Create(ctx context.Context, records []*models.Attachment) error
}

type AttachmentRepository struct {
	dao AttachmentDAO
}

func NewAttachmentRepository(dao AttachmentDAO) *AttachmentRepository {
	attachmentRepositoryOnce.Do(func() {
		attachmentRepository = &AttachmentRepository{
			dao: dao,
		}
	})
	return attachmentRepository
}

func (r *AttachmentRepository) Create(ctx context.Context, records []*models.Attachment) error {
	return r.dao.Create(ctx, records)
}
