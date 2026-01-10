package storage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dkrichards86/brag/internal/models"
)

const (
	defaultBragDir  = ".brag"
	defaultWinsFile = "wins.txt"

	// Environment variables for configuration
	envBragDir  = "BRAG_DIR"
	envWinsFile = "BRAG_FILE"
)

// Storage handles file operations for wins
type Storage struct {
	filePath string
}

// New creates a new Storage instance using default or environment-configured paths
// Respects BRAG_DIR and BRAG_FILE environment variables for testing and customization
func New() (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Allow environment variable override for directory name
	dirName := defaultBragDir
	if envDir := os.Getenv(envBragDir); envDir != "" {
		dirName = envDir
	}

	// Allow environment variable override for file name
	fileName := defaultWinsFile
	if envFile := os.Getenv(envWinsFile); envFile != "" {
		fileName = envFile
	}

	bragPath := filepath.Join(home, dirName)
	filePath := filepath.Join(bragPath, fileName)

	return newStorage(filePath)
}

// NewWithPath creates a new Storage instance with an explicit file path
// This is useful for testing or when you want complete control over the storage location
func NewWithPath(filePath string) (*Storage, error) {
	return newStorage(filePath)
}

// newStorage is the internal constructor that handles directory and file creation
func newStorage(filePath string) (*Storage, error) {
	bragPath := filepath.Dir(filePath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(bragPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create brag directory: %w", err)
	}

	// Create file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return nil, fmt.Errorf("failed to create wins file: %w", err)
		}
	}

	return &Storage{filePath: filePath}, nil
}

// GetFilePath returns the path to the wins file
func (s *Storage) GetFilePath() string {
	return s.filePath
}

// AddWin appends a new win to the file
func (s *Storage) AddWin(message string, timestamp time.Time) error {
	win := &models.Win{
		Timestamp: timestamp,
		Message:   message,
	}

	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open wins file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(win.Format() + "\n")
	if err != nil {
		return fmt.Errorf("failed to write win: %w", err)
	}

	return nil
}

// ReadAllWins reads all wins from the file
func (s *Storage) ReadAllWins() ([]*models.Win, error) {
	f, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open wins file: %w", err)
	}
	defer f.Close()

	var wins []*models.Win
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		win, err := models.ParseWin(line)
		if err != nil {
			// Skip invalid lines with a warning
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid line %d: %v\n", lineNum, err)
			continue
		}

		wins = append(wins, win)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading wins file: %w", err)
	}

	return wins, nil
}

// WriteWins writes all wins back to the file (for edit/delete operations)
func (s *Storage) WriteWins(wins []*models.Win) error {
	f, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open wins file: %w", err)
	}
	defer f.Close()

	for _, win := range wins {
		_, err = f.WriteString(win.Format() + "\n")
		if err != nil {
			return fmt.Errorf("failed to write win: %w", err)
		}
	}

	return nil
}

// DeleteWin deletes a win at the specified index (0-based)
func (s *Storage) DeleteWin(index int) error {
	wins, err := s.ReadAllWins()
	if err != nil {
		return err
	}

	if index < 0 || index >= len(wins) {
		return fmt.Errorf("invalid index: %d", index)
	}

	// Remove the win at the specified index
	wins = append(wins[:index], wins[index+1:]...)

	return s.WriteWins(wins)
}

// UpdateWin updates a win at the specified index (0-based)
func (s *Storage) UpdateWin(index int, message string) error {
	wins, err := s.ReadAllWins()
	if err != nil {
		return err
	}

	if index < 0 || index >= len(wins) {
		return fmt.Errorf("invalid index: %d", index)
	}

	// Update the message, keep the original timestamp
	wins[index].Message = message

	return s.WriteWins(wins)
}
