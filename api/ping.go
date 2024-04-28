package api

import (
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	response.SuccessWithData[string](ctx, "pong")
}

func Exception(ctx *gin.Context) {
	a := 100
	b := 100
	c := 1 / (a - b)
	response.SuccessWithData[int](ctx, c)
}
