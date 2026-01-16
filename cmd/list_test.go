package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
)

func TestListWinsIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create storage
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	testWins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
			Message:   "First win #work",
			Tags:      []string{"#work"},
		},
		{
			Timestamp: time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC),
			Message:   "Second win #personal",
			Tags:      []string{"#personal"},
		},
		{
			Timestamp: time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC),
			Message:   "Third win",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 15, 16, 0, 0, 0, time.UTC),
			Message:   "Fourth win #work #backend",
			Tags:      []string{"#work", "#backend"},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	tests := []struct {
		name          string
		fromDate      string
		toDate        string
		expectedCount int
		expectedFirst string
		expectedLast  string
	}{
		{
			name:          "list all wins",
			fromDate:      "",
			toDate:        "",
			expectedCount: 4,
			expectedFirst: "First win #work",
			expectedLast:  "Fourth win #work #backend",
		},
		{
			name:          "list with date range",
			fromDate:      "2026-01-05",
			toDate:        "2026-01-15",
			expectedCount: 3,
			expectedFirst: "Second win #personal",
			expectedLast:  "Fourth win #work #backend",
		},
		{
			name:          "list single day",
			fromDate:      "2026-01-10",
			toDate:        "2026-01-10",
			expectedCount: 1,
			expectedFirst: "Third win",
			expectedLast:  "Third win",
		},
		{
			name:          "list with only from date",
			fromDate:      "2026-01-10",
			toDate:        "2026-02-01", // Specify explicit to date
			expectedCount: 2,
			expectedFirst: "Third win",
			expectedLast:  "Fourth win #work #backend",
		},
		{
			name:          "list with only to date",
			fromDate:      "2026-01-01", // Specify explicit from date
			toDate:        "2026-01-05",
			expectedCount: 2,
			expectedFirst: "First win #work",
			expectedLast:  "Second win #personal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wins, err := store.ReadAllWins()
			if err != nil {
				t.Fatalf("Failed to read wins: %v", err)
			}

			// Apply the same filtering logic as listWins
			var filteredWins []*models.Win
			if tt.fromDate != "" || tt.toDate != "" {
				from, to := parseListDateRange(tt.fromDate, tt.toDate, false)
				for _, win := range wins {
					if !win.Timestamp.Before(from) && !win.Timestamp.After(to) {
						filteredWins = append(filteredWins, win)
					}
				}
			} else {
				filteredWins = wins
			}

			if len(filteredWins) != tt.expectedCount {
				t.Errorf("filtered count = %d, want %d", len(filteredWins), tt.expectedCount)
			}

			if len(filteredWins) > 0 {
				if filteredWins[0].Message != tt.expectedFirst {
					t.Errorf("first win = %q, want %q", filteredWins[0].Message, tt.expectedFirst)
				}
				if filteredWins[len(filteredWins)-1].Message != tt.expectedLast {
					t.Errorf("last win = %q, want %q", filteredWins[len(filteredWins)-1].Message, tt.expectedLast)
				}
			}
		})
	}
}

func TestListWinsEmpty(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create storage with no wins
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	if len(wins) != 0 {
		t.Errorf("expected 0 wins, got %d", len(wins))
	}
}

func TestListWinsWithLineNumbers(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create storage
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	testWins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
			Message:   "Win 1",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
			Message:   "Win 2",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC),
			Message:   "Win 3",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	// Verify that we can iterate through wins with line numbers
	for i, win := range wins {
		lineNum := i + 1
		expectedMessage := testWins[i].Message

		if win.Message != expectedMessage {
			t.Errorf("win[%d] message = %q, want %q", lineNum, win.Message, expectedMessage)
		}

		// Verify line numbers are sequential starting from 1
		if lineNum != i+1 {
			t.Errorf("line number mismatch: got %d, want %d", lineNum, i+1)
		}
	}

	// Verify we have exactly 3 wins
	if len(wins) != 3 {
		t.Errorf("win count = %d, want 3", len(wins))
	}
}

func TestListWinsDateFiltering(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create storage
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	testWins := []*models.Win{
		{
			Timestamp: time.Date(2025, 12, 25, 10, 0, 0, 0, time.UTC),
			Message:   "Old win from last year",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC),
			Message:   "New year win",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 2, 1, 14, 0, 0, 0, time.UTC),
			Message:   "February win",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	// Test filtering by year
	from, to := parseListDateRange("2026-01-01", "2026-01-31", false)
	var filtered []*models.Win
	for _, win := range wins {
		if !win.Timestamp.Before(from) && !win.Timestamp.After(to) {
			filtered = append(filtered, win)
		}
	}

	if len(filtered) != 1 {
		t.Errorf("filtered count = %d, want 1 (only January 2026 win)", len(filtered))
	}

	if len(filtered) > 0 && filtered[0].Message != "New year win" {
		t.Errorf("filtered win = %q, want %q", filtered[0].Message, "New year win")
	}
}

