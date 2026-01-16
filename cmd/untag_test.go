package cmd

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
)

func TestRemoveTag(t *testing.T) {
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
			Message:   "Fixed bug #work #urgent #backend",
			Tags:      []string{"#work", "#urgent", "#backend"},
		},
		{
			Timestamp: time.Date(2026, 1, 8, 15, 0, 0, 0, time.UTC),
			Message:   "Completed feature #frontend",
			Tags:      []string{"#frontend"},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	tests := []struct {
		name         string
		lineNum      int
		tag          string
		expectedMsg  string
		expectedTags []string
	}{
		{
			name:         "remove one tag from multiple tags",
			lineNum:      1,
			tag:          "#urgent",
			expectedMsg:  "Fixed bug #work #backend",
			expectedTags: []string{"#work", "#backend"},
		},
		{
			name:         "remove only tag",
			lineNum:      2,
			tag:          "frontend",
			expectedMsg:  "Completed feature",
			expectedTags: []string{},
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
			win := wins[index]

			// Remove the tag
			newMessage := win.Message
			newMessage = strings.ReplaceAll(newMessage, " "+tag, "")
			newMessage = strings.ReplaceAll(newMessage, tag, "")
			newMessage = strings.Join(strings.Fields(newMessage), " ")

			if err := store.UpdateWin(index, newMessage); err != nil {
				t.Fatalf("Failed to update win: %v", err)
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

			for _, expectedTag := range tt.expectedTags {
				if !updatedWin.HasTag(expectedTag) {
					t.Errorf("Expected win to still have tag %s", expectedTag)
				}
			}

			// Verify the tag was actually removed
			if updatedWin.HasTag(tag) {
				t.Errorf("Expected tag %s to be removed", tag)
			}
		})
	}
}

func TestRemoveMultipleTags(t *testing.T) {
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
			Message:   "Fixed bug #work #urgent #backend #testing",
			Tags:      []string{"#work", "#urgent", "#backend", "#testing"},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	// Remove multiple tags at once
	tagsToRemove := []string{"#urgent", "#testing"}
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	newMessage := wins[0].Message
	for _, tag := range tagsToRemove {
		newMessage = strings.ReplaceAll(newMessage, " "+tag, "")
		newMessage = strings.ReplaceAll(newMessage, tag, "")
	}
	newMessage = strings.Join(strings.Fields(newMessage), " ")

	if err := store.UpdateWin(0, newMessage); err != nil {
		t.Fatalf("Failed to update win: %v", err)
	}

	// Verify tags were removed
	wins, err = store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins after update: %v", err)
	}

	updatedWin := wins[0]
	expectedMsg := "Fixed bug #work #backend"
	if updatedWin.Message != expectedMsg {
		t.Errorf("Updated message = %q, want %q", updatedWin.Message, expectedMsg)
	}

	if len(updatedWin.Tags) != 2 {
		t.Errorf("Number of tags = %d, want 2", len(updatedWin.Tags))
	}

	// Verify correct tags remain
	if !updatedWin.HasTag("#work") || !updatedWin.HasTag("#backend") {
		t.Error("Expected #work and #backend tags to remain")
	}

	// Verify removed tags are gone
	for _, tag := range tagsToRemove {
		if updatedWin.HasTag(tag) {
			t.Errorf("Expected tag %s to be removed", tag)
		}
	}
}
