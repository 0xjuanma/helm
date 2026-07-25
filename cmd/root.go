package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/0xjuanma/helm/internal/config"
	"github.com/0xjuanma/helm/internal/tui"
)

var quickMinutes int

var rootCmd = &cobra.Command{
	Use:   "helm",
	Short: "Helm - A minimalistic TUI timer",
	Long:  "Helm is a minimalistic TUI timer designed for pure focus",
	RunE: func(cmd *cobra.Command, _ []string) error {
		var m tea.Model
		if cmd.Flags().Changed("quick") {
			if err := validateQuickMinutes(quickMinutes); err != nil {
				return err
			}
			m = tui.NewQuickModel(quickMinutes)
		} else {
			m = tui.NewModel()
		}
		p := tea.NewProgram(m, tea.WithAltScreen())
		_, err := p.Run()
		return err
	},
}

func init() {
	rootCmd.Flags().IntVar(&quickMinutes, "quick", 0, "start a quick timer for the given number of minutes")
}

func validateQuickMinutes(m int) error {
	if m < config.MinStepMinutes || m > config.MaxStepMinutes {
		return fmt.Errorf("--quick must be between %d and %d minutes", config.MinStepMinutes, config.MaxStepMinutes)
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
