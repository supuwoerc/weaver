package utils

import (
	"github.com/pkg/errors"
	"sync"
	"time"
)

type Watchdog struct {
	interval          time.Duration                  // 检测间隔
	timeout           time.Duration                  // 最大心跳间隔
	restartCmd        func() error                   // 重启方法
	restartFail       func(dog *Watchdog, err error) // 重启失败回调
	restartInProgress bool                           // 是否正在重启
	lastHeartBeat     time.Time                      // 最后一次心跳
	mu                sync.Mutex                     // 锁(保护lastHeartBeat的访问安全)
	stopChan          chan struct{}                  // 停止监听的channel
}

// NewWatchdog 创建一个新的看门狗
func NewWatchdog(interval, timeout time.Duration, restartCmd func() error, restartFail func(dog *Watchdog, err error)) *Watchdog {
	return &Watchdog{
		interval:          interval,
		timeout:           timeout,
		restartCmd:        restartCmd,
		restartFail:       restartFail,
		restartInProgress: false,
		lastHeartBeat:     time.Now(),
		stopChan:          make(chan struct{}),
	}
}

// Start 开始监听
func (w *Watchdog) Start() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.stopChan:
			return
		case <-ticker.C:
			w.mu.Lock()
			if w.stopChan == nil {
				w.mu.Unlock()
				return
			}
			if w.restartInProgress {
				w.mu.Unlock()
				continue
			}
			if time.Since(w.lastHeartBeat) > w.timeout {
				// 异步调用重启
				w.restartInProgress = true
				w.mu.Unlock()
				go func() {
					defer func() {
						if r := recover(); r != nil {
							w.restartFail(w, errors.Errorf("Watchdog recovered from panic: %v", r))
							w.Stop()
						}
					}()
					err := w.restartCmd()
					w.mu.Lock()
					w.restartInProgress = false
					if err != nil {
						w.mu.Unlock()
						w.restartFail(w, err)
					} else {
						w.lastHeartBeat = time.Now()
						w.mu.Unlock()
					}
				}()
			} else {
				w.mu.Unlock()
			}
		}
	}
}

// Stop 停止监听
func (w *Watchdog) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopChan != nil {
		close(w.stopChan)
		w.stopChan = nil
	}
}

// HeartBeat 心跳
func (w *Watchdog) HeartBeat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastHeartBeat = time.Now()
}