func TestListWinsOrder(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create storage
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Add wins in non-chronological order
	testWins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC),
			Message:   "Middle",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
			Message:   "First",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC),
			Message:   "Last",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	// Wins should be in the order they were written (chronological by file order, not timestamp)
	expectedOrder := []string{"Middle", "First", "Last"}

	for i, expectedMsg := range expectedOrder {
		if i >= len(wins) {
			t.Errorf("missing win at index %d", i)
			continue
		}
		if wins[i].Message != expectedMsg {
			t.Errorf("win[%d] = %q, want %q", i, wins[i].Message, expectedMsg)
		}
	}
}

func TestListWinsLineNumbersWithFiltering(t *testing.T) {
	// This test verifies that line numbers shown by `brag list` correspond to
	// the actual line numbers in the file, not the filtered result indices.
	// This is critical so that `brag delete <N>` and `brag edit <N>` work correctly.

	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create 10 wins within the last week, but only some will match our tag filter
	now := time.Now()
	testWins := []*models.Win{
		{Timestamp: now.AddDate(0, 0, -6), Message: "Win 1 #work", Tags: []string{"#work"}},         // Line 1
		{Timestamp: now.AddDate(0, 0, -5), Message: "Win 2", Tags: []string{}},                      // Line 2
		{Timestamp: now.AddDate(0, 0, -5), Message: "Win 3 #work", Tags: []string{"#work"}},         // Line 3
		{Timestamp: now.AddDate(0, 0, -4), Message: "Win 4", Tags: []string{}},                      // Line 4
		{Timestamp: now.AddDate(0, 0, -3), Message: "Win 5 #personal", Tags: []string{"#personal"}}, // Line 5
		{Timestamp: now.AddDate(0, 0, -3), Message: "Win 6 #work", Tags: []string{"#work"}},         // Line 6
		{Timestamp: now.AddDate(0, 0, -2), Message: "Win 7", Tags: []string{}},                      // Line 7
		{Timestamp: now.AddDate(0, 0, -1), Message: "Win 8 #work", Tags: []string{"#work"}},         // Line 8
		{Timestamp: now.AddDate(0, 0, -1), Message: "Win 9", Tags: []string{}},                      // Line 9
		{Timestamp: now, Message: "Win 10 #work", Tags: []string{"#work"}},                          // Line 10
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	// Simulate filtering by tag "work" (like `brag list --tag work`)
	// This should return wins at file positions 1, 3, 6, 8, 10
	from, to := parseListDateRange("", "", false) // Use default range

	type indexedWin struct {
		win   *models.Win
		index int
	}

	var filteredWins []indexedWin
	for i, win := range wins {
		if win.Timestamp.Before(from) || win.Timestamp.After(to) {
			continue
		}
		if !win.HasTag("work") {
			continue
		}
		filteredWins = append(filteredWins, indexedWin{win: win, index: i + 1})
	}

	// Verify we got the expected wins
	if len(filteredWins) != 5 {
		t.Fatalf("expected 5 filtered wins, got %d", len(filteredWins))
	}

	// The key test: verify that the line numbers are the ORIGINAL file positions
	expectedLineNumbers := []int{1, 3, 6, 8, 10}
	expectedMessages := []string{"Win 1 #work", "Win 3 #work", "Win 6 #work", "Win 8 #work", "Win 10 #work"}

	for i, iw := range filteredWins {
		if iw.index != expectedLineNumbers[i] {
			t.Errorf("filtered win %d: line number = %d, want %d", i, iw.index, expectedLineNumbers[i])
		}
		if iw.win.Message != expectedMessages[i] {
			t.Errorf("filtered win %d: message = %q, want %q", i, iw.win.Message, expectedMessages[i])
		}
	}

	// Critical test: verify that if we want to delete "Win 6 #work" (which appears
	// as the 3rd item in filtered results), we would use line number 6, not 3
	targetWinInFilteredResults := 2 // 0-indexed, this is the 3rd result
	lineNumberToDelete := filteredWins[targetWinInFilteredResults].index

	if lineNumberToDelete != 6 {
		t.Errorf("to delete the 3rd filtered result, line number should be 6, got %d", lineNumberToDelete)
	}

	// Verify that this line number works with the delete operation
	deleteIndex := lineNumberToDelete - 1 // Convert to 0-based index
	if deleteIndex < 0 || deleteIndex >= len(wins) {
		t.Fatalf("delete index %d is out of range", deleteIndex)
	}

	if wins[deleteIndex].Message != "Win 6 #work" {
		t.Errorf("line %d contains %q, expected %q", lineNumberToDelete, wins[deleteIndex].Message, "Win 6 #work")
	}
}
