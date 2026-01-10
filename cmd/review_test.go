package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
)

func TestParseDateRange(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		fromStr      string
		toStr        string
		checkRelative bool // For relative dates (using now)
		expectedFrom time.Time
		expectedTo   time.Time
	}{
		{
			name:          "no dates - defaults to last 7 days",
			fromStr:       "",
			toStr:         "",
			checkRelative: true,
			expectedFrom:  now.AddDate(0, 0, -7),
			expectedTo:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()),
		},
		{
			name:          "from date YYYY/MM/DD",
			fromStr:       "2026/01/01",
			toStr:         "",
			checkRelative: true,
			expectedFrom:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local),
			expectedTo:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()),
		},
		{
			name:          "from date YYYY-MM-DD",
			fromStr:       "2026-01-01",
			toStr:         "",
			checkRelative: true,
			expectedFrom:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local),
			expectedTo:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()),
		},
		{
			name:          "from date MM/DD/YYYY",
			fromStr:       "01/15/2026",
			toStr:         "",
			checkRelative: true,
			expectedFrom:  time.Date(2026, 1, 15, 0, 0, 0, 0, time.Local),
			expectedTo:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()),
		},
		{
			name:          "from date MM-DD-YYYY",
			fromStr:       "01-15-2026",
			toStr:         "",
			checkRelative: true,
			expectedFrom:  time.Date(2026, 1, 15, 0, 0, 0, 0, time.Local),
			expectedTo:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()),
		},
		{
			name:          "to date YYYY/MM/DD",
			fromStr:       "",
			toStr:         "2026/01/31",
			checkRelative: true,
			expectedFrom:  now.AddDate(0, 0, -7),
			expectedTo:    time.Date(2026, 1, 31, 23, 59, 59, 0, time.Local),
		},
		{
			name:          "both dates",
			fromStr:       "2026/01/01",
			toStr:         "2026/01/31",
			checkRelative: false,
			expectedFrom:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local),
			expectedTo:    time.Date(2026, 1, 31, 23, 59, 59, 0, time.Local),
		},
		{
			name:          "invalid from date - uses default",
			fromStr:       "invalid",
			toStr:         "",
			checkRelative: true,
			expectedFrom:  now.AddDate(0, 0, -7),
			expectedTo:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()),
		},
		{
			name:          "invalid to date - uses default",
			fromStr:       "",
			toStr:         "invalid",
			checkRelative: true,
			expectedFrom:  now.AddDate(0, 0, -7),
			expectedTo:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to := parseDateRange(tt.fromStr, tt.toStr)

			if tt.checkRelative {
				// For relative dates, check they're close enough
				if tt.fromStr == "" || tt.fromStr == "invalid" {
					// From date defaults to 7 days ago
					if from.Before(now.AddDate(0, 0, -7).Add(-1*time.Minute)) ||
						from.After(now.AddDate(0, 0, -7).Add(1*time.Minute)) {
						t.Errorf("from date too far from expected default (7 days ago)")
					}
				}

				if tt.toStr == "" || tt.toStr == "invalid" {
					// To date defaults to end of today
					// Check it's on the same day as now and is set to 23:59:59
					if to.Year() != now.Year() || to.Month() != now.Month() || to.Day() != now.Day() {
						t.Errorf("to date not on same day as now: got %v, expected day %v", to, now)
					}
					if to.Hour() != 23 || to.Minute() != 59 || to.Second() != 59 {
						t.Errorf("to time not set to end of day: %v", to)
					}
				}
			}

			// For explicit dates, check exact values (ignoring nanoseconds)
			if tt.fromStr != "" && tt.fromStr != "invalid" {
				if from.Year() != tt.expectedFrom.Year() ||
					from.Month() != tt.expectedFrom.Month() ||
					from.Day() != tt.expectedFrom.Day() {
					t.Errorf("from = %v, want %v", from, tt.expectedFrom)
				}
			}

			if tt.toStr != "" && tt.toStr != "invalid" {
				// For to date, check it's set to end of day
				if to.Hour() != 23 || to.Minute() != 59 || to.Second() != 59 {
					t.Errorf("to time not set to end of day: %v", to)
				}
				if to.Year() != tt.expectedTo.Year() ||
					to.Month() != tt.expectedTo.Month() ||
					to.Day() != tt.expectedTo.Day() {
					t.Errorf("to = %v, want %v", to, tt.expectedTo)
				}
			}
		})
	}
}

