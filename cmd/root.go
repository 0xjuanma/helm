package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "helm",
	Short: "Helm",
	Long:  "Helm is a minimalistic TUI Pomodoro timer designed for pure focus",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Helm 0.1.0")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
