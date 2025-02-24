package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRoleApiRouter(r *gin.RouterGroup) {
	roleApi := v1.NewRoleApi()
	roleAccessGroup := r.Group("role")
	{
		roleAccessGroup.POST("create", roleApi.CreateRole)
		roleAccessGroup.GET("list", roleApi.GetRoleList)
		roleAccessGroup.GET("detail", roleApi.GetRoleDetail)
		roleAccessGroup.POST("update", roleApi.UpdateRole)
		roleAccessGroup.POST("delete", roleApi.DeleteRole)
	}
}
