package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkrichards86/brag/internal/storage"
	"github.com/dkrichards86/brag/internal/utils"
	"github.com/spf13/cobra"
)

var (
	// Storage configuration flags
	bragDir  string
	winsFile string
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

	// Add persistent flags for storage configuration
	rootCmd.PersistentFlags().StringVar(&bragDir, "brag-dir", "", "Directory for brag storage (default: ~/.brag, env: BRAG_DIR)")
	rootCmd.PersistentFlags().StringVar(&winsFile, "wins-file", "", "Wins file name (default: wins.txt, env: BRAG_FILE)")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// getStorage creates a storage instance using flag values, env vars, or defaults
// Priority: CLI flags > Environment variables > Defaults
func getStorage() (*storage.Storage, error) {
	// Determine directory name: flag > env var > default
	dirName := bragDir
	if dirName == "" {
		dirName = os.Getenv("BRAG_DIR")
	}
	if dirName == "" {
		dirName = ".brag"
	}

	// Determine file name: flag > env var > default
	fileName := winsFile
	if fileName == "" {
		fileName = os.Getenv("BRAG_FILE")
	}
	if fileName == "" {
		fileName = "wins.txt"
	}

	// Build the full path
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	filePath := filepath.Join(home, dirName, fileName)
	return storage.NewWithPath(filePath)
}

// addWin is the default command - adds a new win
func addWin(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	store, err := getStorage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Parse arguments to extract message and optional date
	message, timestamp := parseAddArgs(args)

	// Validate that message is not empty
	if message == "" {
		fmt.Fprintln(os.Stderr, "Error: message cannot be empty")
		os.Exit(1)
	}

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

	// Try to parse the last argument as a date
	if t, ok := utils.ParseDate(lastArg); ok {
		// Valid date found, use it and exclude from message
		timestamp = t
		args = args[:len(args)-1]
	}

	// Join remaining args as the message and trim whitespace
	message := strings.TrimSpace(strings.Join(args, " "))
	return message, timestamp
}

// parseLineNumber converts a line number string to an index
// Supports special keywords like "last", "latest"
// Returns the 0-based index and an error if invalid
func parseLineNumber(lineNumStr string, store *storage.Storage) (int, error) {
	// Handle special keywords
	lineNumStr = strings.ToLower(strings.TrimSpace(lineNumStr))
	if lineNumStr == "last" || lineNumStr == "latest" {
		wins, err := store.ReadAllWins()
		if err != nil {
			return 0, fmt.Errorf("failed to read wins: %w", err)
		}
		if len(wins) == 0 {
			return 0, fmt.Errorf("no wins found (use 'brag \"your message\"' to add your first win)")
		}
		return len(wins) - 1, nil
	}

	// Parse as regular number
	var num int
	n, err := fmt.Sscanf(lineNumStr, "%d", &num)
	if err != nil || n != 1 {
		return 0, fmt.Errorf("invalid line number '%s' (use 'brag list' to see valid line numbers)", lineNumStr)
	}

	// Validate range
	wins, err := store.ReadAllWins()
	if err != nil {
		return 0, fmt.Errorf("failed to read wins: %w", err)
	}

	index := num - 1
	if index < 0 || index >= len(wins) {
		if len(wins) == 0 {
			return 0, fmt.Errorf("no wins found (use 'brag \"your message\"' to add your first win)")
		}
		return 0, fmt.Errorf("line number %d is out of range (valid range: 1-%d)\nUse 'brag list' to see all wins with line numbers", num, len(wins))
	}

	return index, nil
}
