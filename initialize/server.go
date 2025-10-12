package initialize

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/logger"

	"github.com/facebookgo/grace/gracehttp"
)

const (
	host        string = "0.0.0.0"
	defaultPort int    = 8080
)

var (
	isLinux = false
)

type HttpServer struct {
	httpServer *http.Server
	logger     *logger.Logger
	port       int
	conf       *conf.Config
	isLinux    bool
}

// NewHttpServer 创建http服务器
func NewHttpServer(conf *conf.Config, handle http.Handler, logger *logger.Logger) *HttpServer {
	port := wrapPort(conf.System.Port)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handle,
	}
	return &HttpServer{
		httpServer: srv,
		logger:     logger,
		port:       port,
		conf:       conf,
		isLinux:    isLinux,
	}
}

func wrapPort(port int) int {
	if port == 0 {
		port = defaultPort
	}
	return port
}

func (s *HttpServer) Port() int {
	return wrapPort(s.conf.System.Port)
}

func (s *HttpServer) Addr() string {
	return s.httpServer.Addr
}

func (s *HttpServer) Run() {
	if s.isLinux {
		s.graceRunServe(s.httpServer)
	} else {
		s.runServer(s.httpServer)
	}
}

func (s *HttpServer) runServer(srv *http.Server) {
	// 参考地址:https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		s.logger.Infow("server running", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Errorw("server err", "err", err.Error())
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeoutContext); err != nil {
		s.logger.Errorw("server closed", "err", err.Error())
		return
	}
	s.logger.Info("server closed")
}

func (s *HttpServer) graceRunServe(srv *http.Server) {
	s.logger.Infow("grace server running", "addr", srv.Addr)
	err := gracehttp.Serve(srv)
	if err != nil {
		s.logger.Errorw("grace server err", "err", err.Error())
		os.Exit(1)
	}
	s.logger.Info("grace server closed")
}
