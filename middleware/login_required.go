package middleware

import (
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"gin-web/repository"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

const (
	REFRESH_URL = "/api/v1/user/refresh_token"
)

func tokenInvalidResponse(ctx *gin.Context) {
	response.HttpResponse[any](ctx, response.INVALID_TOKEN, nil, nil, nil)
}

func refreshTokenInvalidResponse(ctx *gin.Context) {
	response.HttpResponse[any](ctx, response.INVALID_REFRESH_TOKEN, nil, nil, nil)
}

func LoginRequired() gin.HandlerFunc {
	tokenKey := viper.GetString("jwt.tokenKey")
	refreshTokenKey := viper.GetString("jwt.refreshTokenKey")
	prefix := viper.GetString("jwt.tokenPrefix")
	return func(ctx *gin.Context) {
		jwtBuilder := jwt.NewJwtBuilder()
		token := ctx.GetHeader(tokenKey)
		if token == "" || !strings.HasPrefix(token, prefix) {
			tokenInvalidResponse(ctx)
			return
		}
		claims, err := jwtBuilder.ParseToken(token[len(prefix):])
		userRepository := repository.NewUserRepository()
		if err == nil {
			// token解析正常,判断是不是在不redis中
			exist, existErr := userRepository.TokenPairExist(ctx, claims.User.Email)
			if existErr != nil || !exist {
				tokenInvalidResponse(ctx)
				return
			}
			ctx.Set("user", claims)
			return
		} else if ctx.Request.URL.Path == REFRESH_URL {
			// 短token错误,检查是否满足刷新token的情况
			refreshToken := ctx.GetHeader(refreshTokenKey)
			newToken, newRefreshToken, refreshErr := jwtBuilder.ReGenerateAccessAndRefreshToken(token, refreshToken, func(claims jwt.TokenClaims) error {
				return userRepository.DelTokenPair(ctx, claims.User.Email)
			})
			if refreshErr != nil {
				refreshTokenInvalidResponse(ctx)
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
