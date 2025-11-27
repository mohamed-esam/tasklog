package cmd

import (
	"fmt"
	"os"

	"tasklog/internal/config"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasklog",
	Short: "Interactive time tracking tool with Jira and Tempo integration",
	Long: `Tasklog is an interactive CLI tool for tracking time on Jira tasks.
It integrates with Jira Cloud API and Tempo to help you log time efficiently.`,
	Version: GetVersion(),
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Ensure config directory exists
	if err := config.EnsureConfigDir(); err != nil {
		log.Error().Err(err).Msg("Failed to ensure config directory")
	}
}

func checkConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "Please create a config file at ~/.tasklog/config.yaml\n")
		fmt.Fprintf(os.Stderr, "See config.example.yaml for an example configuration.\n")
		return nil, err
	}
	return cfg, nil
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func SetVersionInfo(v, c, d, b string) {
	version = v
	commit = c
	date = d
	builtBy = b
}

func GetVersion() string {
	return fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(_ *cobra.Command, _ []string) {
		logger := log.Logger.With().Str("component", "version").Logger()
		logger.Info().
			Str("commit", commit).
			Str("built_at", date).
			Str("built_by", builtBy).
			Msg("tasklog version information")
	},
}
