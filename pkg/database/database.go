package database

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gin-web/pkg/constant"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UpsertTime time.Time

func (c UpsertTime) MarshalJSON() ([]byte, error) {
	format := time.Time(c).Format(time.DateTime)
	return json.Marshal(format)
}

func (c UpsertTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	t := time.Time(c)
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t, nil
}

func (c *UpsertTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*c = UpsertTime(value)
		return nil
	}
	return fmt.Errorf("[UpsertTime] can not convert %v to timestamp", v)
}

type BasicModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt UpsertTime     `json:"created_at"`
	UpdatedAt UpsertTime     `json:"updated_at"`
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
