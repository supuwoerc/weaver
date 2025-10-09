package response

type UploadAttachmentResponse struct {
	ID   uint   `json:"id"`   // 文件ID
	Name string `json:"name"` // 文件名
	Size int64  `json:"size"` // 文件大小
	Path string `json:"path"` // 文件路径
}
