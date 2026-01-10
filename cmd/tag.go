package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag <line-number> <tag>",
	Short: "Add a tag to an existing win",
	Long:  `Add a tag to an existing win by line number. The tag will be appended to the message.`,
	Args:  cobra.ExactArgs(2),
	Run:   addTag,
}

func init() {
	rootCmd.AddCommand(tagCmd)
}

func addTag(cmd *cobra.Command, args []string) {
	store, err := getStorage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Parse line number
	lineNum, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid line number '%s'\n", args[0])
		os.Exit(1)
	}

	// Convert to 0-based index
	index := lineNum - 1

	// Get the tag and ensure it starts with #
	tag := strings.TrimSpace(args[1])
	if !strings.HasPrefix(tag, "#") {
		tag = "#" + tag
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading wins: %v\n", err)
		os.Exit(1)
	}

	if index < 0 || index >= len(wins) {
		fmt.Fprintf(os.Stderr, "Error: line number %d is out of range (1-%d)\n", lineNum, len(wins))
		os.Exit(1)
	}

	win := wins[index]

	// Check if the tag already exists
	if win.HasTag(tag) {
		fmt.Printf("Tag %s already exists on this win.\n", tag)
		return
	}

	// Add the tag to the message
	newMessage := win.Message + " " + tag

	// Update the win
	if err := store.UpdateWin(index, newMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating win: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added tag %s to win #%d\n", tag, lineNum)
}
