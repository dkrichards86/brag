package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
)

func TestEditWinUpdate(t *testing.T) {
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
			Message:   "Original message #work",
			Tags:      []string{"#work"},
		},
		{
			Timestamp: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
			Message:   "Second win",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC),
			Message:   "Third win #personal",
			Tags:      []string{"#personal"},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	tests := []struct {
		name           string
		index          int
		newMessage     string
		expectedMsg    string
		expectedTags   []string
		shouldSucceed  bool
		expectedErrMsg string
	}{
		{
			name:          "update first win",
			index:         0,
			newMessage:    "Updated first win #updated",
			expectedMsg:   "Updated first win #updated",
			expectedTags:  []string{"#updated"},
			shouldSucceed: true,
		},
		{
			name:          "update middle win",
			index:         1,
			newMessage:    "Updated second win #new #tags",
			expectedMsg:   "Updated second win #new #tags",
			expectedTags:  []string{"#new", "#tags"},
			shouldSucceed: true,
		},
		{
			name:          "update last win",
			index:         2,
			newMessage:    "Updated third win",
			expectedMsg:   "Updated third win",
			expectedTags:  []string{},
			shouldSucceed: true,
		},
		{
			name:          "remove tags from win",
			index:         0,
			newMessage:    "Win without tags",
			expectedMsg:   "Win without tags",
			expectedTags:  []string{},
			shouldSucceed: true,
		},
		{
			name:          "add multiple tags",
			index:         1,
			newMessage:    "Win with many tags #tag1 #tag2 #tag3 #tag4",
			expectedMsg:   "Win with many tags #tag1 #tag2 #tag3 #tag4",
			expectedTags:  []string{"#tag1", "#tag2", "#tag3", "#tag4"},
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset test wins
			if err := store.WriteWins(testWins); err != nil {
				t.Fatalf("Failed to reset test wins: %v", err)
			}

			// Update the win
			err := store.UpdateWin(tt.index, tt.newMessage)

			if tt.shouldSucceed && err != nil {
				t.Errorf("UpdateWin() error = %v, expected success", err)
				return
			}

			if !tt.shouldSucceed {
				if err == nil {
					t.Error("UpdateWin() expected error, got nil")
				}
				return
			}

			// Verify the update
			wins, err := store.ReadAllWins()
			if err != nil {
				t.Fatalf("Failed to read wins: %v", err)
			}

			if tt.index >= len(wins) {
				t.Fatalf("Index %d out of range", tt.index)
			}

			updatedWin := wins[tt.index]

			if updatedWin.Message != tt.expectedMsg {
				t.Errorf("message = %q, want %q", updatedWin.Message, tt.expectedMsg)
			}

			if len(updatedWin.Tags) != len(tt.expectedTags) {
				t.Errorf("tags count = %d, want %d", len(updatedWin.Tags), len(tt.expectedTags))
			}

			for i, expectedTag := range tt.expectedTags {
				if i >= len(updatedWin.Tags) || updatedWin.Tags[i] != expectedTag {
					t.Errorf("tag[%d] = %q, want %q", i, updatedWin.Tags[i], expectedTag)
				}
			}
		})
	}
}

func TestEditWinInvalidIndex(t *testing.T) {
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
			err := store.UpdateWin(tt.index, "Updated message")

			if err == nil {
				t.Error("Expected error for invalid index, got nil")
			}
		})
	}
}

func TestEditWinPreservesTimestamp(t *testing.T) {
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

	// Note: File format only stores minute precision, so use a timestamp without seconds
	originalTimestamp := time.Date(2026, 1, 5, 14, 30, 0, 0, time.UTC)

	testWins := []*models.Win{
		{
			Timestamp: originalTimestamp,
			Message:   "Original message",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	// Update the win
	err = store.UpdateWin(0, "Updated message")
	if err != nil {
		t.Fatalf("UpdateWin() error = %v", err)
	}

	// Verify timestamp is preserved
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	if len(wins) != 1 {
		t.Fatalf("Expected 1 win, got %d", len(wins))
	}

	if !wins[0].Timestamp.Equal(originalTimestamp) {
		t.Errorf("timestamp = %v, want %v (timestamp should be preserved)", wins[0].Timestamp, originalTimestamp)
	}

	if wins[0].Message != "Updated message" {
		t.Errorf("message = %q, want %q", wins[0].Message, "Updated message")
	}
}

func TestEditWinMultipleUpdates(t *testing.T) {
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
			Message:   "Original",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	// Update multiple times
	updates := []string{
		"First update",
		"Second update #tag1",
		"Third update #tag1 #tag2",
		"Final update",
	}

	for i, newMessage := range updates {
		err := store.UpdateWin(0, newMessage)
		if err != nil {
			t.Fatalf("UpdateWin() iteration %d error = %v", i, err)
		}

		// Verify the update
		wins, err := store.ReadAllWins()
		if err != nil {
			t.Fatalf("Failed to read wins: %v", err)
		}

		if wins[0].Message != newMessage {
			t.Errorf("iteration %d: message = %q, want %q", i, wins[0].Message, newMessage)
		}
	}

	// Verify final state
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	if wins[0].Message != "Final update" {
		t.Errorf("final message = %q, want %q", wins[0].Message, "Final update")
	}
}

func TestEditWinDoesNotAffectOtherWins(t *testing.T) {
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
			Message:   "Win 2 - will be updated",
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

	// Update the middle win
	err = store.UpdateWin(1, "Updated Win 2")
	if err != nil {
		t.Fatalf("UpdateWin() error = %v", err)
	}

	// Verify other wins are unchanged
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	if wins[0].Message != "Win 1" {
		t.Errorf("win[0] message changed: %q, want %q", wins[0].Message, "Win 1")
	}

	if wins[1].Message != "Updated Win 2" {
		t.Errorf("win[1] message = %q, want %q", wins[1].Message, "Updated Win 2")
	}

	if wins[2].Message != "Win 3" {
		t.Errorf("win[2] message changed: %q, want %q", wins[2].Message, "Win 3")
	}
}
