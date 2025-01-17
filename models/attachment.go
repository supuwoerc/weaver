package models

import (
	"gin-web/pkg/database"
)

type Attachment struct {
	Name string `json:"name"`
	Uid  *uint  `json:"uid"` // nil:系统创建
	Type int8   `json:"type"`
	Size int64  `json:"size"`
	Hash string `json:"hash"` // 文件内容hash,作为文件的存储名称
	Path string `json:"path"`
	database.BasicModel
}
