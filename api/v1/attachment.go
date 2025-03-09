package v1

import (
	"context"
	"gin-web/middleware"
	"gin-web/pkg/constant"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mime/multipart"
	"sync"
)

type AttachmentService interface {
	SaveFiles(ctx context.Context, files []*multipart.FileHeader, uid uint) ([]*response.UploadAttachmentResponse, error)
	SaveFile(ctx context.Context, file *multipart.FileHeader, uid uint) (*response.UploadAttachmentResponse, error)
}

type AttachmentApi struct {
	viper   *viper.Viper
	service AttachmentService
}

var (
	attachmentOnce sync.Once
	attachmentApi  *AttachmentApi
)

func NewAttachmentApi(
	route *gin.RouterGroup,
	service AttachmentService,
	authMiddleware *middleware.AuthMiddleware,
	viper *viper.Viper,
) *AttachmentApi {
	attachmentOnce.Do(func() {
		// 初始化controller
		attachmentApi = &AttachmentApi{
			viper:   viper,
			service: service,
		}
		// 挂载路由
		attachmentAccessGroup := route.Group("attachment").Use(authMiddleware.LoginRequired())
		{
			attachmentAccessGroup.POST("multiple-upload", attachmentApi.MultipleUpload)
			attachmentAccessGroup.POST("upload", attachmentApi.Upload)
		}
	})
	return attachmentApi
}

func (a *AttachmentApi) MultipleUpload(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	files := form.File["files"]
	fileLen := len(files)
	maxUploadLength := a.viper.GetInt("system.maxUploadLength")
	if maxUploadLength == 0 {
		maxUploadLength = constant.DefaultMaxLength
	}
	if fileLen == 0 || fileLen > maxUploadLength {
		response.FailWithCode(ctx, response.InvalidAttachmentLength)
		return
	}
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	result, err := a.service.SaveFiles(ctx, files, claims.User.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	if result == nil {
		response.FailWithCode(ctx, response.Error)
		return
	}
	response.SuccessWithData(ctx, result)
}

func (a *AttachmentApi) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	result, err := a.service.SaveFile(ctx, file, claims.User.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	if result == nil {
		response.FailWithCode(ctx, response.Error)
		return
	}
	response.SuccessWithData(ctx, result)
}
