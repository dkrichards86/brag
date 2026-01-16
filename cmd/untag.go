package cmd

import (
	"fmt"
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
	store, index := mustGetStorageAndIndex(args[0])
	tags := normalizeTags(args[1:])
	wins := mustReadWins(store)
	mustValidateIndex(index, len(wins))

	win := wins[index]
	removedTags, missingTags := partitionRemovedAndMissingTags(win, tags)

	if len(removedTags) == 0 {
		displayNoTagsRemovedMessage(missingTags)
		return
	}

	newMessage := removeTagsFromMessage(win.Message, removedTags)
	mustUpdateWin(store, index, newMessage)
	displayRemoveTagResults(removedTags, missingTags, index)
}

func partitionRemovedAndMissingTags(win interface{ HasTag(string) bool }, tags []string) ([]string, []string) {
	var removedTags []string
	var missingTags []string
	for _, tag := range tags {
		if win.HasTag(tag) {
			removedTags = append(removedTags, tag)
		} else {
			missingTags = append(missingTags, tag)
		}
	}
	return removedTags, missingTags
}

func displayNoTagsRemovedMessage(missingTags []string) {
	if len(missingTags) == 1 {
		fmt.Printf("Tag %s doesn't exist on this win.\n", missingTags[0])
	} else {
		fmt.Printf("Tags %s don't exist on this win.\n", strings.Join(missingTags, ", "))
	}
}

func removeTagsFromMessage(message string, tags []string) string {
	newMessage := message
	for _, tag := range tags {
		// Remove tag with space before it, or just the tag if no space
		newMessage = strings.ReplaceAll(newMessage, " "+tag, "")
		newMessage = strings.ReplaceAll(newMessage, tag, "")
	}
	// Clean up any double spaces that might result
	return strings.Join(strings.Fields(newMessage), " ")
}

func displayRemoveTagResults(removedTags, missingTags []string, index int) {
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
