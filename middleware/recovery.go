package middleware

import (
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		global.Logger.Error(err, ":", string(debug.Stack()))
		response.HttpResponse[any](c, response.RecoveryError, nil, nil, nil)
	})
}
