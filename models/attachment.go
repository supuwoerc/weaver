package models

type Attachment struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Uid  uint   `json:"uid"`
	Type int8   `json:"type"`
	Size int64  `json:"size"`
	Hash string `json:"hash"` // 文件内容hash,作为文件的存储名称
}
