package api

import (
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	response.SuccessWithData[string](ctx, "pong")
}
