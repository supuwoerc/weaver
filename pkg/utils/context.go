package utils

import (
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetContextUser 从上下文中获取当前请求接口的用户
func GetContextUser(ctx *gin.Context) (*models.User, error) {
	value, exists := ctx.Get(constant.UserKeyContext)
	if exists {
		user, ok := value.(*models.User)
		if !ok || user.ID == 0 {
			return nil, constant.GetError(ctx, response.UserNotExist)
		}
		return user, nil
	}
	return nil, constant.GetError(ctx, response.UserNotExist)
}
