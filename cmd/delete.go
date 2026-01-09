package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dkrichards86/brag/internal/storage"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <line-number>",
	Short: "Delete a specific win",
	Long:  `Delete a win by its line number (use 'brag list' to see line numbers).`,
	Args:  cobra.ExactArgs(1),
	Run:   deleteWin,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteWin(cmd *cobra.Command, args []string) {
	store, err := storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Parse line number
	lineNum, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid line number '%s'\n", args[0])
		os.Exit(1)
	}

	// Convert to 0-based index
	index := lineNum - 1

	wins, err := store.ReadAllWins()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading wins: %v\n", err)
		os.Exit(1)
	}

	if index < 0 || index >= len(wins) {
		fmt.Fprintf(os.Stderr, "Error: line number %d is out of range (1-%d)\n", lineNum, len(wins))
		os.Exit(1)
	}

	// Show the entry to be deleted
	fmt.Printf("About to delete:\n%s\n\n", wins[index].Format())

	// Confirm deletion
	fmt.Print("Are you sure? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Deletion canceled.")
		return
	}

	// Delete the win
	if err := store.DeleteWin(index); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting win: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Win deleted successfully!")
}
