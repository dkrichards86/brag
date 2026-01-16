package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/utils"
	"github.com/spf13/cobra"
)

var (
	searchFromDate string
	searchToDate   string
	searchTag      string
)

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search for wins containing a keyword",
	Long:  `Search through your wins for entries containing the specified keyword or phrase.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  searchWins,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&searchFromDate, "from", "", "Start date (YYYY/MM/DD or YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&searchToDate, "to", "", "End date (YYYY/MM/DD or YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&searchTag, "tag", "", "Filter by tag (e.g., 'work' or '#work')")
}

func searchWins(cmd *cobra.Command, args []string) error {
	store, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		return fmt.Errorf("failed to read wins: %w", err)
	}

	// Parse date range if provided
	from, to := parseSearchDateRange(searchFromDate, searchToDate)

	// Join all args as search query
	query := strings.Join(args, " ")
	queryLower := strings.ToLower(query)

	var matches []struct {
		win   *models.Win
		index int
	}

	for i, win := range wins {
		// Text match
		if !strings.Contains(strings.ToLower(win.Message), queryLower) {
			continue
		}

		// Date filter
		if from != nil && win.Timestamp.Before(*from) {
			continue
		}
		if to != nil && win.Timestamp.After(*to) {
			continue
		}

		// Tag filter
		if searchTag != "" && !win.HasTag(searchTag) {
			continue
		}

		matches = append(matches, struct {
			win   *models.Win
			index int
		}{win, i})
	}

	if len(matches) == 0 {
		filters := buildSearchFilterDescription(searchFromDate, searchToDate, searchTag)
		fmt.Printf("No wins found containing '%s'%s\n", query, filters)
		return nil
	}

	filters := buildSearchFilterDescription(searchFromDate, searchToDate, searchTag)
	fmt.Printf("Found %d wins containing '%s'%s:\n\n", len(matches), query, filters)
	for _, match := range matches {
		fmt.Printf("%d. %s\n", match.index+1, match.win.Format())
	}
	return nil
}

func parseSearchDateRange(fromStr, toStr string) (*time.Time, *time.Time) {
	var from, to *time.Time

	if fromStr != "" {
		if t, ok := utils.ParseDate(fromStr); ok {
			from = &t
		}
	}

	if toStr != "" {
		if t, ok := utils.ParseDate(toStr); ok {
			// Set to end of day
			endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
			to = &endOfDay
		}
	}

	return from, to
}

func buildSearchFilterDescription(fromStr, toStr, tag string) string {
	var filters []string

	if fromStr != "" || toStr != "" {
		if fromStr != "" && toStr != "" {
			filters = append(filters, fmt.Sprintf("from %s to %s", fromStr, toStr))
		} else if fromStr != "" {
			filters = append(filters, fmt.Sprintf("from %s", fromStr))
		} else {
			filters = append(filters, fmt.Sprintf("to %s", toStr))
		}
	}

	if tag != "" {
		if !strings.HasPrefix(tag, "#") {
			tag = "#" + tag
		}
		filters = append(filters, fmt.Sprintf("with tag %s", tag))
	}

	if len(filters) == 0 {
		return ""
	}

	return " (" + strings.Join(filters, ", ") + ")"
}
