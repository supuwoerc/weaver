package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
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

func NewAttachmentRepository() *AttachmentRepository {
	attachmentRepositoryOnce.Do(func() {
		attachmentRepository = &AttachmentRepository{
			dao: dao.NewAttachmentDAO(),
		}
	})
	return attachmentRepository
}

func (r *AttachmentRepository) Create(ctx context.Context, records []*models.Attachment) error {
	return r.dao.Create(ctx, records)
}
