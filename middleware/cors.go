package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	INCLUDE_PREFIX []string = []string{"http://localhost", "http://127.0.0.1"}
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// TODO:只允许本地和特定域名，读取配置来设置domain
			for _, val := range INCLUDE_PREFIX {
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
