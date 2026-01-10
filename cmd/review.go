package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/spf13/cobra"
)

var (
	fromDate     string
	toDate       string
	tagFilter    string
	taggedOnly   bool
	untaggedOnly bool
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review your micro-wins",
	Long:  `Display your micro-wins within a date range. Defaults to the last 7 days.`,
	Run:   reviewWins,
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().StringVar(&fromDate, "from", "", "Start date (YYYY/MM/DD or YYYY-MM-DD)")
	reviewCmd.Flags().StringVar(&toDate, "to", "", "End date (YYYY/MM/DD or YYYY-MM-DD)")
	reviewCmd.Flags().StringVar(&tagFilter, "tag", "", "Filter by tag (e.g., 'work' or '#work')")
	reviewCmd.Flags().BoolVar(&taggedOnly, "tagged", false, "Show only wins that have tags")
	reviewCmd.Flags().BoolVar(&untaggedOnly, "untagged", false, "Show only wins without tags")
}

func reviewWins(cmd *cobra.Command, args []string) {
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

	// Parse date range
	from, to := parseDateRange(fromDate, toDate)

	// Filter wins
	filteredWins := filterWins(wins, from, to, tagFilter, taggedOnly, untaggedOnly)

	if len(filteredWins) == 0 {
		fmt.Println("No wins found for the specified criteria.")
		return
	}

	// Display wins
	fmt.Printf("Wins from %s to %s:\n\n", from.Format("2006-01-02"), to.Format("2006-01-02"))
	for _, win := range filteredWins {
		fmt.Println(win.Format())
	}

	fmt.Printf("\nTotal: %d wins\n", len(filteredWins))
}

func parseDateRange(fromStr, toStr string) (time.Time, time.Time) {
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

func filterWins(wins []*models.Win, from, to time.Time, tag string, taggedOnly, untaggedOnly bool) []*models.Win {
	var filtered []*models.Win

	for _, win := range wins {
		// Date filter
		if win.Timestamp.Before(from) || win.Timestamp.After(to) {
			continue
		}

		// Tagged filter - only show wins that have at least one tag
		if taggedOnly && len(win.Tags) == 0 {
			continue
		}

		// Untagged filter - only show wins without tags
		if untaggedOnly && len(win.Tags) > 0 {
			continue
		}

		// Tag filter - show wins with a specific tag
		if tag != "" && !win.HasTag(tag) {
			continue
		}

		filtered = append(filtered, win)
	}

	return filtered
}
