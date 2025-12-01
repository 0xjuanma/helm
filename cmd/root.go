package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/0xjuanma/helm/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "helm",
	Short: "Helm - A minimalistic TUI Pomodoro timer",
	Long:  "Helm is a minimalistic TUI Pomodoro timer designed for pure focus",
	RunE: func(_ *cobra.Command, _ []string) error {
		p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
		_, err := p.Run()
		return err
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
