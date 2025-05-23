package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var welcomeCmd = &cobra.Command{
	Use:   "welcome",
	Short: "print welcome",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-signalCh
			cmd.Printf("\n监听到取消信号,取消执行")
			cancel()
		}()
		count := 10
		bar := progressbar.NewOptions(count,
			progressbar.OptionSetDescription("执行进度"),
			progressbar.OptionUseANSICodes(false),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowIts(),
			progressbar.OptionSetTheme(progressbar.ThemeASCII),
			progressbar.OptionShowElapsedTimeOnFinish())
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			default:
				_ = bar.Add(1)
				time.Sleep(1 * time.Second)
				count--
				if count == 0 {
					break loop
				}
			}
		}
		cmd.Printf("\n脚本执行结束\n")
	},
}
