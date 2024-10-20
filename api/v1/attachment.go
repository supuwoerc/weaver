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
)

type AttachmentApi struct {
	*BasicApi
	service func(ctx *gin.Context) *service.AttachmentService
}

const (
	defaultMaxLength = 50
)

func NewAttachmentApi() AttachmentApi {
	return AttachmentApi{
		BasicApi: NewBasicApi(),
		service: func(ctx *gin.Context) *service.AttachmentService {
			return service.NewAttachmentService(ctx)
		},
	}
}

// MultipleUpload TODO：补充文档
func (a AttachmentApi) MultipleUpload(ctx *gin.Context) {
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
		response.FailWithCode(ctx, response.UserNotExist)
		return
	}
	result, err := a.service(ctx).SaveFiles(files, claims.User.UID, nil)
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
			Path: filepath.Join(string(filepath.Separator), item.Path),
		}
	})
	response.SuccessWithData(ctx, ret)
}

// Upload TODO：补充文档
func (a AttachmentApi) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.UserNotExist)
		return
	}
	result, err := a.service(ctx).SaveFile(file, claims.User.UID, nil)
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
