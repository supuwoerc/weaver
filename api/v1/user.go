package v1

import (
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
	BasicApi
}

func NewUserApi() UserApi {
	return UserApi{
		NewBasicApi(),
	}
}

// 注册
func (u UserApi) SignUp(ctx *gin.Context) {
	response.Success(ctx)
}
