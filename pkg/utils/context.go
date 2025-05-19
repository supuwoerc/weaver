package utils

import (
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/jwt"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetContextClaims 从上下文中获取当前请求接口的Claims
func GetContextClaims(ctx *gin.Context) (*jwt.TokenClaims, error) {
	value, exists := ctx.Get(constant.ClaimsContextKey)
	if exists {
		claims, ok := value.(*jwt.TokenClaims)
		if !ok || claims == nil || claims.User.ID == 0 {
			return nil, response.UserNotExist
		}
		return claims, nil
	}
	return nil, response.UserNotExist
}
