package cmd

import (
	"fmt"
	"os"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/spf13/cobra"
)

var (
	listFromDate string
	listToDate   string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all micro-wins with line numbers",
	Long:  `Display all your micro-wins with line numbers for editing/deleting.`,
	Run:   listWins,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listFromDate, "from", "", "Start date (YYYY/MM/DD or YYYY-MM-DD)")
	listCmd.Flags().StringVar(&listToDate, "to", "", "End date (YYYY/MM/DD or YYYY-MM-DD)")
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

	// Apply date filters if specified
	var filteredWins []*models.Win
	if listFromDate != "" || listToDate != "" {
		from, to := parseDateRange(listFromDate, listToDate)
		filteredWins = filterWins(wins, from, to, "", false, false)
	} else {
		filteredWins = wins
	}

	if len(filteredWins) == 0 {
		fmt.Println("No wins found.")
		return
	}

	// Display wins with line numbers
	for i, win := range filteredWins {
		fmt.Printf("%d. %s\n", i+1, win.Format())
	}

	fmt.Printf("\nTotal: %d wins\n", len(filteredWins))
}
