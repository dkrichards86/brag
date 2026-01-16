package cmd

import (
	"fmt"
	"os"
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

	// Check which tags are new and which already exist
	var newTags []string
	var existingTags []string
	for _, tag := range tags {
		if win.HasTag(tag) {
			existingTags = append(existingTags, tag)
		} else {
			newTags = append(newTags, tag)
		}
	}

	// If all tags already exist, nothing to do
	if len(newTags) == 0 {
		if len(existingTags) == 1 {
			fmt.Printf("Tag %s already exists on this win.\n", existingTags[0])
		} else {
			fmt.Printf("Tags %s already exist on this win.\n", strings.Join(existingTags, ", "))
		}
		return
	}

	// Add the new tags to the message
	newMessage := win.Message + " " + strings.Join(newTags, " ")

	// Update the win
	if err := store.UpdateWin(index, newMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating win: %v\n", err)
		os.Exit(1)
	}

	// Report results
	if len(newTags) == 1 {
		fmt.Printf("Added tag %s to win #%d\n", newTags[0], index+1)
	} else {
		fmt.Printf("Added tags %s to win #%d\n", strings.Join(newTags, ", "), index+1)
	}

	if len(existingTags) > 0 {
		if len(existingTags) == 1 {
			fmt.Printf("(Tag %s was already present)\n", existingTags[0])
		} else {
			fmt.Printf("(Tags %s were already present)\n", strings.Join(existingTags, ", "))
		}
	}
}
