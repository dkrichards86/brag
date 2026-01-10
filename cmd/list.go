package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/spf13/cobra"
)

var (
	listFromDate     string
	listToDate       string
	listTagFilter    string
	listTaggedOnly   bool
	listUntaggedOnly bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List micro-wins with line numbers",
	Long:  `Display your micro-wins with line numbers for editing/deleting. Defaults to last 7 days.`,
	Run:   listWins,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listFromDate, "from", "", "Start date (YYYY/MM/DD or YYYY-MM-DD)")
	listCmd.Flags().StringVar(&listToDate, "to", "", "End date (YYYY/MM/DD or YYYY-MM-DD)")
	listCmd.Flags().StringVar(&listTagFilter, "tag", "", "Filter by tag (e.g., 'work' or '#work')")
	listCmd.Flags().BoolVar(&listTaggedOnly, "tagged", false, "Show only wins that have tags")
	listCmd.Flags().BoolVar(&listUntaggedOnly, "untagged", false, "Show only wins without tags")
}

func listWins(cmd *cobra.Command, args []string) {
	store, err := getStorage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading wins: %v\n", err)
		os.Exit(1)
	}

	// Parse date range (defaults to last 7 days)
	from, to := parseListDateRange(listFromDate, listToDate)

	// Filter wins and track their original indices
	type indexedWin struct {
		win   *models.Win
		index int // Original index in the full wins array (1-based for display)
	}

	var filteredWins []indexedWin
	for i, win := range wins {
		// Date filter
		if win.Timestamp.Before(from) || win.Timestamp.After(to) {
			continue
		}

		// Tagged filter - only show wins that have at least one tag
		if listTaggedOnly && len(win.Tags) == 0 {
			continue
		}

		// Untagged filter - only show wins without tags
		if listUntaggedOnly && len(win.Tags) > 0 {
			continue
		}

		// Tag filter - show wins with a specific tag
		if listTagFilter != "" && !win.HasTag(listTagFilter) {
			continue
		}

		filteredWins = append(filteredWins, indexedWin{win: win, index: i + 1})
	}

	if len(filteredWins) == 0 {
		fmt.Println("No wins found for the specified criteria.")
		return
	}

	// Display wins with their actual file line numbers
	fmt.Printf("Wins from %s to %s:\n\n", from.Format("2006-01-02"), to.Format("2006-01-02"))
	for _, iw := range filteredWins {
		fmt.Printf("%d. %s\n", iw.index, iw.win.Format())
	}

	fmt.Printf("\nTotal: %d wins\n", len(filteredWins))
}

func parseListDateRange(fromStr, toStr string) (time.Time, time.Time) {
	// Default: last 7 days to today
	to := time.Now()
	from := to.AddDate(0, 0, -7)

	dateFormats := []string{
		"2006/01/02",
		"2006-01-02",
		"01/02/2006",
		"01-02-2006",
	}

	// Parse from date
	if fromStr != "" {
		for _, format := range dateFormats {
			if t, err := time.Parse(format, fromStr); err == nil {
				from = t
				break
			}
		}
	}

	// Parse to date
	if toStr != "" {
		for _, format := range dateFormats {
			if t, err := time.Parse(format, toStr); err == nil {
				to = t
				break
			}
		}
	}

	// Set to end of day for "to" date
	to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, to.Location())

	return from, to
}
