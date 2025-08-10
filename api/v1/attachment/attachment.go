package attachment

import (
	"context"
	"mime/multipart"

	v1 "github.com/supuwoerc/weaver/api/v1"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Service interface {
	SaveFiles(ctx context.Context, files []*multipart.FileHeader, uid uint) ([]*response.UploadAttachmentResponse, error)
	SaveFile(ctx context.Context, file *multipart.FileHeader, uid uint) (*response.UploadAttachmentResponse, error)
}

type Api struct {
	*v1.BasicApi
	service Service
}

func NewAttachmentApi(basic *v1.BasicApi, service Service) *Api {
	attachmentApi := &Api{
		BasicApi: basic,
		service:  service,
	}
	// 挂载路由
	attachmentAccessGroup := basic.Route.Group("attachment").Use(basic.Auth.LoginRequired())
	{
		attachmentAccessGroup.POST("multiple-upload", attachmentApi.MultipleUpload)
		attachmentAccessGroup.POST("upload", attachmentApi.Upload)
	}
	return attachmentApi
}

func (a *Api) MultipleUpload(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	files := form.File["files"]
	fileLen := len(files)
	maxUploadLength := a.Conf.System.MaxUploadLength
	if maxUploadLength == 0 {
		maxUploadLength = constant.DefaultMaxLength
	}
	if fileLen == 0 || fileLen > maxUploadLength {
		response.FailWithCode(ctx, response.InvalidAttachmentLength)
		return
	}
	// TODO:集中到鉴权中间件
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

func (a *Api) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	// TODO:集中到鉴权中间件
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
