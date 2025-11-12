package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"tasklog/internal/config"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize tasklog configuration",
	Long:  `Creates the configuration directory and an example config file at ~/.tasklog/config.yaml`,
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Ensure config directory exists
	if err := config.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".tasklog", "config.yaml")

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists at: %s\n", configPath)
		fmt.Println("To reinitialize, delete the existing file and run this command again.")
		return nil
	}

	// Read example config from current directory
	examplePath := "config.example.yaml"
	exampleData, err := os.ReadFile(examplePath)
	if err != nil {
		// If example doesn't exist in current directory, create a basic one
		exampleData = []byte(`# Tasklog Configuration
# Fill in your credentials below

# Required: Jira configuration
jira:
  url: "https://your-domain.atlassian.net"
  username: "your-email@example.com"
  api_token: "your-jira-api-token"
  project_key: "PROJ"  # Project key to filter tasks

# Required: Tempo configuration
tempo:
  api_token: "your-tempo-api-token"

# Optional: Filter labels that can be used for time logging
labels:
  allowed_labels:
    - "development"
    - "code-review"
    - "meeting"
    - "testing"
    - "documentation"
    - "bug-fix"

# Optional: Define shortcuts for quick time logging
shortcuts:
  - name: "daily"
    task: "PROJ-123"
    time: "30m"
    label: "meeting"
  
  - name: "standup"
    task: "PROJ-123"
    time: "15m"
    label: "meeting"

# Optional: Database path (defaults to ~/.tasklog/tasklog.db)
database:
  path: ""
`)
	}

	// Write config file
	if err := os.WriteFile(configPath, exampleData, 0600); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Println("✓ Configuration initialized successfully!")
	fmt.Printf("\nConfig file created at: %s\n", configPath)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Edit the config file with your Jira and Tempo credentials")
	fmt.Println("2. Set the Jira project_key for your project (required)")
	fmt.Println("3. Get your Jira API token: https://id.atlassian.com/manage-profile/security/api-tokens")
	fmt.Println("4. Get your Tempo API token from Tempo > Settings > API Integration")
	fmt.Println("5. (Optional) Configure labels and shortcuts")
	fmt.Printf("6. Run: tasklog log\n")

	return nil
}
