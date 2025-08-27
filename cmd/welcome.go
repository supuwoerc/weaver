package cmd

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/supuwoerc/weaver/bootstrap"
)

var welcomeCmd = &cobra.Command{
	Use:   "welcome",
	Short: "print welcome",
	Run: func(cmd *cobra.Command, args []string) {
		cli := bootstrap.WireCli()
		cli.Logger.Infow("welcome cli is running...", "config.env", cli.Conf.Env)
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
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
