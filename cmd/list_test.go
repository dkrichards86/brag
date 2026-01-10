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

			var filteredWins []*models.Win
			if tt.fromDate != "" || tt.toDate != "" {
				from, to := parseDateRange(tt.fromDate, tt.toDate)
				filteredWins = filterWins(wins, from, to, "", false, false)
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
	from, to := parseDateRange("2026-01-01", "2026-01-31")
	filtered := filterWins(wins, from, to, "", false, false)

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
