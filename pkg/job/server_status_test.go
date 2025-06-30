package job

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"
	"go.uber.org/zap/zaptest"
)

func TestNewServerStatus(t *testing.T) {
	t.Run("nil logger", func(t *testing.T) {
		status := NewServerStatus(time.Second, nil)
		assert.Nil(t, status.logger)
		assert.Equal(t, status.cpuStatisticalInterval, time.Second)
	})

	t.Run("create with logger", func(t *testing.T) {
		l := logger.NewLogger(zaptest.NewLogger(t).Sugar())
		status := NewServerStatus(time.Second, l)
		assert.Equal(t, status.logger, l)
		assert.Equal(t, status.cpuStatisticalInterval, time.Second)
	})

	t.Run("create with different intervals", func(t *testing.T) {
		l := logger.NewLogger(zaptest.NewLogger(t).Sugar())

		testCases := []struct {
			name     string
			interval time.Duration
		}{
			{"1 second", time.Second},
			{"5 seconds", 5 * time.Second},
			{"1 minute", time.Minute},
			{"10 milliseconds", 10 * time.Millisecond},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				status := NewServerStatus(tc.interval, l)
				assert.Equal(t, tc.interval, status.cpuStatisticalInterval)
				assert.Equal(t, l, status.logger)
			})
		}
	})

	t.Run("create with zero interval", func(t *testing.T) {
		l := logger.NewLogger(zaptest.NewLogger(t).Sugar())
		status := NewServerStatus(0, l)
		assert.Equal(t, time.Duration(0), status.cpuStatisticalInterval)
		assert.Equal(t, l, status.logger)
	})
}

func TestServerStatus_Name(t *testing.T) {
	t.Run("name is consistent", func(t *testing.T) {
		l := logger.NewLogger(zaptest.NewLogger(t).Sugar())
		status := NewServerStatus(time.Second, l)
		// 多次调用应该返回相同的结果
		name1 := status.Name()
		name2 := status.Name()
		name3 := status.Name()
		assert.Equal(t, name1, name2)
		assert.Equal(t, name2, name3)
		assert.Equal(t, string(constant.ServerStatus), name1)
	})
}

func TestServerStatus_IfStillRunning(t *testing.T) {
	t.Run("mode is consistent", func(t *testing.T) {
		l := logger.NewLogger(zaptest.NewLogger(t).Sugar())
		status := NewServerStatus(time.Second, l)
		// 多次调用应该返回相同的结果
		mode1 := status.IfStillRunning()
		mode2 := status.IfStillRunning()
		mode3 := status.IfStillRunning()
		assert.Equal(t, mode1, mode2)
		assert.Equal(t, mode2, mode3)
		assert.Equal(t, constant.Skip, mode1)
	})
}

func TestServerStatus_Interval(t *testing.T) {
	t.Run("interval is consistent", func(t *testing.T) {
		l := logger.NewLogger(zaptest.NewLogger(t).Sugar())
		status := NewServerStatus(time.Second, l)
		// 多次调用应该返回相同的结果
		interval1 := status.Interval()
		interval2 := status.Interval()
		interval3 := status.Interval()
		assert.Equal(t, interval1, interval2)
		assert.Equal(t, interval2, interval3)
		assert.Equal(t, "0 0 * * * *", interval1)
	})
}

func TestServerStatus_Handle(t *testing.T) {
	t.Run("handle with different intervals", func(t *testing.T) {
		testCases := []struct {
			name     string
			interval time.Duration
		}{
			{"100ms", 100 * time.Millisecond},
			{"500ms", 500 * time.Millisecond},
			{"1s", time.Second},
			{"2s", 2 * time.Second},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				l := logger.NewLogger(zaptest.NewLogger(t).Sugar())
				status := NewServerStatus(tc.interval, l)
				assert.NotPanics(t, func() {
					status.Handle()
				})
			})
		}
	})

	t.Run("handle with nil logger", func(t *testing.T) {
		status := NewServerStatus(time.Second, nil)
		assert.Panics(t, func() {
			status.Handle()
		})
	})
}
