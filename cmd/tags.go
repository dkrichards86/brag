package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags and their frequency",
	Long:  `List all tags found in your wins, sorted by frequency.`,
	RunE:  listTags,
}

func init() {
	rootCmd.AddCommand(tagsCmd)
}

func listTags(cmd *cobra.Command, args []string) error {
	store, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		return fmt.Errorf("failed to read wins: %w", err)
	}

	// Count tag frequencies
	tagCounts := make(map[string]int)
	for _, win := range wins {
		for _, tag := range win.Tags {
			tagCounts[tag]++
		}
	}

	if len(tagCounts) == 0 {
		fmt.Println("No tags found in your wins.")
		return nil
	}

	// Convert to slice for sorting
	type tagFreq struct {
		tag   string
		count int
	}
	var tags []tagFreq
	for tag, count := range tagCounts {
		tags = append(tags, tagFreq{tag, count})
	}

	// Sort by count (descending), then by tag name (ascending)
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].count == tags[j].count {
			return tags[i].tag < tags[j].tag
		}
		return tags[i].count > tags[j].count
	})

	// Print the results
	fmt.Printf("Found %d unique tags:\n\n", len(tags))
	for _, t := range tags {
		fmt.Printf("%-20s %d\n", t.tag, t.count)
	}
	return nil
}
