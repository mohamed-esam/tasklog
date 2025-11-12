package cmd

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"tasklog/internal/jira"
	"tasklog/internal/storage"
	"tasklog/internal/tempo"
	"tasklog/internal/timeparse"
	"tasklog/internal/ui"
)

var (
	shortcutName string
	taskKey      string
	timeSpent    string
	label        string
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Log time to a task",
	Long: `Interactively log time to a Jira task. 
You can use shortcuts, select from in-progress tasks, or search for tasks.`,
	RunE: runLog,
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().StringVarP(&shortcutName, "shortcut", "s", "", "Use a predefined shortcut")
	logCmd.Flags().StringVarP(&taskKey, "task", "t", "", "Task key (e.g., PROJ-123)")
	logCmd.Flags().StringVarP(&timeSpent, "time", "d", "", "Time spent (e.g., 2h 30m, 2.5h, 150m)")
	logCmd.Flags().StringVarP(&label, "label", "l", "", "Work log label")
}

func runLog(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := checkConfig()
	if err != nil {
		return err
	}

	// Initialize clients
	jiraClient := jira.NewClient(cfg.Jira.URL, cfg.Jira.Username, cfg.Jira.APIToken, cfg.Jira.ProjectKey)
	tempoClient := tempo.NewClient(cfg.Tempo.APIToken)

	// Initialize storage
	store, err := storage.NewStorage(cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	var selectedIssue *jira.Issue
	var timeSeconds int
	var selectedLabel string

	// Check if using a shortcut
	if shortcutName != "" {
		log.Debug().Str("shortcut", shortcutName).Msg("Using shortcut")

		shortcut, found := cfg.GetShortcut(shortcutName)
		if !found {
			return fmt.Errorf("shortcut '%s' not found in configuration", shortcutName)
		}

		// Use shortcut values
		if taskKey == "" {
			taskKey = shortcut.Task
		}
		if timeSpent == "" && shortcut.Time != "" {
			timeSpent = shortcut.Time
		}
		if label == "" {
			label = shortcut.Label
		}
	}

	// Get task
	if taskKey != "" {
		log.Debug().Str("task", taskKey).Msg("Fetching specified task")
		issue, err := jiraClient.GetIssue(taskKey)
		if err != nil {
			return fmt.Errorf("failed to fetch task %s: %w", taskKey, err)
		}
		selectedIssue = issue
		fmt.Printf("Task: %s - %s\n", selectedIssue.Key, selectedIssue.Fields.Summary)
	} else {
		// Interactive task selection
		log.Debug().Msg("Fetching in-progress tasks")
		inProgressIssues, err := jiraClient.GetInProgressIssues()
		if err != nil {
			return fmt.Errorf("failed to fetch in-progress tasks: %w", err)
		}

		selectedIssue, err = ui.SelectTask(inProgressIssues)
		if err != nil {
			return fmt.Errorf("failed to select task: %w", err)
		}

		// If user chose to search, perform the search
		if selectedIssue.Fields.Summary == "" {
			searchResults, err := jiraClient.SearchIssues(selectedIssue.Key)
			if err != nil {
				return fmt.Errorf("failed to search tasks: %w", err)
			}

			selectedIssue, err = ui.SelectFromSearchResults(searchResults)
			if err != nil {
				return fmt.Errorf("failed to select from search results: %w", err)
			}

			// Fetch full issue details
			issue, err := jiraClient.GetIssue(selectedIssue.Key)
			if err != nil {
				return fmt.Errorf("failed to fetch task details: %w", err)
			}
			selectedIssue = issue
		}
	}

	// Get time spent
	if timeSpent != "" {
		timeSeconds, err = timeparse.Parse(timeSpent)
		if err != nil {
			return fmt.Errorf("invalid time format: %w", err)
		}
	} else {
		timeStr, err := ui.PromptTimeSpent()
		if err != nil {
			return fmt.Errorf("failed to get time spent: %w", err)
		}

		timeSeconds, err = timeparse.Parse(timeStr)
		if err != nil {
			return fmt.Errorf("invalid time format: %w", err)
		}
	}

	// Get label
	if label != "" {
		if !cfg.IsLabelAllowed(label) {
			return fmt.Errorf("label '%s' is not in the allowed labels list", label)
		}
		selectedLabel = label
	} else {
		selectedLabel, err = ui.SelectLabel(cfg.Labels.AllowedLabels)
		if err != nil {
			return fmt.Errorf("failed to select label: %w", err)
		}

		if !cfg.IsLabelAllowed(selectedLabel) {
			return fmt.Errorf("label '%s' is not allowed", selectedLabel)
		}
	}

	// Get optional comment
	comment, err := ui.PromptComment()
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}

	// Confirm before logging
	fmt.Printf("\n")
	fmt.Printf("Task:    %s - %s\n", selectedIssue.Key, selectedIssue.Fields.Summary)
	fmt.Printf("Time:    %s\n", timeparse.Format(timeSeconds))
	fmt.Printf("Label:   %s\n", selectedLabel)
	if comment != "" {
		fmt.Printf("Comment: %s\n", comment)
	}
	fmt.Printf("\n")

	confirmed, err := ui.Confirm("Log this time entry?")
	if err != nil {
		return fmt.Errorf("failed to confirm: %w", err)
	}

	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	// Create time entry
	now := time.Now()
	entry := &storage.TimeEntry{
		IssueKey:         selectedIssue.Key,
		IssueSummary:     selectedIssue.Fields.Summary,
		TimeSpentSeconds: timeSeconds,
		TimeSpent:        timeparse.Format(timeSeconds),
		Label:            selectedLabel,
		Comment:          comment,
		Started:          now,
		SyncedToJira:     false,
		SyncedToTempo:    false,
	}

	// Save to local storage first
	if err := store.AddTimeEntry(entry); err != nil {
		return fmt.Errorf("failed to save time entry locally: %w", err)
	}

	fmt.Println("✓ Saved to local cache")

	// Log to Jira
	log.Debug().Msg("Logging to Jira")
	worklog, err := jiraClient.AddWorklog(selectedIssue.Key, timeSeconds, now, comment)
	if err != nil {
		log.Error().Err(err).Msg("Failed to log to Jira")
		fmt.Printf("⚠ Failed to log to Jira: %v\n", err)
	} else {
		entry.SyncedToJira = true
		entry.JiraWorklogID = worklog.ID
		fmt.Println("✓ Logged to Jira")
	}

	// Log to Tempo
	log.Debug().Msg("Logging to Tempo")
	tempoWorklog, err := tempoClient.AddWorklog(selectedIssue.Key, timeSeconds, now, selectedLabel, comment)
	if err != nil {
		log.Error().Err(err).Msg("Failed to log to Tempo")
		fmt.Printf("⚠ Failed to log to Tempo: %v\n", err)
	} else {
		entry.SyncedToTempo = true
		entry.TempoWorklogID = fmt.Sprintf("%d", tempoWorklog.TempoWorklogID)
		fmt.Println("✓ Logged to Tempo")
	}

	// Update storage with sync status
	if err := store.UpdateTimeEntry(entry); err != nil {
		log.Error().Err(err).Msg("Failed to update time entry sync status")
	}

	// Show today's summary
	fmt.Println()
	if err := showTodaySummary(store); err != nil {
		log.Error().Err(err).Msg("Failed to show summary")
	}

	return nil
}

func showTodaySummary(store *storage.Storage) error {
	total, err := store.GetTodayTotalSeconds()
	if err != nil {
		return err
	}

	entries, err := store.GetTodayEntries()
	if err != nil {
		return err
	}

	fmt.Println("═══════════════════════════════════════════")
	fmt.Printf("Today's Summary (%d entries)\n", len(entries))
	fmt.Println("═══════════════════════════════════════════")

	if len(entries) > 0 {
		for _, entry := range entries {
			syncStatus := ""
			if entry.SyncedToJira && entry.SyncedToTempo {
				syncStatus = "✓"
			} else if entry.SyncedToJira || entry.SyncedToTempo {
				syncStatus = "⚠"
			} else {
				syncStatus = "✗"
			}

			fmt.Printf("%s %s - %-10s [%s] %s\n",
				syncStatus,
				entry.Started.Format("15:04"),
				entry.TimeSpent,
				entry.Label,
				entry.IssueKey,
			)
		}
	}

	fmt.Println("───────────────────────────────────────────")
	fmt.Printf("Total: %s\n", timeparse.Format(total))
	fmt.Println("═══════════════════════════════════════════")

	return nil
}
