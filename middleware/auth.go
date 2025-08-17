package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/jwt"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	refreshUrl = "/api/v1/user/refresh-token"
)

func tokenInvalidResponse(ctx *gin.Context) {
	response.FailWithError(ctx, response.InvalidToken)
}

func refreshTokenInvalidResponse(ctx *gin.Context) {
	response.FailWithError(ctx, response.InvalidRefreshToken)
}

func unnecessaryRefreshResponse(ctx *gin.Context) {
	response.FailWithError(ctx, response.UnnecessaryRefreshToken)
}

type AuthMiddlewareTokenRepo interface {
	CacheTokenPair(ctx context.Context, email string, pair *models.TokenPair) error
}

type AuthMiddlewarePermissionRepo interface {
	CheckUserPermission(ctx context.Context, uid uint, resource string, permissionType constant.PermissionType) (bool, error)
}

type AuthMiddleware struct {
	conf              *conf.Config
	tokenRepo         AuthMiddlewareTokenRepo
	jwtBuilder        *jwt.TokenBuilder
	permissionChecker AuthMiddlewarePermissionRepo
}

func NewAuthMiddleware(
	conf *conf.Config,
	tokenRepo AuthMiddlewareTokenRepo,
	jwtBuilder *jwt.TokenBuilder,
	permissionCheck AuthMiddlewarePermissionRepo,
) *AuthMiddleware {
	return &AuthMiddleware{
		conf:              conf,
		tokenRepo:         tokenRepo,
		jwtBuilder:        jwtBuilder,
		permissionChecker: permissionCheck,
	}
}

// LoginRequired 检查token和refresh_token的有效性
func (l *AuthMiddleware) LoginRequired() gin.HandlerFunc {
	tokenKey := l.conf.JWT.TokenKey
	refreshTokenKey := l.conf.JWT.RefreshTokenKey
	prefix := l.conf.JWT.TokenPrefix
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(tokenKey)
		if token == "" || !strings.HasPrefix(token, prefix) {
			tokenInvalidResponse(ctx)
			return
		}
		token = strings.TrimPrefix(token, prefix)
		claims, err := l.jwtBuilder.ParseToken(token)
		if err == nil {
			// token解析正常,判断是不是在不redis中
			pair, tempErr := l.jwtBuilder.GetCacheToken(ctx, claims.User.Email)
			if pair == nil || tempErr != nil || pair.AccessToken != token {
				tokenInvalidResponse(ctx)
				return
			}
			if ctx.Request.URL.Path == refreshUrl && pair != nil && pair.AccessToken == token {
				unnecessaryRefreshResponse(ctx)
				return
			}
			ctx.Set(constant.ClaimsContextKey, claims)
		} else if ctx.Request.URL.Path == refreshUrl && ctx.Request.Method == http.MethodGet {
			// 短token错误,检查是否满足刷新token的情况
			refreshToken := ctx.GetHeader(refreshTokenKey)
			if strings.TrimSpace(refreshToken) == "" {
				refreshTokenInvalidResponse(ctx)
				return
			}
			pair, tempErr := l.jwtBuilder.GetCacheToken(ctx, claims.User.Email)
			if pair == nil || tempErr != nil || pair.RefreshToken != refreshToken {
				refreshTokenInvalidResponse(ctx)
				return
			}
			newToken, newRefreshToken, refreshErr := l.jwtBuilder.ReGenerateTokenPairs(token, refreshToken, func(claims *jwt.TokenClaims, newToken, newRefreshToken string) error {
				return l.tokenRepo.CacheTokenPair(ctx, claims.User.Email, &models.TokenPair{
					AccessToken:  newToken,
					RefreshToken: newRefreshToken,
				})
			})
			if errors.Is(refreshErr, response.InvalidToken) {
				tokenInvalidResponse(ctx)
				return
			}
			if errors.Is(refreshErr, response.InvalidRefreshToken) {
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

func permissionInvalidResponse(ctx *gin.Context) {
	response.FailWithError(ctx, response.AuthErr)
}

func (l *AuthMiddleware) PermissionRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		claims, exists := c.Get(constant.ClaimsContextKey)
		if !exists {
			permissionInvalidResponse(c)
			return
		}
		tokenClaims, ok := claims.(*jwt.TokenClaims)
		if !ok {
			permissionInvalidResponse(c)
			return
		}
		// 检查API权限
		hasPermission, err := l.permissionChecker.CheckUserPermission(c, tokenClaims.User.ID, c.Request.URL.Path, constant.ApiRoute)
		if err != nil {
			response.FailWithError(c, err)
			return
		}
		if !hasPermission {
			permissionInvalidResponse(c)
		}
	}
}
