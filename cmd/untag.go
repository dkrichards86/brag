package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var untagCmd = &cobra.Command{
	Use:   "untag <line-number|last> <tag> [tags...]",
	Short: "Remove one or more tags from an existing win",
	Long:  `Remove one or more tags from an existing win by line number. Use 'last' or 'latest' to untag the most recent win.`,
	Args:  cobra.MinimumNArgs(2),
	Run:   removeTag,
}

func init() {
	rootCmd.AddCommand(untagCmd)
}

func removeTag(cmd *cobra.Command, args []string) {
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

	// Check which tags exist and which don't
	var removedTags []string
	var missingTags []string
	for _, tag := range tags {
		if win.HasTag(tag) {
			removedTags = append(removedTags, tag)
		} else {
			missingTags = append(missingTags, tag)
		}
	}

	// If no tags to remove, nothing to do
	if len(removedTags) == 0 {
		if len(missingTags) == 1 {
			fmt.Printf("Tag %s doesn't exist on this win.\n", missingTags[0])
		} else {
			fmt.Printf("Tags %s don't exist on this win.\n", strings.Join(missingTags, ", "))
		}
		return
	}

	// Remove the tags from the message
	newMessage := win.Message
	for _, tag := range removedTags {
		// Remove tag with space before it, or just the tag if no space
		newMessage = strings.ReplaceAll(newMessage, " "+tag, "")
		newMessage = strings.ReplaceAll(newMessage, tag, "")
	}

	// Clean up any double spaces that might result
	newMessage = strings.Join(strings.Fields(newMessage), " ")

	// Update the win
	if err := store.UpdateWin(index, newMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating win: %v\n", err)
		os.Exit(1)
	}

	// Report results
	if len(removedTags) == 1 {
		fmt.Printf("Removed tag %s from win #%d\n", removedTags[0], index+1)
	} else {
		fmt.Printf("Removed tags %s from win #%d\n", strings.Join(removedTags, ", "), index+1)
	}

	if len(missingTags) > 0 {
		if len(missingTags) == 1 {
			fmt.Printf("(Tag %s was not present)\n", missingTags[0])
		} else {
			fmt.Printf("(Tags %s were not present)\n", strings.Join(missingTags, ", "))
		}
	}
}
