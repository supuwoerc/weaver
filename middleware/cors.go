package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

func Cors() gin.HandlerFunc {
	originPrefix := viper.GetStringSlice("cors.originPrefix")
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			for _, val := range originPrefix {
				if strings.HasPrefix(origin, val) {
					return true
				}
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Locale"},
		AllowCredentials: true,
	})
}
