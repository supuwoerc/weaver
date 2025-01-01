package jwt

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/response"
	"gin-web/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

type TokenClaimsBasic struct {
	UID      uint
	Email    string
	Nickname string
	Roles    []uint // token中仅存储角色ID
}

type TokenClaims struct {
	jwt.RegisteredClaims
	User *TokenClaimsBasic
}

type TokenBuilder struct {
}

func NewJwtBuilder() *TokenBuilder {
	return &TokenBuilder{}
}

// 生成token
func (j *TokenBuilder) generateToken(user *TokenClaimsBasic, createAt time.Time, duration time.Duration) (string, error) {
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    viper.GetString("jwt.issuer"),
			IssuedAt:  jwt.NewNumericDate(createAt),
			ExpiresAt: jwt.NewNumericDate(createAt.Add(duration)),
		},
		User: user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("jwt.secret")))
}

// 生成短token
func (j *TokenBuilder) generateAccessToken(user *TokenClaimsBasic, createAt time.Time) (string, error) {
	return j.generateToken(user, createAt, viper.GetDuration("jwt.expires")*time.Minute)
}

// generateRefreshToken 生成长token
func (j *TokenBuilder) generateRefreshToken(createAt time.Time) (string, error) {
	return j.generateToken(&TokenClaimsBasic{}, createAt, viper.GetDuration("jwt.refreshTokenExpires")*time.Minute)
}

type RefreshTokenCallback func(claims *TokenClaims, accessToken, refreshToken string) error

// ReGenerateAccessAndRefreshToken 校验并生成长短token
func (j *TokenBuilder) ReGenerateAccessAndRefreshToken(accessToken, refreshToken string, callback RefreshTokenCallback) (string, string, error) {
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
	if claims == nil || claims.User.UID == 0 {
		return "", "", response.InvalidToken
	}
	createAt := time.Now()
	newAccessToken, err := j.generateAccessToken(&TokenClaimsBasic{
		UID:      claims.User.UID,
		Email:    claims.User.Email,
		Nickname: claims.User.Nickname,
		Roles:    claims.User.Roles,
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
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil || !token.Valid {
		return &claims, response.InvalidToken
	}
	return &claims, nil
}

// GetCacheToken 获取缓存的Token对
func (j *TokenBuilder) GetCacheToken(ctx context.Context, email string) (*models.TokenPair, error) {
	userRepository := repository.NewUserRepository()
	return userRepository.GetTokenPair(ctx, email)
}
