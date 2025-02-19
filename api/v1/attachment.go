package v1

import (
	"context"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/service"
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
	*BasicApi
	service AttachmentService
}

const (
	defaultMaxLength = 50
)

var (
	attachmentOnce sync.Once
	attachmentApi  *AttachmentApi
)

func NewAttachmentApi() *AttachmentApi {
	attachmentOnce.Do(func() {
		attachmentApi = &AttachmentApi{
			BasicApi: NewBasicApi(),
			service:  service.NewAttachmentService(),
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
	maxUploadLength := viper.GetInt("system.maxUploadLength")
	if maxUploadLength == 0 {
		maxUploadLength = defaultMaxLength
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
