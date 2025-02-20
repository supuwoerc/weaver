package initialize

import (
	"context"
	"errors"
	"fmt"
	"gin-web/pkg/global"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	host        string = "127.0.0.1"
	defaultPort int    = 8080
)

var (
	isLinux = false
)

// InitServer 创建http服务器
func InitServer(handle http.Handler) {
	port := viper.GetInt("server.port")
	if port == 0 {
		port = defaultPort
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handle,
	}
	if isLinux {
		graceHttpServe(srv)
	} else {
		httpServer(srv)
	}
}

func httpServer(srv *http.Server) {
	pid := os.Getpid()
	// 参考地址:https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		global.Logger.Infow("服务启动", "addr", srv.Addr, "PID", pid)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.Logger.Errorw("服务启动失败", "err", err.Error())
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeoutContext); err != nil {
		global.Logger.Errorw("服务关闭", "pid", pid, "err", err.Error())
		return
	}
	global.Logger.Infow("服务关闭", "pid", pid)
}

func graceHttpServe(srv *http.Server) {
	pid := os.Getpid()
	global.Logger.Infow("服务启动", "addr", srv.Addr, "pid", pid)
	err := gracehttp.Serve(srv)
	if err != nil {
		global.Logger.Errorw("服务启动失败", "err", err.Error())
		os.Exit(1)
	}
	global.Logger.Infow("服务关闭", "pid", pid)
}
