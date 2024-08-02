package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var welcomeCmd = &cobra.Command{
	Use:   "welcome",
	Short: "print welcome",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("welcome!!!")
	},
}
