package jwt

import (
	"context"
	"errors"
	"gin-web/conf"
	"gin-web/models"
	"gin-web/pkg/redis"
	"gin-web/pkg/response"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type TokenClaimsBasic struct {
	ID       uint    `json:"uid"`
	Email    string  `json:"email"`
	Nickname *string `json:"nickname"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	User *TokenClaimsBasic
}

type TokenBuilderRepo interface {
	GetTokenPair(ctx context.Context, email string) (*models.TokenPair, error)
}
type TokenBuilder struct {
	db          *gorm.DB
	redisClient *redis.CommonRedisClient
	conf        *conf.Config
	repo        TokenBuilderRepo
}

func NewJwtBuilder(db *gorm.DB, r *redis.CommonRedisClient, conf *conf.Config, repo TokenBuilderRepo) *TokenBuilder {
	return &TokenBuilder{
		db:          db,
		redisClient: r,
		conf:        conf,
		repo:        repo,
	}
}

// 生成token
func (j *TokenBuilder) generateToken(user *TokenClaimsBasic, createAt time.Time, duration time.Duration) (string, error) {
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.conf.JWT.Issuer,
			IssuedAt:  jwt.NewNumericDate(createAt),
			ExpiresAt: jwt.NewNumericDate(createAt.Add(duration)),
		},
		User: user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.conf.JWT.Secret))
}

// 生成短token
func (j *TokenBuilder) generateAccessToken(user *TokenClaimsBasic, createAt time.Time) (string, error) {
	return j.generateToken(user, createAt, j.conf.JWT.Expires*time.Minute)
}

// generateRefreshToken 生成长token
func (j *TokenBuilder) generateRefreshToken(createAt time.Time) (string, error) {
	return j.generateToken(&TokenClaimsBasic{}, createAt, j.conf.JWT.RefreshTokenExpires*time.Minute)
}

type RefreshTokenCallback func(claims *TokenClaims, accessToken, refreshToken string) error

// ReGenerateTokenPairs 校验并生成长短token
func (j *TokenBuilder) ReGenerateTokenPairs(accessToken, refreshToken string, callback RefreshTokenCallback) (string, string, error) {
	if _, err := j.ParseToken(refreshToken); err != nil {
		return "", "", response.InvalidRefreshToken
	}
	claims, err := j.ParseToken(accessToken)
	if err == nil {
		return "", "", response.UnnecessaryRefreshToken
	}
	if !errors.Is(response.InvalidToken, err) {
		return "", "", err
	}
	if claims == nil || claims.User.ID == 0 {
		return "", "", response.InvalidToken
	}
	createAt := time.Now()
	newAccessToken, err := j.generateAccessToken(&TokenClaimsBasic{
		ID:       claims.User.ID,
		Email:    claims.User.Email,
		Nickname: claims.User.Nickname,
	}, createAt)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := j.generateRefreshToken(createAt)
	if err != nil {
		return "", "", err
	}
	if callback != nil {
		callbackErr := callback(claims, newAccessToken, newRefreshToken)
		if callbackErr != nil {
			return "", "", callbackErr
		}
	}
	return newAccessToken, newRefreshToken, nil
}

// GenerateAccessAndRefreshToken 生成长短token
func (j *TokenBuilder) GenerateAccessAndRefreshToken(user *TokenClaimsBasic) (string, string, error) {
	createAt := time.Now()
	newAccessToken, err := j.generateAccessToken(user, createAt)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := j.generateRefreshToken(createAt)
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}

// ParseToken 解析token
func (j *TokenBuilder) ParseToken(tokenString string) (*TokenClaims, error) {
	claims := TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.conf.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return &claims, response.InvalidToken
	}
	return &claims, nil
}

// GetCacheToken 获取缓存的Token对
func (j *TokenBuilder) GetCacheToken(ctx context.Context, email string) (*models.TokenPair, error) {
	return j.repo.GetTokenPair(ctx, email)
}
