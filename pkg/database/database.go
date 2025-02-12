package database

import (
	"context"
	"database/sql/driver"
	"fmt"
	"gin-web/pkg/constant"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"strings"
	"time"
)

type UpsertTime time.Time

func (c *UpsertTime) UnmarshalJSON(bytes []byte) error {
	str := string(bytes)
	if str == "null" || str == `""` {
		return nil
	}
	if parse, err := time.Parse(`"`+time.DateTime+`"`, str); err != nil {
		return err
	} else {
		*c = UpsertTime(parse)
		return nil
	}
}

func (c UpsertTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.DateTime)+len(`""`))
	b = append(b, '"')
	b = time.Time(c).AppendFormat(b, time.DateTime)
	b = append(b, '"')
	return b, nil
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
	ID        uint                  `json:"id" gorm:"primarykey"`
	CreatedAt UpsertTime            `json:"created_at"`
	UpdatedAt UpsertTime            `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `json:"-" gorm:"softDelete:milli;index"`
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
