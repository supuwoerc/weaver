package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing cli '%s'", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(welcomeCmd)
}
