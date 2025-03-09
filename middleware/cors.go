package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

type CorsMiddleware struct {
	viper *viper.Viper
}

func NewCorsMiddleware(v *viper.Viper) *CorsMiddleware {
	return &CorsMiddleware{
		viper: v,
	}
}

func (c *CorsMiddleware) Cors() gin.HandlerFunc {
	originPrefix := c.viper.GetStringSlice("cors.originPrefix")
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
		AllowHeaders:     []string{"Content-Type", "Authorization", "Locale", "Refresh-Token"},
		AllowCredentials: true,
	})
}
