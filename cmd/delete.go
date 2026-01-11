package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <line-number>",
	Short: "Delete a specific win",
	Long:  `Delete a win by its line number (use 'brag list' to see line numbers).`,
	Args:  cobra.ExactArgs(1),
	RunE:  deleteWin,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteWin(cmd *cobra.Command, args []string) error {
	store, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Parse line number
	lineNum, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid line number '%s'", args[0])
	}

	// Convert to 0-based index
	index := lineNum - 1

	wins, err := store.ReadAllWins()
	if err != nil {
		return fmt.Errorf("failed to read wins: %w", err)
	}

	if index < 0 || index >= len(wins) {
		return fmt.Errorf("line number %d is out of range (1-%d)", lineNum, len(wins))
	}

	// Show the entry to be deleted
	fmt.Printf("About to delete:\n%s\n\n", wins[index].Format())

	// Confirm deletion
	fmt.Print("Are you sure? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Deletion canceled.")
		return nil
	}

	// Delete the win
	if err := store.DeleteWin(index); err != nil {
		return fmt.Errorf("failed to delete win: %w", err)
	}

	fmt.Println("Win deleted successfully!")
	return nil
}
