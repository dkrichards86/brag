package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	deleteForce bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <line-number|last>",
	Short: "Delete a specific win",
	Long:  `Delete a win by its line number (use 'brag list' to see line numbers). Use 'last' or 'latest' to delete the most recent win.`,
	Args:  cobra.ExactArgs(1),
	RunE:  deleteWin,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "y", false, "Skip confirmation prompt")
}

func deleteWin(cmd *cobra.Command, args []string) error {
	store, err := getStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Parse line number (supports "last" keyword)
	index, err := parseLineNumber(args[0], store)
	if err != nil {
		return err
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		return fmt.Errorf("failed to read wins: %w", err)
	}

	// Show the entry to be deleted
	fmt.Printf("About to delete:\n%s\n\n", wins[index].Format())

	// Confirm deletion unless --force/-y is set
	if !deleteForce {
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
	}

	// Delete the win
	if err := store.DeleteWin(index); err != nil {
		return fmt.Errorf("failed to delete win: %w", err)
	}

	fmt.Println("Win deleted successfully!")
	return nil
}