func TestFilterWins(t *testing.T) {
	testWins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
			Message:   "Win 1 #work",
			Tags:      []string{"#work"},
		},
		{
			Timestamp: time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC),
			Message:   "Win 2 #personal",
			Tags:      []string{"#personal"},
		},
		{
			Timestamp: time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC),
			Message:   "Win 3",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 15, 16, 0, 0, 0, time.UTC),
			Message:   "Win 4 #work #backend",
			Tags:      []string{"#work", "#backend"},
		},
		{
			Timestamp: time.Date(2026, 1, 20, 18, 0, 0, 0, time.UTC),
			Message:   "Win 5 #personal #health",
			Tags:      []string{"#personal", "#health"},
		},
	}

	tests := []struct {
		name          string
		from          time.Time
		to            time.Time
		tag           string
		taggedOnly    bool
		expectedCount int
		expectedMsgs  []string
	}{
		{
			name:          "all wins in range",
			from:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
			tag:           "",
			taggedOnly:    false,
			expectedCount: 5,
			expectedMsgs:  []string{"Win 1 #work", "Win 2 #personal", "Win 3", "Win 4 #work #backend", "Win 5 #personal #health"},
		},
		{
			name:          "date range filter",
			from:          time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 15, 23, 59, 59, 0, time.UTC),
			tag:           "",
			taggedOnly:    false,
			expectedCount: 3,
			expectedMsgs:  []string{"Win 2 #personal", "Win 3", "Win 4 #work #backend"},
		},
		{
			name:          "filter by specific tag",
			from:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
			tag:           "#work",
			taggedOnly:    false,
			expectedCount: 2,
			expectedMsgs:  []string{"Win 1 #work", "Win 4 #work #backend"},
		},
		{
			name:          "filter by tag without hash",
			from:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
			tag:           "personal",
			taggedOnly:    false,
			expectedCount: 2,
			expectedMsgs:  []string{"Win 2 #personal", "Win 5 #personal #health"},
		},
		{
			name:          "tagged only",
			from:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
			tag:           "",
			taggedOnly:    true,
			expectedCount: 4,
			expectedMsgs:  []string{"Win 1 #work", "Win 2 #personal", "Win 4 #work #backend", "Win 5 #personal #health"},
		},
		{
			name:          "tagged only with specific tag",
			from:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
			tag:           "backend",
			taggedOnly:    true,
			expectedCount: 1,
			expectedMsgs:  []string{"Win 4 #work #backend"},
		},
		{
			name:          "narrow date range",
			from:          time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 10, 23, 59, 59, 0, time.UTC),
			tag:           "",
			taggedOnly:    false,
			expectedCount: 1,
			expectedMsgs:  []string{"Win 3"},
		},
		{
			name:          "no matches - wrong tag",
			from:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
			tag:           "nonexistent",
			taggedOnly:    false,
			expectedCount: 0,
			expectedMsgs:  []string{},
		},
		{
			name:          "no matches - date range before wins",
			from:          time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			tag:           "",
			taggedOnly:    false,
			expectedCount: 0,
			expectedMsgs:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterWins(testWins, tt.from, tt.to, tt.tag, tt.taggedOnly)

			if len(filtered) != tt.expectedCount {
				t.Errorf("filtered count = %d, want %d", len(filtered), tt.expectedCount)
			}

			for i, expectedMsg := range tt.expectedMsgs {
				if i >= len(filtered) {
					t.Errorf("missing expected win: %q", expectedMsg)
					continue
				}
				if filtered[i].Message != expectedMsg {
					t.Errorf("win[%d] message = %q, want %q", i, filtered[i].Message, expectedMsg)
				}
			}
		})
	}
}

func TestReviewWinsIntegration(t *testing.T) {
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

	// Add test wins
	testWins := []*models.Win{
		{
			Timestamp: time.Now().AddDate(0, 0, -2),
			Message:   "Recent win #work",
			Tags:      []string{"#work"},
		},
		{
			Timestamp: time.Now().AddDate(0, 0, -5),
			Message:   "Last week win #personal",
			Tags:      []string{"#personal"},
		},
		{
			Timestamp: time.Now().AddDate(0, 0, -10),
			Message:   "Old win",
			Tags:      []string{},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	// Read and filter wins (default: last 7 days)
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	from, to := parseDateRange("", "")
	filtered := filterWins(wins, from, to, "", false)

	// Should get the two recent wins (within last 7 days)
	if len(filtered) != 2 {
		t.Errorf("filtered count = %d, want 2 (wins from last 7 days)", len(filtered))
	}

	// Filter by tag
	tagFiltered := filterWins(wins, from, to, "work", false)
	if len(tagFiltered) != 1 {
		t.Errorf("tag filtered count = %d, want 1", len(tagFiltered))
	}
	if len(tagFiltered) > 0 && tagFiltered[0].Message != "Recent win #work" {
		t.Errorf("wrong win filtered: %q", tagFiltered[0].Message)
	}
}

func TestFilterWinsEdgeCases(t *testing.T) {
	testWins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			Message:   "Exactly at midnight",
			Tags:      []string{},
		},
		{
			Timestamp: time.Date(2026, 1, 10, 23, 59, 59, 0, time.UTC),
			Message:   "End of day",
			Tags:      []string{},
		},
	}

	tests := []struct {
		name          string
		from          time.Time
		to            time.Time
		expectedCount int
	}{
		{
			name:          "inclusive start boundary",
			from:          time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 10, 23, 59, 59, 0, time.UTC),
			expectedCount: 2,
		},
		{
			name:          "exclusive end boundary",
			from:          time.Date(2026, 1, 10, 0, 0, 1, 0, time.UTC),
			to:            time.Date(2026, 1, 10, 23, 59, 58, 0, time.UTC),
			expectedCount: 0,
		},
		{
			name:          "single second range",
			from:          time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterWins(testWins, tt.from, tt.to, "", false)

			if len(filtered) != tt.expectedCount {
				t.Errorf("filtered count = %d, want %d", len(filtered), tt.expectedCount)
			}
		})
	}
}
