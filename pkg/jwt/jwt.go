package jwt

import (
	"context"
	"time"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/redis"
	"github.com/supuwoerc/weaver/pkg/response"

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
	GetRefreshToken(ctx context.Context, email string) (string, error)
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

// generateToken 生成token
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

// GenerateAccessToken 生成短token
func (j *TokenBuilder) GenerateAccessToken(user *TokenClaimsBasic, createAt time.Time) (string, error) {
	return j.generateToken(user, createAt, j.getAccessTokenExpiration())
}

// generateRefreshToken 生成长token
func (j *TokenBuilder) generateRefreshToken(createAt time.Time) (string, error) {
	return j.generateToken(&TokenClaimsBasic{}, createAt, j.GetRefreshTokenExpiration())
}

func (j *TokenBuilder) getAccessTokenExpiration() time.Duration {
	return j.conf.JWT.Expires * time.Minute
}

func (j *TokenBuilder) GetRefreshTokenExpiration() time.Duration {
	return j.conf.JWT.RefreshTokenExpires * time.Minute
}

// GenerateAccessAndRefreshToken 生成长短token
func (j *TokenBuilder) GenerateAccessAndRefreshToken(user *TokenClaimsBasic) (string, string, error) {
	createAt := time.Now()
	newAccessToken, err := j.GenerateAccessToken(user, createAt)
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
	var claims TokenClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.conf.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return &claims, response.InvalidToken
	}
	return &claims, nil
}

// GetRefreshToken 获取缓存的RefreshToken
func (j *TokenBuilder) GetRefreshToken(ctx context.Context, email string) (string, error) {
	return j.repo.GetRefreshToken(ctx, email)
}
