package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
)

func TestDeleteWin(t *testing.T) {
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
			Timestamp: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
			Message:   "Second win #personal",
			Tags:      []string{"#personal"},
		},
		{
			Timestamp: time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC),
			Message:   "Third win",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	tests := []struct {
		name               string
		index              int
		expectedCountAfter int
		expectedRemaining  []string
		shouldSucceed      bool
	}{
		{
			name:               "delete first win",
			index:              0,
			expectedCountAfter: 2,
			expectedRemaining:  []string{"Second win #personal", "Third win"},
			shouldSucceed:      true,
		},
		{
			name:               "delete middle win",
			index:              1,
			expectedCountAfter: 2,
			expectedRemaining:  []string{"First win #work", "Third win"},
			shouldSucceed:      true,
		},
		{
			name:               "delete last win",
			index:              2,
			expectedCountAfter: 2,
			expectedRemaining:  []string{"First win #work", "Second win #personal"},
			shouldSucceed:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset test wins
			if err := store.WriteWins(testWins); err != nil {
				t.Fatalf("Failed to reset test wins: %v", err)
			}

			// Delete the win
			err := store.DeleteWin(tt.index)

			if tt.shouldSucceed && err != nil {
				t.Errorf("DeleteWin() error = %v, expected success", err)
				return
			}

			if !tt.shouldSucceed && err == nil {
				t.Error("DeleteWin() expected error, got nil")
				return
			}

			// Verify the deletion
			wins, err := store.ReadAllWins()
			if err != nil {
				t.Fatalf("Failed to read wins: %v", err)
			}

			if len(wins) != tt.expectedCountAfter {
				t.Errorf("win count = %d, want %d", len(wins), tt.expectedCountAfter)
			}

			for i, expectedMsg := range tt.expectedRemaining {
				if i >= len(wins) {
					t.Errorf("missing expected win: %q", expectedMsg)
					continue
				}
				if wins[i].Message != expectedMsg {
					t.Errorf("win[%d] = %q, want %q", i, wins[i].Message, expectedMsg)
				}
			}
		})
	}
}

func TestDeleteWinInvalidIndex(t *testing.T) {
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
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	tests := []struct {
		name  string
		index int
	}{
		{
			name:  "negative index",
			index: -1,
		},
		{
			name:  "index too large",
			index: 10,
		},
		{
			name:  "index equals length",
			index: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.DeleteWin(tt.index)

			if err == nil {
				t.Error("Expected error for invalid index, got nil")
			}

			// Verify no wins were deleted
			wins, err := store.ReadAllWins()
			if err != nil {
				t.Fatalf("Failed to read wins: %v", err)
			}

			if len(wins) != 2 {
				t.Errorf("win count = %d, want 2 (no deletion should occur)", len(wins))
			}
		})
	}
}

func TestDeleteAllWins(t *testing.T) {
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

	// Delete wins one by one
	err = store.DeleteWin(2) // Delete last
	if err != nil {
		t.Fatalf("DeleteWin(2) error = %v", err)
	}

	err = store.DeleteWin(1) // Delete new last
	if err != nil {
		t.Fatalf("DeleteWin(1) error = %v", err)
	}

	err = store.DeleteWin(0) // Delete remaining
	if err != nil {
		t.Fatalf("DeleteWin(0) error = %v", err)
	}

	// Verify all wins are deleted
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	if len(wins) != 0 {
		t.Errorf("win count = %d, want 0 (all wins should be deleted)", len(wins))
	}
}

func TestDeleteWinFromEmpty(t *testing.T) {
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

	// Try to delete from empty list
	err = store.DeleteWin(0)

	if err == nil {
		t.Error("Expected error when deleting from empty list, got nil")
	}
}

func TestDeleteWinMultipleTimes(t *testing.T) {
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
		{
			Timestamp: time.Date(2026, 1, 4, 13, 0, 0, 0, time.UTC),
			Message:   "Win 4",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 5, 14, 0, 0, 0, time.UTC),
			Message:   "Win 5",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	// Delete multiple wins (always delete index 1 to test shifting)
	deletions := []struct {
		index           int
		expectedRemaining []string
	}{
		{
			index:           1, // Delete "Win 2"
			expectedRemaining: []string{"Win 1", "Win 3", "Win 4", "Win 5"},
		},
		{
			index:           1, // Delete "Win 3" (now at index 1)
			expectedRemaining: []string{"Win 1", "Win 4", "Win 5"},
		},
		{
			index:           1, // Delete "Win 4" (now at index 1)
			expectedRemaining: []string{"Win 1", "Win 5"},
		},
	}

	for i, d := range deletions {
		err := store.DeleteWin(d.index)
		if err != nil {
			t.Fatalf("DeleteWin() iteration %d error = %v", i, err)
		}

		wins, err := store.ReadAllWins()
		if err != nil {
			t.Fatalf("Failed to read wins: %v", err)
		}

		if len(wins) != len(d.expectedRemaining) {
			t.Errorf("iteration %d: win count = %d, want %d", i, len(wins), len(d.expectedRemaining))
		}

		for j, expectedMsg := range d.expectedRemaining {
			if j >= len(wins) {
				t.Errorf("iteration %d: missing expected win: %q", i, expectedMsg)
				continue
			}
			if wins[j].Message != expectedMsg {
				t.Errorf("iteration %d: win[%d] = %q, want %q", i, j, wins[j].Message, expectedMsg)
			}
		}
	}
}

func TestDeleteWinPreservesOrder(t *testing.T) {
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
			Message:   "First",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
			Message:   "Second - to be deleted",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC),
			Message:   "Third",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 4, 13, 0, 0, 0, time.UTC),
			Message:   "Fourth",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	// Delete the second win
	err = store.DeleteWin(1)
	if err != nil {
		t.Fatalf("DeleteWin() error = %v", err)
	}

	// Verify order is preserved
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	expectedOrder := []string{"First", "Third", "Fourth"}

	if len(wins) != len(expectedOrder) {
		t.Errorf("win count = %d, want %d", len(wins), len(expectedOrder))
	}

	for i, expectedMsg := range expectedOrder {
		if i >= len(wins) {
			t.Errorf("missing expected win: %q", expectedMsg)
			continue
		}
		if wins[i].Message != expectedMsg {
			t.Errorf("win[%d] = %q, want %q (order not preserved)", i, wins[i].Message, expectedMsg)
		}
	}
}
