package v1

import (
	"gin-web/models"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"path/filepath"
	"sync"
)

type AttachmentApi struct {
	*BasicApi
	service *service.AttachmentService
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
	ret := lo.Map(result, func(item *models.Attachment, _ int) *response.UploadAttachmentResponse {
		return &response.UploadAttachmentResponse{
			ID:   item.ID,
			Name: item.Name,
			Size: item.Size,
			Path: filepath.Join(string(filepath.Separator), item.Path), // TODO:返回预览/下载接口完整路径
		}
	})
	response.SuccessWithData(ctx, ret)
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
	response.SuccessWithData(ctx, &response.UploadAttachmentResponse{
		ID:   result.ID,
		Name: result.Name,
		Size: result.Size,
		Path: result.Path,
	})
}
