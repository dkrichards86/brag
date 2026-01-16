package cmd

import (
	"fmt"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/utils"
	"github.com/spf13/cobra"
)

var (
	listFromDate     string
	listToDate       string
	listTagFilter    string
	listTaggedOnly   bool
	listUntaggedOnly bool
	listAll          bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List micro-wins with line numbers",
	Long:  `Display your micro-wins with line numbers for editing/deleting. Defaults to last 7 days.`,
	RunE:  listWins,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listFromDate, "from", "", "Start date (YYYY/MM/DD or YYYY-MM-DD)")
	listCmd.Flags().StringVar(&listToDate, "to", "", "End date (YYYY/MM/DD or YYYY-MM-DD)")
	listCmd.Flags().StringVar(&listTagFilter, "tag", "", "Filter by tag (e.g., 'work' or '#work')")
	listCmd.Flags().BoolVar(&listTaggedOnly, "tagged", false, "Show only wins that have tags")
	listCmd.Flags().BoolVar(&listUntaggedOnly, "untagged", false, "Show only wins without tags")
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all wins (no date filter)")
}

func listWins(cmd *cobra.Command, args []string) error {
	if err := validateListFlags(); err != nil {
		return err
	}

	store, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		return fmt.Errorf("failed to read wins: %w", err)
	}

	from, to := parseListDateRange(listFromDate, listToDate, listAll)
	filteredWins := filterAndIndexWins(wins, from, to)

	displayListResults(filteredWins, wins, from, to)
	return nil
}

func validateListFlags() error {
	if listTaggedOnly && listUntaggedOnly {
		return fmt.Errorf("cannot use --tagged and --untagged together")
	}

	if listTagFilter != "" && listUntaggedOnly {
		return fmt.Errorf("cannot use --tag and --untagged together (--tag filters for wins with a specific tag)")
	}

	if listAll && (listFromDate != "" || listToDate != "") {
		return fmt.Errorf("cannot use --all with --from or --to (--all shows all wins regardless of date)")
	}

	return nil
}

type indexedWin struct {
	win   *models.Win
	index int // Original index in the full wins array (1-based for display)
}

func filterAndIndexWins(wins []*models.Win, from, to time.Time) []indexedWin {
	var filteredWins []indexedWin
	for i, win := range wins {
		if shouldIncludeWin(win, from, to) {
			filteredWins = append(filteredWins, indexedWin{win: win, index: i + 1})
		}
	}
	return filteredWins
}

func shouldIncludeWin(win *models.Win, from, to time.Time) bool {
	// Date filter
	if win.Timestamp.Before(from) || win.Timestamp.After(to) {
		return false
	}

	// Tagged filter - only show wins that have at least one tag
	if listTaggedOnly && len(win.Tags) == 0 {
		return false
	}

	// Untagged filter - only show wins without tags
	if listUntaggedOnly && len(win.Tags) > 0 {
		return false
	}

	// Tag filter - show wins with a specific tag
	if listTagFilter != "" && !win.HasTag(listTagFilter) {
		return false
	}

	return true
}

func displayListResults(filteredWins []indexedWin, allWins []*models.Win, from, to time.Time) {
	if len(filteredWins) == 0 {
		displayNoWinsMessage(allWins)
		return
	}

	displayListHeader(from, to)
	for _, iw := range filteredWins {
		fmt.Printf("%d. %s\n", iw.index, iw.win.Format())
	}
	fmt.Printf("\nTotal: %d wins\n", len(filteredWins))
}

func displayNoWinsMessage(allWins []*models.Win) {
	if len(allWins) == 0 {
		fmt.Println("No wins found. Add your first win with: brag \"your achievement here\"")
	} else {
		fmt.Println("No wins found for the specified criteria.")
		fmt.Println("Try adjusting your filters or use 'brag list' to see all recent wins.")
	}
}

func displayListHeader(from, to time.Time) {
	if listAll {
		fmt.Println("All wins:")
		fmt.Println()
	} else {
		fmt.Printf("Wins from %s to %s:\n\n", from.Format("2006-01-02"), to.Format("2006-01-02"))
	}
}

func parseListDateRange(fromStr, toStr string, showAll bool) (time.Time, time.Time) {
	// If --all is specified, return a very wide date range
	if showAll {
		// From: year 2000 (before brag existed)
		from := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		// To: far future
		to := time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC)
		return from, to
	}

	// Default: last 7 days to today
	to := time.Now()
	from := to.AddDate(0, 0, -7)

	// Parse from date
	if fromStr != "" {
		if t, ok := utils.ParseDate(fromStr); ok {
			from = t
		}
	}

	// Parse to date
	if toStr != "" {
		if t, ok := utils.ParseDate(toStr); ok {
			to = t
		}
	}

	// Set to end of day for "to" date
	to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, to.Location())

	return from, to
}
