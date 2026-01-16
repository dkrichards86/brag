package cmd

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
)

func TestAddTag(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create storage and add test wins
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	testWins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
			Message:   "Fixed bug #work",
			Tags:      []string{"#work"},
		},
		{
			Timestamp: time.Date(2026, 1, 8, 15, 0, 0, 0, time.UTC),
			Message:   "Completed feature",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 8, 16, 30, 0, 0, time.UTC),
			Message:   "Code review #review #backend",
			Tags:      []string{"#review", "#backend"},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	tests := []struct {
		name          string
		lineNum       int
		tag           string
		expectedMsg   string
		expectedTags  []string
		shouldSucceed bool
	}{
		{
			name:          "add tag with hash to untagged win",
			lineNum:       2,
			tag:           "#frontend",
			expectedMsg:   "Completed feature #frontend",
			expectedTags:  []string{"#frontend"},
			shouldSucceed: true,
		},
		{
			name:          "add tag without hash to untagged win",
			lineNum:       2,
			tag:           "testing",
			expectedMsg:   "Completed feature #testing",
			expectedTags:  []string{"#testing"},
			shouldSucceed: true,
		},
		{
			name:          "add tag to already tagged win",
			lineNum:       1,
			tag:           "urgent",
			expectedMsg:   "Fixed bug #work #urgent",
			expectedTags:  []string{"#work", "#urgent"},
			shouldSucceed: true,
		},
		{
			name:          "add multiple tags to win with existing tags",
			lineNum:       3,
			tag:           "#urgent",
			expectedMsg:   "Code review #review #backend #urgent",
			expectedTags:  []string{"#review", "#backend", "#urgent"},
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Re-write test wins to reset state for each test
			if err := store.WriteWins(testWins); err != nil {
				t.Fatalf("Failed to reset test wins: %v", err)
			}

			// Prepare the tag
			tag := tt.tag
			if !strings.HasPrefix(tag, "#") {
				tag = "#" + tag
			}

			// Read wins
			wins, err := store.ReadAllWins()
			if err != nil {
				t.Fatalf("Failed to read wins: %v", err)
			}

			index := tt.lineNum - 1
			if index < 0 || index >= len(wins) {
				t.Fatalf("Invalid line number: %d", tt.lineNum)
			}

			win := wins[index]

			// Check if tag already exists (should not add duplicate)
			if win.HasTag(tag) {
				return
			}

			// Add the tag
			newMessage := win.Message + " " + tag
			if err := store.UpdateWin(index, newMessage); err != nil {
				if tt.shouldSucceed {
					t.Errorf("Failed to update win: %v", err)
				}
				return
			}

			if !tt.shouldSucceed {
				t.Error("Expected operation to fail, but it succeeded")
				return
			}

			// Verify the update
			wins, err = store.ReadAllWins()
			if err != nil {
				t.Fatalf("Failed to read wins after update: %v", err)
			}

			updatedWin := wins[index]
			if updatedWin.Message != tt.expectedMsg {
				t.Errorf("Updated message = %q, want %q", updatedWin.Message, tt.expectedMsg)
			}

			if len(updatedWin.Tags) != len(tt.expectedTags) {
				t.Errorf("Number of tags = %d, want %d", len(updatedWin.Tags), len(tt.expectedTags))
			}

			for i, expectedTag := range tt.expectedTags {
				if i >= len(updatedWin.Tags) || updatedWin.Tags[i] != expectedTag {
					t.Errorf("Tag[%d] = %q, want %q", i, updatedWin.Tags[i], expectedTag)
				}
			}
		})
	}
}

func TestAddTagDuplicate(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create storage and add test win
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	testWins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
			Message:   "Fixed bug #work #urgent",
			Tags:      []string{"#work", "#urgent"},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	// Read the win
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	win := wins[0]

	// Try to add a tag that already exists
	if !win.HasTag("#work") {
		t.Error("Expected win to have #work tag")
	}

	// Verify that HasTag works case-insensitively
	if !win.HasTag("WORK") {
		t.Error("Expected HasTag to work case-insensitively")
	}

	if !win.HasTag("work") {
		t.Error("Expected HasTag to work without # prefix")
	}
}

func TestAddMultipleTags(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create storage and add test win
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	testWins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
			Message:   "Completed feature",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	// Add multiple tags at once
	tags := []string{"#frontend", "#testing", "#urgent"}
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	newMessage := wins[0].Message + " " + strings.Join(tags, " ")
	if err := store.UpdateWin(0, newMessage); err != nil {
		t.Fatalf("Failed to update win: %v", err)
	}

	// Verify all tags were added
	wins, err = store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins after update: %v", err)
	}

	updatedWin := wins[0]
	expectedMsg := "Completed feature #frontend #testing #urgent"
	if updatedWin.Message != expectedMsg {
		t.Errorf("Updated message = %q, want %q", updatedWin.Message, expectedMsg)
	}

	if len(updatedWin.Tags) != 3 {
		t.Errorf("Number of tags = %d, want 3", len(updatedWin.Tags))
	}

	for _, tag := range tags {
		if !updatedWin.HasTag(tag) {
			t.Errorf("Expected win to have tag %s", tag)
		}
	}
}
