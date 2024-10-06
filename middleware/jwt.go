package middleware

import (
	"gin-web/pkg/constant"
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"gin-web/repository"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

const (
	refreshUrl = "/api/v1/user/refresh_token"
)

func tokenInvalidResponse(ctx *gin.Context) {
	response.FailWithError(ctx, constant.GetError(ctx, response.InvalidToken))
}

func refreshTokenInvalidResponse(ctx *gin.Context) {
	response.FailWithError(ctx, constant.GetError(ctx, response.InvalidRefreshToken))
}

func LoginRequired() gin.HandlerFunc {
	tokenKey := viper.GetString("jwt.tokenKey")
	refreshTokenKey := viper.GetString("jwt.refreshTokenKey")
	prefix := viper.GetString("jwt.tokenPrefix")
	return func(ctx *gin.Context) {
		jwtBuilder := jwt.NewJwtBuilder(ctx)
		token := ctx.GetHeader(tokenKey)
		if token == "" || !strings.HasPrefix(token, prefix) {
			tokenInvalidResponse(ctx)
			return
		}
		token = strings.TrimPrefix(token, prefix)
		claims, err := jwtBuilder.ParseToken(token)
		userRepository := repository.NewUserRepository(ctx)
		if err == nil {
			// token解析正常,判断是不是在不redis中
			exist, existErr := userRepository.TokenPairExist(ctx, claims.User.Email)
			if existErr != nil || !exist {
				tokenInvalidResponse(ctx)
				return
			}
			ctx.Set(constant.ClaimKeyContext, claims)
			return
		} else if ctx.Request.URL.Path == refreshUrl {
			// 短token错误,检查是否满足刷新token的情况
			refreshToken := ctx.GetHeader(refreshTokenKey)
			newToken, newRefreshToken, refreshErr := jwtBuilder.ReGenerateAccessAndRefreshToken(token, refreshToken, func(claims *jwt.TokenClaims) error {
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
