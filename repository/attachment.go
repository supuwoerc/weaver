package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
	"gorm.io/gorm"
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

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	attachmentRepositoryOnce.Do(func() {
		attachmentRepository = &AttachmentRepository{
			dao: dao.NewAttachmentDAO(db),
		}
	})
	return attachmentRepository
}

func (r *AttachmentRepository) Create(ctx context.Context, records []*models.Attachment) error {
	return r.dao.Create(ctx, records)
}
