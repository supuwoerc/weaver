package utils

import (
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetContextUser 从上下文中获取当前请求接口的用户
func GetContextUser(ctx *gin.Context) (*models.User, error) {
	value, exists := ctx.Get(constant.UserKeyContext)
	if exists {
		user, ok := value.(*models.User)
		if !ok || user.ID == 0 {
			return nil, response.UserNotExist
		}
		return user, nil
	}
	return nil, response.UserNotExist
}

// GetContextClaims 从上下文中获取当前请求接口的Claims
func GetContextClaims(ctx *gin.Context) (*jwt.TokenClaims, error) {
	value, exists := ctx.Get(constant.ClaimsKeyContext)
	if exists {
		claims, ok := value.(*jwt.TokenClaims)
		if !ok || claims == nil || claims.User.UID == 0 {
			return nil, response.UserNotExist
		}
		return claims, nil
	}
	return nil, response.UserNotExist
}
