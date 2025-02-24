package router

import (
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApiRouter(r *gin.Engine) {
	group := r.Group("api/v1")
	// 国际化中间件
	group.Use(middleware.I18N(), middleware.InjectTranslator())
	// 系统基础测试
	InitPingApiRouter(group)
	// 开放api(不需要走鉴权中间件)
	InitPublicApiRouter(group)
	// 登录鉴权中间件
	group.Use(middleware.LoginRequired())
	// 用户模块
	InitUserApiRouter(group)
	// 角色模块
	InitRoleApiRouter(group)
	// 附件模块
	InitAttachmentApiRouter(group)
	// 权限模块
	InitPermissionApiRouter(group)
	// 部门模块
	InitDepartmentApiRouter(group)
}

func InitWebRouter(r *gin.Engine) {
	group := r.Group("view/v1")
	// 开放页面
	InitPublicWebRouter(group)
}

func InitSystemWebRouter(r *gin.Engine) {
	r.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusOK, "404.html", nil)
	})
}
