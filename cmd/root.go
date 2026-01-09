package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dkrichards86/brag/internal/storage"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "brag [message] [date]",
	Short: "Track your micro-wins",
	Long:  `Brag is a CLI tool for tracking daily micro-wins and achievements.`,
	// By default, if a subcommand is not found, cobra will error
	// We need to handle this differently
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If args provided, treat as add command
		if len(args) > 0 {
			addWin(cmd, args)
			return nil
		}
		// No args, show help
		return cmd.Help()
	},
}

func init() {
	// Disable suggestions for unknown commands
	rootCmd.DisableSuggestions = true
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// addWin is the default command - adds a new win
func addWin(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	store, err := storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Parse arguments to extract message and optional date
	message, timestamp := parseAddArgs(args)

	if err := store.AddWin(message, timestamp); err != nil {
		fmt.Fprintf(os.Stderr, "Error adding win: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added: %s\n", message)
}

// parseAddArgs parses the arguments for the add command
// Returns the message and timestamp
func parseAddArgs(args []string) (string, time.Time) {
	// Check if the last argument is a date
	lastArg := args[len(args)-1]
	timestamp := time.Now()

	// Try to parse various date formats
	dateFormats := []string{
		"2006/01/02",
		"2006-01-02",
		"01/02/2006",
		"01-02-2006",
	}

	for _, format := range dateFormats {
		if t, err := time.Parse(format, lastArg); err == nil {
			// Valid date found, use it and exclude from message
			timestamp = t
			args = args[:len(args)-1]
			break
		}
	}

	// Join remaining args as the message
	message := strings.Join(args, " ")
	return message, timestamp
}
