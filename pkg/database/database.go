package database

import (
	"context"
	"gin-web/pkg/constant"
	"gorm.io/gorm"
	"strings"
	"time"
)

type BasicModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Action func(ctx context.Context) error

type TransactionManager struct {
	DB                           *gorm.DB
	AlreadyCommittedOrRolledBack bool // 是否已经提交或者回滚了
}

func LoadManager(ctx context.Context) *TransactionManager {
	value := ctx.Value(constant.TransactionKey)
	if value == nil {
		return nil
	} else {
		if m, ok := value.(*TransactionManager); !ok {
			return nil
		} else {
			return m
		}
	}
}

func InjectManager(ctx context.Context, m any) context.Context {
	return context.WithValue(ctx, constant.TransactionKey, m)
}

func FuzzKeyword(s string) string {
	if s == "" {
		return ""
	}
	str := strings.ReplaceAll(strings.ReplaceAll(s, "%", "\\%"), "_", "\\_")
	return "%" + str + "%"
}
