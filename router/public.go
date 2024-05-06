package router

import "github.com/gin-gonic/gin"

func InitPublicRouter(r *gin.RouterGroup) *gin.RouterGroup {
	group := r.Group("public")
	{

	}
	return group
}
