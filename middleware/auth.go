package middleware

import (
	"errors"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/jwt"
	"gin-web/pkg/redis"
	"gin-web/pkg/response"
	"gin-web/repository"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
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

type AuthMiddleware struct {
	db          *gorm.DB
	redisClient *redis.CommonRedisClient
}

var (
	authMiddlewareOnce sync.Once
	authMiddleware     *AuthMiddleware
)

func NewAuthMiddleware(db *gorm.DB, redisClient *redis.CommonRedisClient) *AuthMiddleware {
	authMiddlewareOnce.Do(func() {
		authMiddleware = &AuthMiddleware{
			db:          db,
			redisClient: redisClient,
		}
	})
	return authMiddleware
}

// LoginRequired 检查token和refresh_token的有效性
func (l *AuthMiddleware) LoginRequired() gin.HandlerFunc {
	tokenKey := viper.GetString("jwt.tokenKey")
	refreshTokenKey := viper.GetString("jwt.refreshTokenKey")
	prefix := viper.GetString("jwt.tokenPrefix")
	jwtBuilder := jwt.NewJwtBuilder(l.db, l.redisClient)
	userRepository := repository.NewUserRepository(l.db, l.redisClient)
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(tokenKey)
		if token == "" || !strings.HasPrefix(token, prefix) {
			tokenInvalidResponse(ctx)
			return
		}
		token = strings.TrimPrefix(token, prefix)
		claims, err := jwtBuilder.ParseToken(token)
		if err == nil {
			// token解析正常,判断是不是在不redis中
			pair, tempErr := jwtBuilder.GetCacheToken(ctx, claims.User.Email)
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
			pair, tempErr := jwtBuilder.GetCacheToken(ctx, claims.User.Email)
			if pair == nil || tempErr != nil || pair.RefreshToken != refreshToken {
				refreshTokenInvalidResponse(ctx)
				return
			}
			newToken, newRefreshToken, refreshErr := jwtBuilder.ReGenerateTokenPairs(token, refreshToken, func(claims *jwt.TokenClaims, newToken, newRefreshToken string) error {
				return userRepository.CacheTokenPair(ctx, claims.User.Email, &models.TokenPair{
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

func (l *AuthMiddleware) PermissionRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求资源路径
		//obj := c.Request.URL.Path
		// 请求方法
		//act := c.Request.Method
		// TODO:根据用户确认请求主体
		//sub := "admin"
		//ok, err := true, nil
		//if err != nil {
		//	response.FailWithError(c, constant.GetError(c, response.CasbinErr))
		//	return
		//}
		//if ok {
		//	c.Next()
		//} else {
		//	response.FailWithError(c, constant.GetError(c, response.CasbinInvalid))
		//}
	}
}
