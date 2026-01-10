package cmd

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
)

func TestSearchWins(t *testing.T) {
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
			Message:   "Fixed critical bug in authentication #work #backend",
			Tags:      []string{"#work", "#backend"},
		},
		{
			Timestamp: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
			Message:   "Completed feature for user profile #work #frontend",
			Tags:      []string{"#work", "#frontend"},
		},
		{
			Timestamp: time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC),
			Message:   "Had a great workout session #personal #health",
			Tags:      []string{"#personal", "#health"},
		},
		{
			Timestamp: time.Date(2026, 1, 4, 13, 0, 0, 0, time.UTC),
			Message:   "Fixed typo in documentation",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 5, 14, 0, 0, 0, time.UTC),
			Message:   "Deployed new backend service to production #work #backend #deployment",
			Tags:      []string{"#work", "#backend", "#deployment"},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	tests := []struct {
		name          string
		query         string
		expectedCount int
		expectedMsgs  []string
	}{
		{
			name:          "search for 'bug'",
			query:         "bug",
			expectedCount: 1,
			expectedMsgs:  []string{"Fixed critical bug in authentication #work #backend"},
		},
		{
			name:          "search for 'fixed'",
			query:         "fixed",
			expectedCount: 2,
			expectedMsgs:  []string{"Fixed critical bug in authentication #work #backend", "Fixed typo in documentation"},
		},
		{
			name:          "search for 'work' (tag and partial)",
			query:         "work",
			expectedCount: 4,
			expectedMsgs: []string{
				"Fixed critical bug in authentication #work #backend",
				"Completed feature for user profile #work #frontend",
				"Had a great workout session #personal #health", // Contains "work" in "workout"
				"Deployed new backend service to production #work #backend #deployment",
			},
		},
		{
			name:          "search for 'backend'",
			query:         "backend",
			expectedCount: 2,
			expectedMsgs: []string{
				"Fixed critical bug in authentication #work #backend",
				"Deployed new backend service to production #work #backend #deployment",
			},
		},
		{
			name:          "search for 'personal'",
			query:         "personal",
			expectedCount: 1,
			expectedMsgs:  []string{"Had a great workout session #personal #health"},
		},
		{
			name:          "case insensitive search",
			query:         "FIXED",
			expectedCount: 2,
			expectedMsgs:  []string{"Fixed critical bug in authentication #work #backend", "Fixed typo in documentation"},
		},
		{
			name:          "partial word match",
			query:         "auth",
			expectedCount: 1,
			expectedMsgs:  []string{"Fixed critical bug in authentication #work #backend"},
		},
		{
			name:          "no matches",
			query:         "nonexistent",
			expectedCount: 0,
			expectedMsgs:  []string{},
		},
		{
			name:          "search with hashtag",
			query:         "#backend",
			expectedCount: 2,
			expectedMsgs: []string{
				"Fixed critical bug in authentication #work #backend",
				"Deployed new backend service to production #work #backend #deployment",
			},
		},
		{
			name:          "multi-word search",
			query:         "user profile",
			expectedCount: 1,
			expectedMsgs:  []string{"Completed feature for user profile #work #frontend"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wins, err := store.ReadAllWins()
			if err != nil {
				t.Fatalf("Failed to read wins: %v", err)
			}

			// Perform search (case-insensitive)
			queryLower := strings.ToLower(tt.query)
			var matches []*models.Win

			for _, win := range wins {
				if strings.Contains(strings.ToLower(win.Message), queryLower) {
					matches = append(matches, win)
				}
			}

			if len(matches) != tt.expectedCount {
				t.Errorf("match count = %d, want %d", len(matches), tt.expectedCount)
			}

			for i, expectedMsg := range tt.expectedMsgs {
				if i >= len(matches) {
					t.Errorf("missing expected match: %q", expectedMsg)
					continue
				}
				if matches[i].Message != expectedMsg {
					t.Errorf("match[%d] = %q, want %q", i, matches[i].Message, expectedMsg)
				}
			}
		})
	}
}

func TestSearchWinsEmpty(t *testing.T) {
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

	// Search in empty list
	query := "anything"
	queryLower := strings.ToLower(query)
	var matches []*models.Win

	for _, win := range wins {
		if strings.Contains(strings.ToLower(win.Message), queryLower) {
			matches = append(matches, win)
		}
	}

	if len(matches) != 0 {
		t.Errorf("match count = %d, want 0 (searching in empty list)", len(matches))
	}
}

