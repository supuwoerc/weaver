package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

type CorsMiddleware struct {
	viper *viper.Viper
}

var (
	corsMiddlewareOnce sync.Once
	corsMiddleware     *CorsMiddleware
)

// TODO:确认是否需要单例
func NewCorsMiddleware(v *viper.Viper) *CorsMiddleware {
	corsMiddlewareOnce.Do(func() {
		corsMiddleware = &CorsMiddleware{
			viper: v,
		}
	})
	return corsMiddleware
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
