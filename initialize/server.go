package initialize

import (
	"context"
	"errors"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

type HttpServer struct {
	httpServer *http.Server
	logger     *zap.SugaredLogger
	isLinux    bool
}

// NewServer 创建http服务器
func NewServer(v *viper.Viper, handle http.Handler, logger *zap.SugaredLogger) *HttpServer {
	port := v.GetInt("server.port")
	if port == 0 {
		port = defaultPort
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handle,
	}
	return &HttpServer{
		httpServer: srv,
		logger:     logger,
		isLinux:    isLinux,
	}
}

func (s *HttpServer) Run() {
	if s.isLinux {
		s.graceRunServe(s.httpServer, s.logger)
	} else {
		s.runServer(s.httpServer, s.logger)
	}
}

func (s *HttpServer) runServer(srv *http.Server, logger *zap.SugaredLogger) {
	pid := os.Getpid()
	// 参考地址:https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		logger.Infow("服务启动", "addr", srv.Addr, "pid", pid)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorw("服务启动失败", "err", err.Error())
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeoutContext); err != nil {
		logger.Errorw("服务关闭", "pid", pid, "err", err.Error())
		return
	}
	logger.Infow("服务关闭", "pid", pid)
}

func (s *HttpServer) graceRunServe(srv *http.Server, logger *zap.SugaredLogger) {
	pid := os.Getpid()
	logger.Infow("服务启动", "addr", srv.Addr, "pid", pid)
	err := gracehttp.Serve(srv)
	if err != nil {
		logger.Errorw("服务启动失败", "err", err.Error())
		os.Exit(1)
	}
	logger.Infow("服务关闭", "pid", pid)
}
