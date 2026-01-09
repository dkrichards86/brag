package models

import (
	"testing"
	"time"
)

func TestParseWin(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    *Win
		wantErr bool
	}{
		{
			name: "valid win with tags",
			line: "2026-01-08 14:32 | Fixed bug #work #backend",
			want: &Win{
				Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
				Message:   "Fixed bug #work #backend",
				Tags:      []string{"#work", "#backend"},
			},
			wantErr: false,
		},
		{
			name: "valid win without tags",
			line: "2026-01-08 10:00 | Completed project milestone",
			want: &Win{
				Timestamp: time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC),
				Message:   "Completed project milestone",
				Tags:      []string{},
			},
			wantErr: false,
		},
		{
			name: "valid win with single tag",
			line: "2026-01-08 09:15 | Morning standup #meeting",
			want: &Win{
				Timestamp: time.Date(2026, 1, 8, 9, 15, 0, 0, time.UTC),
				Message:   "Morning standup #meeting",
				Tags:      []string{"#meeting"},
			},
			wantErr: false,
		},
		{
			name:    "invalid format - missing separator",
			line:    "2026-01-08 14:32 No separator here",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid format - bad timestamp",
			line:    "invalid-date | Some message",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid format - empty line",
			line:    "",
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid win with extra whitespace",
			line: "  2026-01-08 14:32  |  Message with spaces  ",
			want: &Win{
				Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
				Message:   "Message with spaces",
				Tags:      []string{},
			},
			wantErr: false,
		},
		{
			name: "valid win with multiple tags in middle",
			line: "2026-01-08 11:00 | Reviewed #code and #documentation today",
			want: &Win{
				Timestamp: time.Date(2026, 1, 8, 11, 0, 0, 0, time.UTC),
				Message:   "Reviewed #code and #documentation today",
				Tags:      []string{"#code", "#documentation"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseWin(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if !got.Timestamp.Equal(tt.want.Timestamp) {
				t.Errorf("ParseWin() Timestamp = %v, want %v", got.Timestamp, tt.want.Timestamp)
			}
			if got.Message != tt.want.Message {
				t.Errorf("ParseWin() Message = %v, want %v", got.Message, tt.want.Message)
			}
			if len(got.Tags) != len(tt.want.Tags) {
				t.Errorf("ParseWin() Tags length = %d, want %d", len(got.Tags), len(tt.want.Tags))
				return
			}
			for i := range got.Tags {
				if got.Tags[i] != tt.want.Tags[i] {
					t.Errorf("ParseWin() Tags[%d] = %v, want %v", i, got.Tags[i], tt.want.Tags[i])
				}
			}
		})
	}
}

func TestWin_Format(t *testing.T) {
	tests := []struct {
		name string
		win  *Win
		want string
	}{
		{
			name: "win with tags",
			win: &Win{
				Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
				Message:   "Fixed bug #work #backend",
				Tags:      []string{"#work", "#backend"},
			},
			want: "2026-01-08 14:32 | Fixed bug #work #backend",
		},
		{
			name: "win without tags",
			win: &Win{
				Timestamp: time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC),
				Message:   "Completed milestone",
				Tags:      []string{},
			},
			want: "2026-01-08 10:00 | Completed milestone",
		},
		{
			name: "win with single digit hour",
			win: &Win{
				Timestamp: time.Date(2026, 1, 8, 9, 5, 0, 0, time.UTC),
				Message:   "Early win",
				Tags:      []string{},
			},
			want: "2026-01-08 09:05 | Early win",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.win.Format(); got != tt.want {
				t.Errorf("Win.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWin_HasTag(t *testing.T) {
	win := &Win{
		Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
		Message:   "Fixed bug #work #backend #API",
		Tags:      []string{"#work", "#backend", "#API"},
	}

	tests := []struct {
		name string
		tag  string
		want bool
	}{
		{
			name: "has tag with hash",
			tag:  "#work",
			want: true,
		},
		{
			name: "has tag without hash",
			tag:  "backend",
			want: true,
		},
		{
			name: "case insensitive match",
			tag:  "WORK",
			want: true,
		},
		{
			name: "case insensitive match with hash",
			tag:  "#Backend",
			want: true,
		},
		{
			name: "tag not present",
			tag:  "frontend",
			want: false,
		},
		{
			name: "mixed case tag present",
			tag:  "api",
			want: true,
		},
		{
			name: "partial match should fail",
			tag:  "wor",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := win.HasTag(tt.tag); got != tt.want {
				t.Errorf("Win.HasTag(%v) = %v, want %v", tt.tag, got, tt.want)
			}
		})
	}
}

func TestExtractTags(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    []string
	}{
		{
			name:    "multiple tags",
			message: "Fixed bug #work #backend #urgent",
			want:    []string{"#work", "#backend", "#urgent"},
		},
		{
			name:    "no tags",
			message: "Just a regular message",
			want:    []string{},
		},
		{
			name:    "single tag",
			message: "Meeting notes #meeting",
			want:    []string{"#meeting"},
		},
		{
			name:    "tags in middle",
			message: "Reviewed #code and #tests today",
			want:    []string{"#code", "#tests"},
		},
		{
			name:    "tag with numbers",
			message: "Sprint #sprint2024 planning",
			want:    []string{"#sprint2024"},
		},
		{
			name:    "tag with underscore",
			message: "Fixed #bug_fix today",
			want:    []string{"#bug_fix"},
		},
		{
			name:    "hash without word chars not a tag",
			message: "Cost is #$100",
			want:    []string{},
		},
		{
			name:    "multiple hashes together",
			message: "Tags: #tag1 #tag2 #tag3",
			want:    []string{"#tag1", "#tag2", "#tag3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTags(tt.message)
			if len(got) != len(tt.want) {
				t.Errorf("extractTags() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("extractTags()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestParseWinRoundTrip tests that Format() output can be parsed back
func TestParseWinRoundTrip(t *testing.T) {
	original := &Win{
		Timestamp: time.Date(2026, 1, 8, 14, 32, 0, 0, time.UTC),
		Message:   "Fixed bug #work #backend",
		Tags:      []string{"#work", "#backend"},
	}

	formatted := original.Format()
	parsed, err := ParseWin(formatted)
	if err != nil {
		t.Fatalf("ParseWin() failed on formatted output: %v", err)
	}

	if !parsed.Timestamp.Equal(original.Timestamp) {
		t.Errorf("Round trip Timestamp = %v, want %v", parsed.Timestamp, original.Timestamp)
	}
	if parsed.Message != original.Message {
		t.Errorf("Round trip Message = %v, want %v", parsed.Message, original.Message)
	}
	if len(parsed.Tags) != len(original.Tags) {
		t.Errorf("Round trip Tags length = %d, want %d", len(parsed.Tags), len(original.Tags))
	}
}
