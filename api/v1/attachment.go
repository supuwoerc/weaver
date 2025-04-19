package v1

import (
	"context"
	"gin-web/conf"
	"gin-web/middleware"
	"gin-web/pkg/constant"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type AttachmentService interface {
	SaveFiles(ctx context.Context, files []*multipart.FileHeader, uid uint) ([]*response.UploadAttachmentResponse, error)
	SaveFile(ctx context.Context, file *multipart.FileHeader, uid uint) (*response.UploadAttachmentResponse, error)
}

type AttachmentApi struct {
	service AttachmentService
	conf    *conf.Config
}

func NewAttachmentApi(
	route *gin.RouterGroup,
	service AttachmentService,
	authMiddleware *middleware.AuthMiddleware,
	conf *conf.Config,
) *AttachmentApi {
	// 初始化controller
	attachmentApi := &AttachmentApi{
		conf:    conf,
		service: service,
	}
	// 挂载路由
	attachmentAccessGroup := route.Group("attachment").Use(authMiddleware.LoginRequired())
	{
		attachmentAccessGroup.POST("multiple-upload", attachmentApi.MultipleUpload)
		attachmentAccessGroup.POST("upload", attachmentApi.Upload)
	}
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
	maxUploadLength := a.conf.System.MaxUploadLength
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
