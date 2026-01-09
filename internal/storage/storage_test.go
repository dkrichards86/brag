package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/models"
)

func TestNew(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Logf("Warning: failed to restore HOME: %v", err)
		}
	}()

	// Set HOME to temp directory
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("Failed to set HOME: %v", err)
	}

	storage, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, ".brag", "wins.txt")
	if storage.GetFilePath() != expectedPath {
		t.Errorf("GetFilePath() = %v, want %v", storage.GetFilePath(), expectedPath)
	}

	// Check directory was created
	if _, err := os.Stat(filepath.Join(tmpDir, ".brag")); os.IsNotExist(err) {
		t.Error("Expected .brag directory to be created")
	}

	// Check file was created
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Error("Expected wins.txt file to be created")
	}
}

func TestAddWin(t *testing.T) {
	storage := setupTestStorage(t)

	timestamp := time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC)
	message := "Test win #test"

	err := storage.AddWin(message, timestamp)
	if err != nil {
		t.Fatalf("AddWin() error = %v", err)
	}

	// Read the file and verify content
	content, err := os.ReadFile(storage.GetFilePath())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expected := "2026-01-08 14:32 | Test win #test\n"
	if string(content) != expected {
		t.Errorf("File content = %q, want %q", string(content), expected)
	}
}

func TestAddMultipleWins(t *testing.T) {
	storage := setupTestStorage(t)

	wins := []struct {
		message   string
		timestamp time.Time
	}{
		{"First win #work", time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC)},
		{"Second win #personal", time.Date(2026, 1, 8, 14, 30, 0, 0, time.UTC)},
		{"Third win #project", time.Date(2026, 1, 8, 16, 45, 0, 0, time.UTC)},
	}

	for _, w := range wins {
		if err := storage.AddWin(w.message, w.timestamp); err != nil {
			t.Fatalf("AddWin() error = %v", err)
		}
	}

	// Read all wins back
	allWins, err := storage.ReadAllWins()
	if err != nil {
		t.Fatalf("ReadAllWins() error = %v", err)
	}

	if len(allWins) != len(wins) {
		t.Errorf("ReadAllWins() returned %d wins, want %d", len(allWins), len(wins))
	}

	for i, got := range allWins {
		if got.Message != wins[i].message {
			t.Errorf("Win %d message = %v, want %v", i, got.Message, wins[i].message)
		}
		if !got.Timestamp.Equal(wins[i].timestamp) {
			t.Errorf("Win %d timestamp = %v, want %v", i, got.Timestamp, wins[i].timestamp)
		}
	}
}

func TestReadAllWins_EmptyFile(t *testing.T) {
	storage := setupTestStorage(t)

	wins, err := storage.ReadAllWins()
	if err != nil {
		t.Fatalf("ReadAllWins() error = %v", err)
	}

	if len(wins) != 0 {
		t.Errorf("ReadAllWins() returned %d wins, want 0", len(wins))
	}
}

