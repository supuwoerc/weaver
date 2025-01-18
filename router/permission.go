package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
)

func InitPermissionRouter(r *gin.RouterGroup) {
	permissionApi := v1.NewPermissionApi()
	permissionAccessGroup := r.Group("permission")
	{
		// TODO:添加权限限制
		permissionAccessGroup.POST("create", permissionApi.CreatePermission)
		permissionAccessGroup.GET("list", permissionApi.GetPermissionList)
		permissionAccessGroup.GET("detail", permissionApi.GetPermissionDetail)
	}
}
