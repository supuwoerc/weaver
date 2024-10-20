package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
)

func InitAttachmentRouter(r *gin.RouterGroup) {
	attachmentApi := v1.NewAttachmentApi()
	attachmentAccessGroup := r.Group("attachment")
	{
		attachmentAccessGroup.POST("multiple-upload", attachmentApi.MultipleUpload)
		attachmentAccessGroup.POST("upload", attachmentApi.Upload)
	}
}
