package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Win represents a single micro-win entry
type Win struct {
	Timestamp time.Time
	Message   string
	Tags      []string
}

// ParseWin parses a line from the wins file into a Win struct
// Expected format: "2026-01-08 14:32 | message with #tags"
func ParseWin(line string) (*Win, error) {
	parts := strings.SplitN(line, "|", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format: missing separator '|'")
	}

	// Parse timestamp
	timestampStr := strings.TrimSpace(parts[0])
	timestamp, err := time.Parse("2006-01-02 15:04", timestampStr)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	// Parse message and extract tags
	message := strings.TrimSpace(parts[1])
	tags := extractTags(message)

	return &Win{
		Timestamp: timestamp,
		Message:   message,
		Tags:      tags,
	}, nil
}

// Format returns the formatted string representation of a Win
func (w *Win) Format() string {
	return fmt.Sprintf("%s | %s", w.Timestamp.Format("2006-01-02 15:04"), w.Message)
}

// HasTag checks if the win has a specific tag
func (w *Win) HasTag(tag string) bool {
	tag = strings.ToLower(tag)
	if !strings.HasPrefix(tag, "#") {
		tag = "#" + tag
	}
	for _, t := range w.Tags {
		if strings.ToLower(t) == tag {
			return true
		}
	}
	return false
}

// extractTags extracts all hashtags from a message
func extractTags(message string) []string {
	re := regexp.MustCompile(`#\w+`)
	return re.FindAllString(message, -1)
}
