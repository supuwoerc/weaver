package api

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Ping(ctx *gin.Context) {
	response.SuccessWithData[string](ctx, "pong")
}

func Exception(ctx *gin.Context) {
	num := 100 - 100
	response.SuccessWithData[int](ctx, 1/num)
}

func CheckPermission(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}

func SlowResponse(ctx *gin.Context) {
	value := ctx.Query("t")
	second, err := strconv.Atoi(value)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
	}
	time.Sleep(time.Duration(second) * time.Second)
	ctx.String(http.StatusOK, fmt.Sprintf("sleep %ds,PID %d", second, os.Getpid()))
}

func LockResponse(ctx *gin.Context) {
	lock := utils.NewRedisLock("test", time.Second*3)
	if err := utils.TryLock(ctx, lock, true); err != nil {
		response.FailWithError(ctx, err)
		return
	}
	global.Logger.Info("lock success")
	defer func(lock *utils.RedisLock) {
		e := utils.Unlock(lock)
		if e != nil {
			global.Logger.Infof("unlock fail %s", e.Error())
			return
		}
		global.Logger.Info("unlock success")
	}(lock)
	time.Sleep(time.Second * 20)
	response.Success(ctx)
}
