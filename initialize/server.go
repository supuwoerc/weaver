package initialize

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	HOST string = "127.0.0.1"
	PORT int    = 8080
)

// 创建http服务器
func InitServer(handle http.Handler) {
	// 参考地址:https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", HOST, PORT),
		Handler: handle,
	}
	go func() {
		//TODO：记录日志
		fmt.Printf("服务启动，地址:%s\n", fmt.Sprintf("%s:%d", HOST, PORT))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			//TODO：记录日志
			fmt.Printf("服务启动失败：%s\n", err.Error())
			return
		}
	}()
	<-ctx.Done()
	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(timeoutContext); err != nil {
		//TODO：记录日志
		fmt.Printf("服务关闭：%s\n", err.Error())
		return
	}
	//TODO：记录日志
	fmt.Println("服务关闭...")
}