func TestReadAllWins_SkipsInvalidLines(t *testing.T) {
	storage := setupTestStorage(t)

	// Write a mix of valid and invalid lines
	content := `2026-01-08 10:00 | Valid win #test
invalid line without separator
2026-01-08 14:30 | Another valid win
bad-timestamp | Some message
2026-01-08 16:00 | Third valid win #work
`
	if err := os.WriteFile(storage.GetFilePath(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	wins, err := storage.ReadAllWins()
	if err != nil {
		t.Fatalf("ReadAllWins() error = %v", err)
	}

	// Should have 3 valid wins
	if len(wins) != 3 {
		t.Errorf("ReadAllWins() returned %d wins, want 3", len(wins))
	}

	expectedMessages := []string{
		"Valid win #test",
		"Another valid win",
		"Third valid win #work",
	}

	for i, win := range wins {
		if win.Message != expectedMessages[i] {
			t.Errorf("Win %d message = %v, want %v", i, win.Message, expectedMessages[i])
		}
	}
}

func TestReadAllWins_SkipsBlankLines(t *testing.T) {
	storage := setupTestStorage(t)

	content := `2026-01-08 10:00 | First win

2026-01-08 14:30 | Second win

2026-01-08 16:00 | Third win
`
	if err := os.WriteFile(storage.GetFilePath(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	wins, err := storage.ReadAllWins()
	if err != nil {
		t.Fatalf("ReadAllWins() error = %v", err)
	}

	if len(wins) != 3 {
		t.Errorf("ReadAllWins() returned %d wins, want 3", len(wins))
	}
}

func TestWriteWins(t *testing.T) {
	storage := setupTestStorage(t)

	wins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC),
			Message:   "First win #test",
			Tags:      []string{"#test"},
		},
		{
			Timestamp: time.Date(2026, 1, 8, 14, 30, 0, 0, time.UTC),
			Message:   "Second win #work",
			Tags:      []string{"#work"},
		},
	}

	err := storage.WriteWins(wins)
	if err != nil {
		t.Fatalf("WriteWins() error = %v", err)
	}

	// Read back and verify
	readWins, err := storage.ReadAllWins()
	if err != nil {
		t.Fatalf("ReadAllWins() error = %v", err)
	}

	if len(readWins) != len(wins) {
		t.Errorf("ReadAllWins() returned %d wins, want %d", len(readWins), len(wins))
	}

	for i, got := range readWins {
		if got.Message != wins[i].Message {
			t.Errorf("Win %d message = %v, want %v", i, got.Message, wins[i].Message)
		}
	}
}

func TestDeleteWin(t *testing.T) {
	storage := setupTestStorage(t)

	// Add some wins
	wins := []struct {
		message   string
		timestamp time.Time
	}{
		{"First win", time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC)},
		{"Second win", time.Date(2026, 1, 8, 11, 0, 0, 0, time.UTC)},
		{"Third win", time.Date(2026, 1, 8, 12, 0, 0, 0, time.UTC)},
	}

	for _, w := range wins {
		if err := storage.AddWin(w.message, w.timestamp); err != nil {
			t.Fatalf("AddWin() error = %v", err)
		}
	}

	// Delete the middle win (index 1)
	err := storage.DeleteWin(1)
	if err != nil {
		t.Fatalf("DeleteWin() error = %v", err)
	}

	// Verify it was deleted
	allWins, err := storage.ReadAllWins()
	if err != nil {
		t.Fatalf("ReadAllWins() error = %v", err)
	}

	if len(allWins) != 2 {
		t.Errorf("After delete, got %d wins, want 2", len(allWins))
	}

	if allWins[0].Message != "First win" {
		t.Errorf("First win message = %v, want 'First win'", allWins[0].Message)
	}
	if allWins[1].Message != "Third win" {
		t.Errorf("Second win message = %v, want 'Third win'", allWins[1].Message)
	}
}

func TestDeleteWin_InvalidIndex(t *testing.T) {
	storage := setupTestStorage(t)

	// Add one win
	if err := storage.AddWin("Test win", time.Now()); err != nil {
		t.Fatalf("AddWin() error = %v", err)
	}

	tests := []struct {
		name  string
		index int
	}{
		{"negative index", -1},
		{"index too large", 5},
		{"index equals length", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.DeleteWin(tt.index)
			if err == nil {
				t.Error("DeleteWin() expected error, got nil")
			}
		})
	}
}

func TestUpdateWin(t *testing.T) {
	storage := setupTestStorage(t)

	// Add some wins
	originalTime := time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC)
	if err := storage.AddWin("Original message", originalTime); err != nil {
		t.Fatalf("AddWin() error = %v", err)
	}

	// Update the win
	newMessage := "Updated message #updated"
	err := storage.UpdateWin(0, newMessage)
	if err != nil {
		t.Fatalf("UpdateWin() error = %v", err)
	}

	// Verify the update
	wins, err := storage.ReadAllWins()
	if err != nil {
		t.Fatalf("ReadAllWins() error = %v", err)
	}

	if len(wins) != 1 {
		t.Fatalf("Expected 1 win, got %d", len(wins))
	}

	if wins[0].Message != newMessage {
		t.Errorf("Updated message = %v, want %v", wins[0].Message, newMessage)
	}

	// Verify timestamp was preserved
	if !wins[0].Timestamp.Equal(originalTime) {
		t.Errorf("Timestamp changed from %v to %v", originalTime, wins[0].Timestamp)
	}
}

func TestUpdateWin_InvalidIndex(t *testing.T) {
	storage := setupTestStorage(t)

	// Add one win
	if err := storage.AddWin("Test win", time.Now()); err != nil {
		t.Fatalf("AddWin() error = %v", err)
	}

	tests := []struct {
		name  string
		index int
	}{
		{"negative index", -1},
		{"index too large", 5},
		{"index equals length", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.UpdateWin(tt.index, "New message")
			if err == nil {
				t.Error("UpdateWin() expected error, got nil")
			}
		})
	}
}

// setupTestStorage creates a Storage instance in a temporary directory
func setupTestStorage(t *testing.T) *Storage {
	t.Helper()

	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Logf("Warning: failed to restore HOME: %v", err)
		}
	})

	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("Failed to set HOME: %v", err)
	}

	storage, err := New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	return storage
}
