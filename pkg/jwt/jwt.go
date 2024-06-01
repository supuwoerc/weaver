package jwt

import (
	"gin-web/models"
	"gin-web/pkg/constant"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	ID       uint
	Email    string
	NickName string
	Gender   models.UserGender
	About    string
	Birthday string
}

var (
	TOKEN_SECRET          = viper.GetString("jwt.secret")
	TOKEN_ISSUER          = viper.GetString("jwt.secret")
	TOKEN_EXPIRES         = viper.GetDuration("jwt.expires") * time.Minute
	REFRESH_TOKEN_EXPIRES = viper.GetDuration("jwt.refreshTokenExpires") * time.Minute
)

type JwtBuilder struct {
}

var jwtBuilder *JwtBuilder

func NewJwtBuilder() *JwtBuilder {
	if jwtBuilder == nil {
		jwtBuilder = &JwtBuilder{}
	}
	return jwtBuilder
}

// 生成token
func (j *JwtBuilder) generateToken(id uint, name string, gender models.UserGender, createAt time.Time, duration time.Duration) (string, error) {
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TOKEN_ISSUER,
			IssuedAt:  jwt.NewNumericDate(createAt),
			ExpiresAt: jwt.NewNumericDate(createAt.Add(duration)),
		},
		ID:       id,
		NickName: name,
		Gender:   gender,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(TOKEN_SECRET))
}

// 生成短token
func (j *JwtBuilder) generateAccessToken(id uint, name string, gender models.UserGender, createAt time.Time) (string, error) {
	return j.generateToken(id, name, gender, createAt, TOKEN_EXPIRES)
}

// 生成长token
func (j *JwtBuilder) generateRefreshToken(createAt time.Time) (string, error) {
	return j.generateToken(0, "", 0, createAt, REFRESH_TOKEN_EXPIRES)
}

type RefreshTokenCallback func(claims TokenClaims) error

// 校验并生成长短token
func (j *JwtBuilder) ReGenerateAccessAndRefreshToken(accessToken, refreshToken string, callback RefreshTokenCallback) (string, string, error) {
	if _, err := j.ParseToken(refreshToken); err != nil {
		return "", "", constant.REFRESH_TOKEN_PARSE_ERROR
	}
	claims, err := j.ParseToken(accessToken)
	if err == nil {
		return "", "", constant.UNNECESSARY_REFRESH_TOKEN_ERROR
	}
	if err != constant.TOKEN_PARSE_ERROR {
		return "", "", err
	}
	createAt := time.Now()
	newAccessToken, err := j.generateAccessToken(claims.ID, claims.NickName, claims.Gender, createAt)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := j.generateRefreshToken(createAt)
	if err != nil {
		return "", "", err
	}
	if callback != nil {
		callbackErr := callback(claims)
		if callbackErr != nil {
			return "", "", callbackErr
		}
	}
	return newAccessToken, newRefreshToken, nil
}

// 生成长短token
func (j *JwtBuilder) GenerateAccessAndRefreshToken(id uint, name string, gender models.UserGender) (string, string, error) {
	createAt := time.Now()
	newAccessToken, err := j.generateAccessToken(id, name, gender, createAt)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := j.generateRefreshToken(createAt)
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}

// 解析token
func (j *JwtBuilder) ParseToken(tokenString string) (TokenClaims, error) {
	claims := TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return TOKEN_SECRET, nil
	})
	if err != nil {
		return claims, err
	}
	if !token.Valid {
		return TokenClaims{}, constant.TOKEN_PARSE_ERROR
	}
	return claims, nil
}
