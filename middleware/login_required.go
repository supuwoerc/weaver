package middleware

import (
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

const (
	REFRESH_URL = "/api/v1/user/refresh_token"
)

func tokenInvalidResponse(ctx *gin.Context) {
	response.HttpResponse[any](ctx, response.INVALID_TOKEN, nil, response.GetMessage(response.INVALID_TOKEN))
}

func LoginRequired() gin.HandlerFunc {
	tokenKey := viper.GetString("jwt.tokenKey")
	refreshTokenKey := viper.GetString("jwt.refreshTokenKey")
	prefix := viper.GetString("jwt.tokenPrefix")
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(tokenKey)
		if token == "" || !strings.HasPrefix(token, prefix) {
			tokenInvalidResponse(ctx)
			return
		}
		claims, err := jwt.ParseToken(token[len(prefix):])
		if err == nil {
			// token正常且未过期
			ctx.Set("user", claims)
			return
		} else if ctx.Request.URL.Path == REFRESH_URL {
			// 短token错误,检查是否满足刷新token的情况
			refreshToken := ctx.GetHeader(refreshTokenKey)
			newToken, newRefreshToken, refreshErr := jwt.ReGenerateAccessAndRefreshToken(token, refreshToken)
			if refreshErr != nil {
				return
			}
			response.SuccessWithData[response.RefreshTokenResponse](ctx, response.RefreshTokenResponse{
				Token:        newToken,
				RefreshToken: newRefreshToken,
			})
			return
		} else {
			// 其他情况直接返回错误信息
			tokenInvalidResponse(ctx)
			return
		}

	}
}
