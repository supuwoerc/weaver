package router

import (
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	group := r.Group("api/v1")
	// 国际化中间件
	group.Use(middleware.I18N(), middleware.InjectTranslator())
	// 系统基础测试
	InitPingRouter(group)
	// 开放api(不需要走鉴权中间件)
	InitPublicRouter(group)
	// 登录鉴权中间件
	group.Use(middleware.LoginRequired())
	// 用户模块
	InitUserRouter(group)
	// 角色模块
	InitRoleRouter(group)
	// 附件模块
	InitAttachmentRouter(group)
	// 权限模块
	InitPermissionRouter(group)
	// 部门模块
	InitDepartmentRouter(group)
}
