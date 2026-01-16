package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag <line-number|last> <tag> [tags...]",
	Short: "Add one or more tags to an existing win",
	Long:  `Add one or more tags to an existing win by line number. Use 'last' or 'latest' to tag the most recent win. The tags will be appended to the message.`,
	Args:  cobra.MinimumNArgs(2),
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

	// Parse line number (supports "last" keyword)
	index, err := parseLineNumber(args[0], store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Process all tags and ensure they start with #
	var tags []string
	for _, arg := range args[1:] {
		tag := strings.TrimSpace(arg)
		if !strings.HasPrefix(tag, "#") {
			tag = "#" + tag
		}
		tags = append(tags, tag)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading wins: %v\n", err)
		os.Exit(1)
	}

	if index < 0 || index >= len(wins) {
		fmt.Fprintf(os.Stderr, "Error: line number is out of range (1-%d)\n", len(wins))
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
