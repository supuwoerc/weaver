package middleware

import (
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

func tokenInvalidResponse(ctx *gin.Context) {
	response.HttpResponse[any](ctx, response.INVALID_TOKEN, nil, response.GetMessage(response.INVALID_TOKEN))
}

func LoginRequired() gin.HandlerFunc {
	// TODO:长短token实现
	tokenKey := viper.GetString("jwt.tokenKey")
	prefix := viper.GetString("jwt.tokenPrefix")
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(tokenKey)
		if token == "" || strings.HasPrefix(token, prefix) {
			tokenInvalidResponse(ctx)
			return
		}
		claims, err := jwt.ParseToken(token[len(prefix):])
		if err != nil {
			tokenInvalidResponse(ctx)
			return
		}
		ctx.Set("user", claims)
	}
}
