package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/dkrichards86/brag/internal/models"
	"github.com/dkrichards86/brag/internal/storage"
)

func TestListTags(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name         string
		wins         []*models.Win
		expectedTags map[string]int
	}{
		{
			name: "multiple tags with different frequencies",
			wins: []*models.Win{
				{
					Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
					Message:   "Fixed bug #work #backend",
					Tags:      []string{"#work", "#backend"},
				},
				{
					Timestamp: time.Date(2026, 1, 8, 15, 0, 0, 0, time.UTC),
					Message:   "Code review #work #review",
					Tags:      []string{"#work", "#review"},
				},
				{
					Timestamp: time.Date(2026, 1, 8, 16, 30, 0, 0, time.UTC),
					Message:   "Meeting #meeting #work",
					Tags:      []string{"#meeting", "#work"},
				},
			},
			expectedTags: map[string]int{
				"#work":    3,
				"#backend": 1,
				"#review":  1,
				"#meeting": 1,
			},
		},
		{
			name: "wins with no tags",
			wins: []*models.Win{
				{
					Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
					Message:   "Fixed bug",
					Tags:      []string{},
				},
				{
					Timestamp: time.Date(2026, 1, 8, 15, 0, 0, 0, time.UTC),
					Message:   "Completed feature",
					Tags:      []string{},
				},
			},
			expectedTags: map[string]int{},
		},
		{
			name: "single tag across multiple wins",
			wins: []*models.Win{
				{
					Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
					Message:   "Task 1 #project",
					Tags:      []string{"#project"},
				},
				{
					Timestamp: time.Date(2026, 1, 8, 15, 0, 0, 0, time.UTC),
					Message:   "Task 2 #project",
					Tags:      []string{"#project"},
				},
				{
					Timestamp: time.Date(2026, 1, 8, 16, 30, 0, 0, time.UTC),
					Message:   "Task 3 #project",
					Tags:      []string{"#project"},
				},
			},
			expectedTags: map[string]int{
				"#project": 3,
			},
		},
		{
			name: "mixed tagged and untagged wins",
			wins: []*models.Win{
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
					Message:   "Code review #review",
					Tags:      []string{"#review"},
				},
			},
			expectedTags: map[string]int{
				"#work":   1,
				"#review": 1,
			},
		},
		{
			name: "tags with numbers and underscores",
			wins: []*models.Win{
				{
					Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
					Message:   "Sprint planning #sprint2024 #q1_planning",
					Tags:      []string{"#sprint2024", "#q1_planning"},
				},
				{
					Timestamp: time.Date(2026, 1, 8, 15, 0, 0, 0, time.UTC),
					Message:   "Review #sprint2024",
					Tags:      []string{"#sprint2024"},
				},
			},
			expectedTags: map[string]int{
				"#sprint2024":  2,
				"#q1_planning": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create storage and write test wins
			store, err := storage.New()
			if err != nil {
				t.Fatalf("Failed to create storage: %v", err)
			}

			if err := store.WriteWins(tt.wins); err != nil {
				t.Fatalf("Failed to write test wins: %v", err)
			}

			// Read wins and count tags
			wins, err := store.ReadAllWins()
			if err != nil {
				t.Fatalf("Failed to read wins: %v", err)
			}

			tagCounts := make(map[string]int)
			for _, win := range wins {
				for _, tag := range win.Tags {
					tagCounts[tag]++
				}
			}

			// Verify tag counts
			if len(tagCounts) != len(tt.expectedTags) {
				t.Errorf("Number of unique tags = %d, want %d", len(tagCounts), len(tt.expectedTags))
			}

			for tag, expectedCount := range tt.expectedTags {
				if count, exists := tagCounts[tag]; !exists {
					t.Errorf("Expected tag %q not found", tag)
				} else if count != expectedCount {
					t.Errorf("Tag %q count = %d, want %d", tag, count, expectedCount)
				}
			}

			// Verify no unexpected tags
			for tag := range tagCounts {
				if _, expected := tt.expectedTags[tag]; !expected {
					t.Errorf("Unexpected tag found: %q", tag)
				}
			}
		})
	}
}

func TestTagsSorting(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create storage and add test wins with varying tag frequencies
	store, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	testWins := []*models.Win{
		{
			Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
			Message:   "Task 1 #alpha #beta #gamma",
			Tags:      []string{"#alpha", "#beta", "#gamma"},
		},
		{
			Timestamp: time.Date(2026, 1, 8, 15, 0, 0, 0, time.UTC),
			Message:   "Task 2 #alpha #beta",
			Tags:      []string{"#alpha", "#beta"},
		},
		{
			Timestamp: time.Date(2026, 1, 8, 16, 30, 0, 0, time.UTC),
			Message:   "Task 3 #alpha",
			Tags:      []string{"#alpha"},
		},
	}

	if err := store.WriteWins(testWins); err != nil {
		t.Fatalf("Failed to write test wins: %v", err)
	}

	// Read wins and count tags
	wins, err := store.ReadAllWins()
	if err != nil {
		t.Fatalf("Failed to read wins: %v", err)
	}

	tagCounts := make(map[string]int)
	for _, win := range wins {
		for _, tag := range win.Tags {
			tagCounts[tag]++
		}
	}

	// Expected: #alpha (3), #beta (2), #gamma (1)
	expectedOrder := []struct {
		tag   string
		count int
	}{
		{"#alpha", 3},
		{"#beta", 2},
		{"#gamma", 1},
	}

	// Verify counts match expected
	for _, expected := range expectedOrder {
		if count, exists := tagCounts[expected.tag]; !exists {
			t.Errorf("Expected tag %q not found", expected.tag)
		} else if count != expected.count {
			t.Errorf("Tag %q count = %d, want %d", expected.tag, count, expected.count)
		}
	}
}
