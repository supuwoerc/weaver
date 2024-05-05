package api

import (
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	response.SuccessWithData[string](ctx, "pong")
}

func Exception(ctx *gin.Context) {
	response.SuccessWithData[any](ctx, nil)
}
