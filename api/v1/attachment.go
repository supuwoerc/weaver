package v1

import (
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

// TODO：补充文档
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
	user, err := utils.GetContextUser(ctx)
	if err != nil || user == nil {
		response.FailWithCode(ctx, response.UserNotExist)
		return
	}
	result, err := a.service(ctx).SaveFiles(files, user.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, result)
}
