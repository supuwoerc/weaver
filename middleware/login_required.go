package middleware

import (
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"strings"
)

func tokenInvalidResponse(ctx *gin.Context) {
	response.HttpResponse[any](ctx, response.INVALID_TOKEN, nil, response.GetMessage(response.INVALID_TOKEN))
}

func LoginRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO:从配置文件中读取
		token := ctx.GetHeader("Authorization")
		// TODO:从配置文件中读取
		var prefix = "Bearer "
		if token == "" || strings.HasPrefix(token, prefix) {
			tokenInvalidResponse(ctx)
			return
		}
		claims, err := jwt.ParseToken(token[len(prefix):])
		if err != nil {
			tokenInvalidResponse(ctx)
			return
		}
		// TODO:从配置文件中读取key
		ctx.Set("user", claims)
	}
}
