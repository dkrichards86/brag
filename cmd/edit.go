package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dkrichards86/brag/internal/storage"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [line-number]",
	Short: "Edit wins",
	Long:  `Edit the entire wins file or a specific entry by line number.`,
	Run:   editWin,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func editWin(cmd *cobra.Command, args []string) {
	store, err := storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// If no arguments, open the entire file in $EDITOR
	if len(args) == 0 {
		openInEditor(store.GetFilePath())
		return
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

	// Show current entry
	fmt.Printf("Current entry:\n%s\n\n", wins[index].Format())

	// Prompt for new message
	fmt.Print("Enter new message (or press Ctrl+C to cancel): ")
	reader := bufio.NewReader(os.Stdin)
	newMessage, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	newMessage = strings.TrimSpace(newMessage)
	if newMessage == "" {
		fmt.Println("No changes made.")
		return
	}

	// Update the win
	if err := store.UpdateWin(index, newMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating win: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Win updated successfully!")
}

func openInEditor(filePath string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // Default to vi if EDITOR is not set
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
		os.Exit(1)
	}
}
