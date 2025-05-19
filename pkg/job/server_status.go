package job

import (
	"fmt"
	"runtime"
	"time"

	"github.com/supuwoerc/weaver/pkg/constant"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

type ServerStatus struct {
	cpuStatisticalInterval time.Duration
	logger                 *zap.SugaredLogger
}

func NewServerStatus(t time.Duration, logger *zap.SugaredLogger) *ServerStatus {
	return &ServerStatus{
		cpuStatisticalInterval: t,
		logger:                 logger,
	}
}

func (s *ServerStatus) Name() string {
	return string(constant.ServerStatus)
}
func (s *ServerStatus) IfStillRunning() constant.JobStillMode {
	return constant.Skip
}
func (s *ServerStatus) Interval() string {
	return "0 0 * * * *"
}

func (s *ServerStatus) Handle() {
	memory, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	percent, err := cpu.Percent(s.cpuStatisticalInterval, false)
	if err != nil {
		panic(err)
	}
	args := []any{
		"CPU Number", runtime.NumCPU(),
		"Goroutine Number", runtime.NumGoroutine(),
		"OS", runtime.GOOS,
		"Architecture", runtime.GOARCH,
		"Total Memory", bytes2MB(memory.Total),
		"Available Memory", bytes2MB(memory.Available),
		"Memory UsedPercent", float2Percent(memory.UsedPercent),
		"CPU UsedPercent", float2Percent(percent[0]),
	}
	s.logger.Infow("server status", args...)
}

func bytes2MB(kbs uint64) string {
	return fmt.Sprintf("%.2fMB", float64(kbs)/(1024*1024))
}

func float2Percent(f float64) string {
	return fmt.Sprintf("%.2f%%", f)
}
