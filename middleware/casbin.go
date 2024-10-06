package middleware

import (
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

func PermissionRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求资源路径
		obj := c.Request.URL.Path
		// 请求方法
		act := c.Request.Method
		// TODO:根据用户确认请求主体
		sub := "admin"
		ok, err := global.CasbinEnforcer.Enforce(sub, obj, act)
		if err != nil {
			response.FailWithError(c, constant.GetError(c, response.CasbinErr))
			return
		}
		if ok {
			c.Next()
		} else {
			response.FailWithError(c, constant.GetError(c, response.CasbinInvalid))
		}
	}
}
