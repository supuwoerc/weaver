package models

import (
	"github.com/supuwoerc/weaver/pkg/database"
)

type Attachment struct {
	Name      string `json:"name"`
	CreatorId *uint  `json:"-"` // nil:系统创建
	Creator   *User  `json:"creator" gorm:"foreignKey:CreatorId;references:ID"`
	Type      int8   `json:"type"`
	Size      int64  `json:"size"`
	Hash      string `json:"hash"` // 文件内容hash,作为文件的存储名称
	Path      string `json:"path"`
	database.BasicModel
}
