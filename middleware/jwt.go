package middleware

import (
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"gin-web/repository"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

const (
	refreshUrl = "/api/v1/user/refresh-token"
)

func tokenInvalidResponse(ctx *gin.Context) {
	response.FailWithError(ctx, constant.GetError(ctx, response.InvalidToken))
}

func refreshTokenInvalidResponse(ctx *gin.Context) {
	response.FailWithError(ctx, constant.GetError(ctx, response.InvalidRefreshToken))
}

func unnecessaryRefreshResponse(ctx *gin.Context) {
	response.FailWithError(ctx, constant.GetError(ctx, response.UnnecessaryRefreshToken))
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
			pair, tempErr := jwtBuilder.GetCacheToken(claims.User.Email)
			if pair == nil || tempErr != nil || pair.AccessToken != token {
				tokenInvalidResponse(ctx)
				return
			}
			if ctx.Request.URL.Path == refreshUrl && pair != nil && pair.AccessToken == token {
				unnecessaryRefreshResponse(ctx)
				return
			}
			ctx.Set(constant.UserKeyContext, claims)
		} else if ctx.Request.URL.Path == refreshUrl && ctx.Request.Method == http.MethodGet {
			// 短token错误,检查是否满足刷新token的情况
			refreshToken := ctx.GetHeader(refreshTokenKey)
			if strings.TrimSpace(refreshToken) == "" {
				refreshTokenInvalidResponse(ctx)
				return
			}
			pair, tempErr := jwtBuilder.GetCacheToken(claims.User.Email)
			if pair == nil || tempErr != nil || pair.RefreshToken != refreshToken {
				refreshTokenInvalidResponse(ctx)
				return
			}
			newToken, newRefreshToken, refreshErr := jwtBuilder.ReGenerateAccessAndRefreshToken(token, refreshToken, func(claims *jwt.TokenClaims, newToken, newRefreshToken string) error {
				return userRepository.CacheTokenPair(ctx, claims.User.Email, &models.TokenPair{
					AccessToken:  newToken,
					RefreshToken: newRefreshToken,
				})
			})
			if refreshErr == constant.GetError(ctx, response.InvalidToken) {
				tokenInvalidResponse(ctx)
				return
			}
			if refreshErr == constant.GetError(ctx, response.InvalidRefreshToken) {
				refreshTokenInvalidResponse(ctx)
				return
			}
			if refreshErr != nil {
				response.FailWithError(ctx, err)
				return
			}
			response.SuccessWithData[response.RefreshTokenResponse](ctx, response.RefreshTokenResponse{
				Token:        newToken,
				RefreshToken: newRefreshToken,
			})
			ctx.Abort()
		} else {
			// 其他情况直接返回错误信息
			tokenInvalidResponse(ctx)
		}
	}
}
