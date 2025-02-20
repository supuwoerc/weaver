package job

import (
	"fmt"
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
	"runtime"
)

type ServerStatus struct{}

func NewServerStatus() *ServerStatus {
	return &ServerStatus{}
}

func (s *ServerStatus) Name() string {
	return constant.ServerStatus.String()
}

func (s *ServerStatus) Handle() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	args := []any{
		"CPU Number", runtime.NumCPU(),
		"Goroutine Number", runtime.NumGoroutine(),
		"System Memory", bytes2MB(memStats.Sys),
		"Alloc Memory", bytes2MB(memStats.Alloc),
		"Malloc Memory", bytes2MB(memStats.Mallocs),
		"OS", runtime.GOOS,
		"Architecture", runtime.GOARCH,
	}
	global.Logger.Infow("server status", args...)
}

func bytes2MB(kbs uint64) string {
	return fmt.Sprintf("%.2fMB", float64(kbs)/(1024))
}
