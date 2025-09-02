package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/jwt"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	refreshUrl = "/api/v1/user/refresh-token"
)

type AuthMiddlewareTokenRepo interface {
	CacheRefreshToken(ctx context.Context, email, refreshToken string, expiration time.Duration) error
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
			response.FailWithError(ctx, response.InvalidToken)
			return
		}
		token = strings.TrimPrefix(token, prefix)
		claims, err := l.jwtBuilder.ParseToken(token)
		isRefresh := errors.Is(err, response.InvalidToken) &&
			ctx.Request.URL.Path == refreshUrl &&
			ctx.Request.Method == http.MethodPost &&
			claims != nil &&
			claims.User != nil &&
			claims.User.Email != ""
		if err == nil {
			if ctx.Request.URL.Path == refreshUrl {
				response.FailWithError(ctx, response.UnnecessaryRefreshToken)
				return
			}
			ctx.Set(constant.ClaimsContextKey, claims)
		} else if isRefresh {
			// 短token错误,检查是否满足刷新token的情况
			refreshToken := ctx.GetHeader(refreshTokenKey)
			if strings.TrimSpace(refreshToken) == "" {
				response.FailWithError(ctx, response.InvalidRefreshToken)
				return
			}
			cachedRefreshToken, tempErr := l.jwtBuilder.GetRefreshToken(ctx, claims.User.Email)
			if cachedRefreshToken == "" || tempErr != nil || cachedRefreshToken != refreshToken {
				response.FailWithError(ctx, response.InvalidRefreshToken)
				return
			}
			newToken, refreshErr := l.jwtBuilder.GenerateAccessToken(claims.User, time.Now())
			if refreshErr != nil {
				response.FailWithError(ctx, err)
				return
			}
			response.SuccessWithData[response.RefreshTokenResponse](ctx, response.RefreshTokenResponse{
				Token: newToken,
			})
			ctx.Abort()
		} else {
			// 其他情况直接返回错误信息
			response.FailWithError(ctx, response.InvalidToken)
		}
	}
}

func (l *AuthMiddleware) PermissionRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		claims, exists := c.Get(constant.ClaimsContextKey)
		if !exists {
			response.FailWithError(c, response.AuthErr)
			return
		}
		tokenClaims, ok := claims.(*jwt.TokenClaims)
		if !ok {
			response.FailWithError(c, response.AuthErr)
			return
		}
		if tokenClaims.User.Email == l.conf.System.Admin.Email {
			c.Next()
			return
		}
		// 检查API权限
		hasPermission, err := l.permissionChecker.CheckUserPermission(c, tokenClaims.User.ID, c.Request.URL.Path, constant.ApiRoute)
		if err != nil {
			response.FailWithError(c, err)
			return
		}
		if !hasPermission {
			response.FailWithError(c, response.AuthErr)
		}
	}
}
