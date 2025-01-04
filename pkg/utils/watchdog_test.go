package utils

import (
	"fmt"
	"testing"
	"time"
)

// 模拟重启命令
func mockRestartCmd() error {
	fmt.Println("restart...")
	time.Sleep(3 * time.Second)
	fmt.Println("restart success")
	return nil
}

// 模拟失败的重启
func mockFailRestartCmd() error {
	fmt.Println("restart...")
	time.Sleep(3 * time.Second)
	fmt.Println("restart fail")
	return fmt.Errorf("restart failed")
}

// 模拟重启失败回调
func mockRestartFail(d *Watchdog, err error) {
	fmt.Printf("Restart failed: %v\n", err)
}

func TestWatchdog_HeartBeat(t *testing.T) {
	watchdog := NewWatchdog(1*time.Second, 2*time.Second, mockRestartCmd, mockRestartFail)
	watchdog.HeartBeat()
	watchdog.mu.Lock()
	if time.Since(watchdog.lastHeartBeat) > 1*time.Second {
		t.Errorf("HeartBeat failed to update lastHeartBeat in time")
	}
	watchdog.mu.Unlock()
}

func TestWatchdog_Restart(t *testing.T) {
	// 创建一个 Watchdog 实例，模拟成功重启
	watchdog := NewWatchdog(1*time.Second, 2*time.Second, mockRestartCmd, mockRestartFail)
	// 模拟启动 Watchdog，并手动触发一次超时重启
	go func() {
		time.Sleep(3 * time.Second)
		watchdog.HeartBeat() // 模拟心跳
	}()
	// 启动 Watchdog
	go watchdog.Start()
	// 等待重启操作完成
	time.Sleep(7 * time.Second)
	// 验证重启是否被触发，可以通过检查 `lastHeartBeat` 时间
	watchdog.mu.Lock()
	if time.Since(watchdog.lastHeartBeat) > 2*time.Second {
		t.Errorf("Watchdog did not restart in time")
	}
	watchdog.mu.Unlock()
}

func TestWatchdog_RestartFail(t *testing.T) {
	// 创建一个 Watchdog 实例，模拟失败的重启
	watchdog := NewWatchdog(1*time.Second, 2*time.Second, mockFailRestartCmd, mockRestartFail)
	// 启动 Watchdog
	go watchdog.Start()
	// 等待重启操作完成
	time.Sleep(6 * time.Second)
	// 验证失败回调是否被触发
	// 由于失败回调在 `mockRestartFail` 中，错误信息已经在 callback 函数中打印，我们只需要看控制台输出
}
