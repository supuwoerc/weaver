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
	ID     uint
	Name   string
	Gender models.UserGender
}

var (
	TOKEN_SECRET          = viper.GetString("jwt.secret")
	TOKEN_ISSUER          = viper.GetString("jwt.secret")
	TOKEN_EXPIRES         = viper.GetDuration("jwt.expires") * time.Minute
	REFRESH_TOKEN_EXPIRES = viper.GetDuration("jwt.refreshTokenExpires") * time.Minute
)

// 生成token
func generateToken(id uint, name string, gender models.UserGender, createAt time.Time, duration time.Duration) (string, error) {
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TOKEN_ISSUER,
			IssuedAt:  jwt.NewNumericDate(createAt),
			ExpiresAt: jwt.NewNumericDate(createAt.Add(duration)),
		},
		ID:     id,
		Name:   name,
		Gender: gender,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(TOKEN_SECRET))
}

// 生成短token
func generateAccessToken(id uint, name string, gender models.UserGender, createAt time.Time) (string, error) {
	return generateToken(id, name, gender, createAt, TOKEN_EXPIRES)
}

// 生成长token
func generateRefreshToken(createAt time.Time) (string, error) {
	return generateToken(0, "", 0, createAt, REFRESH_TOKEN_EXPIRES)
}

// 校验并生成长短token
func ReGenerateAccessAndRefreshToken(accessToken, refreshToken string) (string, string, error) {
	if _, err := ParseToken(refreshToken); err != nil {
		return "", "", constant.REFRESH_TOKEN_PARSE_ERROR
	}
	claims, err := ParseToken(accessToken)
	if err == nil {
		return "", "", constant.UNNECESSARY_REFRESH_TOKEN_ERROR
	}
	if err != constant.TOKEN_PARSE_ERROR {
		return "", "", err
	}
	createAt := time.Now()
	newAccessToken, err := generateAccessToken(claims.ID, claims.Name, claims.Gender, createAt)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := generateRefreshToken(createAt)
	if err != nil {
		return "", "", err
	}
	// TODO：将旧的token和refresh_token加入黑名单
	return newAccessToken, newRefreshToken, nil
}

// 生成长短token
func GenerateAccessAndRefreshToken(id uint, name string, gender models.UserGender) (string, string, error) {
	createAt := time.Now()
	newAccessToken, err := generateAccessToken(id, name, gender, createAt)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := generateRefreshToken(createAt)
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}

// 解析token
func ParseToken(tokenString string) (TokenClaims, error) {
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
