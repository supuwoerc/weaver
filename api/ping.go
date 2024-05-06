package api

import (
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	response.SuccessWithData[string](ctx, "pong")
}

func Exception(ctx *gin.Context) {
	num := 100 - 100
	response.SuccessWithData[int](ctx, 1/num)
}
