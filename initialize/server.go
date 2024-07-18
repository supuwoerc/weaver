package initialize

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/spf13/viper"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	host        string = "127.0.0.1"
	defaultPort int    = 8080
)

// 创建http服务器
func InitServer(handle http.Handler) {
	// 参考地址:https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	port := viper.GetInt("server.port")
	if port == 0 {
		port = defaultPort
	}
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handle,
	}
	go func() {
		global.Logger.Info(fmt.Sprintf("服务启动，地址:%s\n", fmt.Sprintf("%s:%d", host, port)))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Error(fmt.Sprintf("服务启动失败：%s\n", err.Error()))
			return
		}
	}()
	<-ctx.Done()
	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(timeoutContext); err != nil {
		global.Logger.Error(fmt.Sprintf("服务关闭：%s\n", err.Error()))
		return
	}
	global.Logger.Info(fmt.Sprintf("服务关闭..."))
}
