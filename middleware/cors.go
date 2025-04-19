package middleware

import (
	"gin-web/conf"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

type CorsMiddleware struct {
	conf *conf.Config
}

func NewCorsMiddleware(conf *conf.Config) *CorsMiddleware {
	return &CorsMiddleware{
		conf: conf,
	}
}

func (c *CorsMiddleware) Cors() gin.HandlerFunc {
	originPrefix := c.conf.Cors.OriginPrefix
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