func TestSearchWinsCaseInsensitive(t *testing.T) {
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
			Message:   "Fixed BUG in Authentication",
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

	// Test various case combinations
	queries := []string{"bug", "BUG", "Bug", "authentication", "AUTHENTICATION", "Authentication"}

	for _, query := range queries {
		t.Run("query="+query, func(t *testing.T) {
			queryLower := strings.ToLower(query)
			var matches []*models.Win

			for _, win := range wins {
				if strings.Contains(strings.ToLower(win.Message), queryLower) {
					matches = append(matches, win)
				}
			}

			if len(matches) != 1 {
				t.Errorf("query %q: match count = %d, want 1", query, len(matches))
			}
		})
	}
}

func TestSearchWinsSpecialCharacters(t *testing.T) {
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
			Message:   "Fixed bug: null pointer exception",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
			Message:   "Updated README.md documentation",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC),
			Message:   "Merged PR #123 into main branch",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	tests := []struct {
		name          string
		query         string
		expectedCount int
		expectedMsg   string
	}{
		{
			name:          "search with colon",
			query:         "bug:",
			expectedCount: 1,
			expectedMsg:   "Fixed bug: null pointer exception",
		},
		{
			name:          "search with dot",
			query:         "README.md",
			expectedCount: 1,
			expectedMsg:   "Updated README.md documentation",
		},
		{
			name:          "search with hash",
			query:         "#123",
			expectedCount: 1,
			expectedMsg:   "Merged PR #123 into main branch",
		},
		{
			name:          "search for 'PR'",
			query:         "PR",
			expectedCount: 1,
			expectedMsg:   "Merged PR #123 into main branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wins, err := store.ReadAllWins()
			if err != nil {
				t.Fatalf("Failed to read wins: %v", err)
			}

			queryLower := strings.ToLower(tt.query)
			var matches []*models.Win

			for _, win := range wins {
				if strings.Contains(strings.ToLower(win.Message), queryLower) {
					matches = append(matches, win)
				}
			}

			if len(matches) != tt.expectedCount {
				t.Errorf("match count = %d, want %d", len(matches), tt.expectedCount)
			}

			if len(matches) > 0 && matches[0].Message != tt.expectedMsg {
				t.Errorf("matched message = %q, want %q", matches[0].Message, tt.expectedMsg)
			}
		})
	}
}

func TestSearchWinsPartialMatches(t *testing.T) {
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
			Message:   "Implemented authentication",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
			Message:   "Added authorization checks",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC),
			Message:   "Automated testing framework",
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

	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{
			name:          "prefix match 'auth'",
			query:         "auth",
			expectedCount: 2, // authentication and authorization
		},
		{
			name:          "prefix match 'auto'",
			query:         "auto",
			expectedCount: 1, // automated only
		},
		{
			name:          "full word 'authentication'",
			query:         "authentication",
			expectedCount: 1,
		},
		{
			name:          "partial word 'tion'",
			query:         "tion",
			expectedCount: 2, // authentication and authorization
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryLower := strings.ToLower(tt.query)
			var matches []*models.Win

			for _, win := range wins {
				if strings.Contains(strings.ToLower(win.Message), queryLower) {
					matches = append(matches, win)
				}
			}

			if len(matches) != tt.expectedCount {
				t.Errorf("match count = %d, want %d", len(matches), tt.expectedCount)
			}
		})
	}
}

func TestSearchWinsMultipleWords(t *testing.T) {
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
			Message:   "Fixed critical bug in user authentication",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
			Message:   "Updated user profile page",
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

	// Multi-word queries should match as a phrase
	query := "user authentication"
	queryLower := strings.ToLower(query)
	var matches []*models.Win

	for _, win := range wins {
		if strings.Contains(strings.ToLower(win.Message), queryLower) {
			matches = append(matches, win)
		}
	}

	// Should only match the first win (exact phrase match)
	if len(matches) != 1 {
		t.Errorf("match count = %d, want 1 (exact phrase match)", len(matches))
	}

	if len(matches) > 0 && matches[0].Message != "Fixed critical bug in user authentication" {
		t.Errorf("matched wrong win: %q", matches[0].Message)
	}
}
