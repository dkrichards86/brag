package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search for wins containing a keyword",
	Long:  `Search through your wins for entries containing the specified keyword or phrase.`,
	Args:  cobra.MinimumNArgs(1),
	Run:   searchWins,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func searchWins(cmd *cobra.Command, args []string) {
	store, err := storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading wins: %v\n", err)
		os.Exit(1)
	}

	// Join all args as search query
	query := strings.Join(args, " ")
	queryLower := strings.ToLower(query)

	var matches []struct {
		win   *models.Win
		index int
	}

	for i, win := range wins {
		if strings.Contains(strings.ToLower(win.Message), queryLower) {
			matches = append(matches, struct {
				win   *models.Win
				index int
			}{win, i})
		}
	}

	if len(matches) == 0 {
		fmt.Printf("No wins found containing '%s'\n", query)
		return
	}

	fmt.Printf("Found %d wins containing '%s':\n\n", len(matches), query)
	for _, match := range matches {
		fmt.Printf("%d. %s\n", match.index+1, match.win.Format())
	}
}
