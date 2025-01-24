package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
)

func InitPermissionRouter(r *gin.RouterGroup) {
	permissionApi := v1.NewPermissionApi()
	permissionAccessGroup := r.Group("permission")
	{
		permissionAccessGroup.POST("create", permissionApi.CreatePermission)
		permissionAccessGroup.GET("list", permissionApi.GetPermissionList)
		permissionAccessGroup.GET("detail", permissionApi.GetPermissionDetail)
		permissionAccessGroup.POST("update", permissionApi.UpdatePermission)
		permissionAccessGroup.POST("delete", permissionApi.DeletePermission)
	}
}
