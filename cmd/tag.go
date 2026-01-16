package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
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
	store, index := mustGetStorageAndIndex(args[0])
	tags := normalizeTags(args[1:])
	wins := mustReadWins(store)
	mustValidateIndex(index, len(wins))

	win := wins[index]
	newTags, existingTags := partitionNewAndExistingTags(win, tags)

	if len(newTags) == 0 {
		displayAllTagsExistMessage(existingTags)
		return
	}

	newMessage := win.Message + " " + strings.Join(newTags, " ")
	mustUpdateWin(store, index, newMessage)
	displayAddTagResults(newTags, existingTags, index)
}

func normalizeTags(args []string) []string {
	var tags []string
	for _, arg := range args {
		tag := strings.TrimSpace(arg)
		if !strings.HasPrefix(tag, "#") {
			tag = "#" + tag
		}
		tags = append(tags, tag)
	}
	return tags
}

func partitionNewAndExistingTags(win interface{ HasTag(string) bool }, tags []string) ([]string, []string) {
	var newTags []string
	var existingTags []string
	for _, tag := range tags {
		if win.HasTag(tag) {
			existingTags = append(existingTags, tag)
		} else {
			newTags = append(newTags, tag)
		}
	}
	return newTags, existingTags
}

func displayAllTagsExistMessage(existingTags []string) {
	if len(existingTags) == 1 {
		fmt.Printf("Tag %s already exists on this win.\n", existingTags[0])
	} else {
		fmt.Printf("Tags %s already exist on this win.\n", strings.Join(existingTags, ", "))
	}
}

func displayAddTagResults(newTags, existingTags []string, index int) {
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

func mustGetStorageAndIndex(lineArg string) (*storage.Storage, int) {
	store, err := getStorage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	index, err := parseLineNumber(lineArg, store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return store, index
}

func mustReadWins(store *storage.Storage) []*models.Win {
	wins, err := store.ReadAllWins()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading wins: %v\n", err)
		os.Exit(1)
	}
	return wins
}

func mustValidateIndex(index, winsCount int) {
	if index < 0 || index >= winsCount {
		fmt.Fprintf(os.Stderr, "Error: line number is out of range (1-%d)\n", winsCount)
		os.Exit(1)
	}
}

func mustUpdateWin(store *storage.Storage, index int, newMessage string) {
	if err := store.UpdateWin(index, newMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating win: %v\n", err)
		os.Exit(1)
	}
}
