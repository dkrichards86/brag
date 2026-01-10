package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/storage"
)

func TestParseAddArgs(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedMsg   string
		expectedYear  int
		expectedMonth time.Month
		expectedDay   int
	}{
		{
			name:          "message only",
			args:          []string{"Fixed", "bug", "#work"},
			expectedMsg:   "Fixed bug #work",
			expectedYear:  time.Now().Year(),
			expectedMonth: time.Now().Month(),
			expectedDay:   time.Now().Day(),
		},
		{
			name:          "message with date YYYY/MM/DD",
			args:          []string{"Fixed", "bug", "2026/01/08"},
			expectedMsg:   "Fixed bug",
			expectedYear:  2026,
			expectedMonth: time.January,
			expectedDay:   8,
		},
		{
			name:          "message with date YYYY-MM-DD",
			args:          []string{"Completed", "feature", "2026-01-10"},
			expectedMsg:   "Completed feature",
			expectedYear:  2026,
			expectedMonth: time.January,
			expectedDay:   10,
		},
		{
			name:          "message with date MM/DD/YYYY",
			args:          []string{"Code", "review", "01/15/2026"},
			expectedMsg:   "Code review",
			expectedYear:  2026,
			expectedMonth: time.January,
			expectedDay:   15,
		},
		{
			name:          "message with date MM-DD-YYYY",
			args:          []string{"Deployed", "app", "01-20-2026"},
			expectedMsg:   "Deployed app",
			expectedYear:  2026,
			expectedMonth: time.January,
			expectedDay:   20,
		},
		{
			name:          "single word message",
			args:          []string{"Done!"},
			expectedMsg:   "Done!",
			expectedYear:  time.Now().Year(),
			expectedMonth: time.Now().Month(),
			expectedDay:   time.Now().Day(),
		},
		{
			name:          "message with hashtag and date",
			args:          []string{"Meeting", "#important", "2026/02/01"},
			expectedMsg:   "Meeting #important",
			expectedYear:  2026,
			expectedMonth: time.February,
			expectedDay:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, timestamp := parseAddArgs(tt.args)

			if message != tt.expectedMsg {
				t.Errorf("message = %q, want %q", message, tt.expectedMsg)
			}

			if timestamp.Year() != tt.expectedYear {
				t.Errorf("year = %d, want %d", timestamp.Year(), tt.expectedYear)
			}

			if timestamp.Month() != tt.expectedMonth {
				t.Errorf("month = %v, want %v", timestamp.Month(), tt.expectedMonth)
			}

			if timestamp.Day() != tt.expectedDay {
				t.Errorf("day = %d, want %d", timestamp.Day(), tt.expectedDay)
			}
		})
	}
}

func TestGetStorage(t *testing.T) {
	// Save original env vars
	originalBragDir := os.Getenv("BRAG_DIR")
	originalBragFile := os.Getenv("BRAG_FILE")
	defer func() {
		os.Setenv("BRAG_DIR", originalBragDir)
		os.Setenv("BRAG_FILE", originalBragFile)
	}()

	tests := []struct {
		name         string
		bragDirFlag  string
		winsFileFlag string
		bragDirEnv   string
		bragFileEnv  string
		expectedPath string // relative to home
	}{
		{
			name:         "defaults",
			bragDirFlag:  "",
			winsFileFlag: "",
			bragDirEnv:   "",
			bragFileEnv:  "",
			expectedPath: ".brag/wins.txt",
		},
		{
			name:         "custom dir via env",
			bragDirFlag:  "",
			winsFileFlag: "",
			bragDirEnv:   ".brag-test",
			bragFileEnv:  "",
			expectedPath: ".brag-test/wins.txt",
		},
		{
			name:         "custom file via env",
			bragDirFlag:  "",
			winsFileFlag: "",
			bragDirEnv:   "",
			bragFileEnv:  "2026.txt",
			expectedPath: ".brag/2026.txt",
		},
		{
			name:         "custom dir and file via env",
			bragDirFlag:  "",
			winsFileFlag: "",
			bragDirEnv:   ".brag-work",
			bragFileEnv:  "work-wins.txt",
			expectedPath: ".brag-work/work-wins.txt",
		},
		{
			name:         "flag overrides env for dir",
			bragDirFlag:  ".brag-override",
			winsFileFlag: "",
			bragDirEnv:   ".brag-env",
			bragFileEnv:  "",
			expectedPath: ".brag-override/wins.txt",
		},
		{
			name:         "flag overrides env for file",
			bragDirFlag:  "",
			winsFileFlag: "override.txt",
			bragDirEnv:   "",
			bragFileEnv:  "env.txt",
			expectedPath: ".brag/override.txt",
		},
		{
			name:         "flags override everything",
			bragDirFlag:  ".brag-flag",
			winsFileFlag: "flag.txt",
			bragDirEnv:   ".brag-env",
			bragFileEnv:  "env.txt",
			expectedPath: ".brag-flag/flag.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flags
			bragDir = tt.bragDirFlag
			winsFile = tt.winsFileFlag

			// Set env vars
			os.Setenv("BRAG_DIR", tt.bragDirEnv)
			os.Setenv("BRAG_FILE", tt.bragFileEnv)

			store, err := getStorage()
			if err != nil {
				t.Fatalf("getStorage() error = %v", err)
			}

			// Get home directory
			home, err := os.UserHomeDir()
			if err != nil {
				t.Fatalf("failed to get home dir: %v", err)
			}

			expectedPath := home + "/" + tt.expectedPath
			actualPath := store.GetFilePath()

			if actualPath != expectedPath {
				t.Errorf("path = %q, want %q", actualPath, expectedPath)
			}
		})
	}

	// Clean up
	bragDir = ""
	winsFile = ""
	os.Unsetenv("BRAG_DIR")
	os.Unsetenv("BRAG_FILE")
}

func TestAddWinIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Reset global flags
	bragDir = ""
	winsFile = ""

	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	tests := []struct {
		name        string
		message     string
		timestamp   time.Time
		expectError bool
	}{
		{
			name:        "simple message",
			message:     "Fixed a bug",
			timestamp:   time.Date(2026, 1, 8, 14, 30, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "message with tag",
			message:     "Completed feature #work",
			timestamp:   time.Date(2026, 1, 8, 15, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "message with multiple tags",
			message:     "Code review #work #backend #urgent",
			timestamp:   time.Date(2026, 1, 8, 16, 30, 0, 0, time.UTC),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.AddWin(tt.message, tt.timestamp)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				// Verify the win was added
				wins, err := store.ReadAllWins()
				if err != nil {
					t.Fatalf("Failed to read wins: %v", err)
				}

				found := false
				for _, win := range wins {
					if win.Message == tt.message {
						found = true
						if !win.Timestamp.Equal(tt.timestamp) {
							t.Errorf("timestamp = %v, want %v", win.Timestamp, tt.timestamp)
						}
						break
					}
				}

				if !found {
					t.Errorf("Win with message %q not found", tt.message)
				}
			}
		})
	}
}

func TestRootCmdWithArgs(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Reset global flags
	bragDir = ""
	winsFile = ""

	// Create storage
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Test that we can add a win via the root command
	testMessage := "Test win #testing"
	testTime := time.Now()

	err = store.AddWin(testMessage, testTime)
	if err != nil {
		t.Fatalf("Failed to add win: %v", err)
	}

	// Verify the win was added
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	if len(wins) == 0 {
		t.Fatal("No wins found after adding")
	}

	lastWin := wins[len(wins)-1]
	if lastWin.Message != testMessage {
		t.Errorf("message = %q, want %q", lastWin.Message, testMessage)
	}

	if !lastWin.HasTag("#testing") {
		t.Error("Win should have #testing tag")
	}
}

func TestParseAddArgsEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedMsg string
	}{
		{
			name:        "date-like string that's not a valid date",
			args:        []string{"Meeting", "on", "2026/99/99"},
			expectedMsg: "Meeting on 2026/99/99",
		},
		{
			name:        "partial date",
			args:        []string{"Task", "due", "01/08"},
			expectedMsg: "Task due 01/08",
		},
		{
			name:        "date with extra characters",
			args:        []string{"Done", "2026-01-08extra"},
			expectedMsg: "Done 2026-01-08extra",
		},
		{
			name:        "message ending with valid date format year",
			args:        []string{"Started", "in", "2026"},
			expectedMsg: "Started in 2026",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, _ := parseAddArgs(tt.args)

			if message != tt.expectedMsg {
				t.Errorf("message = %q, want %q", message, tt.expectedMsg)
			}
		})
	}
}
