package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"tasklog/internal/storage"
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show today's time tracking summary",
	Long:  `Displays a summary of all time entries logged today.`,
	RunE:  runSummary,
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}

func runSummary(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := checkConfig()
	if err != nil {
		return err
	}

	// Initialize storage
	store, err := storage.NewStorage(cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	return showTodaySummary(store)
}
