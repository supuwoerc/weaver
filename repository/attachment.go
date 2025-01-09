package repository

import (
	"context"
	"gin-web/models"
	"gin-web/repository/dao"
	"gin-web/repository/transducer"
	"github.com/samber/lo"
	"sync"
)

var (
	attachmentRepository     *AttachmentRepository
	attachmentRepositoryOnce sync.Once
)

type AttachmentRepository struct {
	dao *dao.AttachmentDAO
}

func NewAttachmentRepository() *AttachmentRepository {
	attachmentRepositoryOnce.Do(func() {
		attachmentRepository = &AttachmentRepository{
			dao: dao.NewAttachmentDAO(),
		}
	})
	return attachmentRepository
}

func toModelAttachment(record *dao.Attachment) *models.Attachment {
	return &models.Attachment{
		ID:   record.ID,
		Name: record.Name,
		Uid:  transducer.NullValue(record.Uid),
		Type: record.Type,
		Size: record.Size,
		Hash: record.Hash,
		Path: record.Path,
	}
}

func toModelAttachments(records []*dao.Attachment) []*models.Attachment {
	return lo.Map[*dao.Attachment, *models.Attachment](records, func(item *dao.Attachment, _ int) *models.Attachment {
		return toModelAttachment(item)
	})
}

func (r *AttachmentRepository) Create(ctx context.Context, records []*models.Attachment) ([]*models.Attachment, error) {
	attachments := lo.Map[*models.Attachment, *dao.Attachment](records, func(item *models.Attachment, _ int) *dao.Attachment {
		return &dao.Attachment{
			Name: item.Name,
			Uid:  transducer.SqlNullValue(item.Uid),
			Type: item.Type,
			Size: item.Size,
			Hash: item.Hash,
			Path: item.Path,
		}
	})
	ret, err := r.dao.Insert(ctx, attachments)
	return toModelAttachments(ret), err
}

func (r *AttachmentRepository) IsExistByHash(ctx context.Context, hash string) (bool, error) {
	return r.dao.IsExistByHash(ctx, hash)
}
