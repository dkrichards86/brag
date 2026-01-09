# CLAUDE.md - AI Assistant Context

This file provides context for AI assistants (like Claude) working on the brag project.

## Project Overview

**Brag** is a CLI tool for tracking daily micro-wins and achievements. It's designed to be simple, fast, and frictionless - allowing users to quickly log accomplishments throughout their day.

## Architecture

### Technology Stack
- **Language**: Go 1.x
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Storage**: Plain text file (`~/.brag/wins.txt`)
- **No external dependencies** for runtime (beyond Go stdlib and Cobra)

### Project Structure

```
brag/
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI pipeline
├── cmd/                        # Command implementations
│   ├── root.go                # Root command & default add functionality
│   ├── review.go              # Review/filter wins
│   ├── list.go                # List wins with line numbers
│   ├── edit.go                # Edit wins (file or specific entry)
│   ├── delete.go              # Delete specific win
│   └── search.go              # Search through wins
├── internal/
│   ├── models/
│   │   ├── win.go             # Win data model & parsing
│   │   └── win_test.go        # Unit tests (100% coverage)
│   └── storage/
│       ├── storage.go         # File I/O operations
│       └── storage_test.go    # Unit tests (82.8% coverage)
├── main.go                    # Entry point
├── Makefile                   # Build and test automation
├── .golangci.yml              # Linter configuration
├── README.md                  # User documentation
├── TESTING.md                 # Testing and CI documentation
├── CLAUDE.md                  # This file
├── LICENSE                    # MIT License
└── .gitignore
```

## Design Principles

1. **Simplicity First**: The tool should be dead simple to use. Default command is add.
2. **Plain Text Storage**: Human-readable format that can be manually edited
3. **Fast**: No database overhead, minimal dependencies
4. **Flexible Tagging**: Use hashtags for categorization, no rigid structure
5. **Shell-Friendly**: Works with standard shell conventions

## Key Design Decisions

### Default Command Pattern
- Running `brag <message>` (no subcommand) adds a new win
- This is the most common operation, so it should be the easiest
- Implemented using Cobra's `RunE` on the root command

### File Format
```
YYYY-MM-DD HH:MM | message with #tags
```

**Why this format?**
- Human-readable and easily parseable
- Timestamp first for chronological sorting
- Pipe separator is visually clear
- Hashtags are natural and familiar

### Macro Wins Philosophy
- No hardcoded `#macro` tag
- **Any tagged win** is considered a "macro win"
- Reasoning: If you bothered to tag it, it's worth highlighting
- Use `--tagged` flag to filter for wins with any tags
- Use `--tag <name>` to filter for specific tags

### Shell Quoting Requirement
Users must quote messages with hashtags:
```bash
brag "Fixed bug #work #backend"  # Correct
brag Fixed bug #work #backend     # Wrong - shell treats # as comment
```

This is a shell limitation, not a tool limitation. Documented clearly in README.

## Data Models

### Win Structure
```go
type Win struct {
    Timestamp time.Time
    Message   string
    Tags      []string  // Extracted from message
}
```

- Tags are extracted via regex: `#\w+`
- Tags are stored inline in the message (not separately)
- This keeps the file format simple and portable

## Common Operations

### Adding a Win
1. Parse arguments
2. Check if last arg is a date
3. Join remaining args as message
4. Append to file with timestamp

### Filtering Wins
1. Read all wins from file
2. Apply date range filter (default: last 7 days)
3. Apply tag filters if specified
4. Display results

### Editing/Deleting
1. Read all wins into memory
2. Modify the slice
3. Write entire slice back to file
4. Simple but works fine for typical usage (hundreds of wins)

## Testing Infrastructure

### Test Coverage
The project has comprehensive unit tests:
- **internal/models/win_test.go**: 100% coverage
  - Tests for ParseWin, Format, HasTag, extractTags
  - Edge cases, error handling, round-trip testing
- **internal/storage/storage_test.go**: 82.8% coverage
  - Tests for all CRUD operations
  - File I/O edge cases, invalid input handling

### Running Tests
```bash
# Quick test run
make test

# With coverage report
make test-coverage

# Run linter
make lint

# Full CI pipeline locally
make ci
```

See [TESTING.md](TESTING.md) for comprehensive testing documentation.

### Continuous Integration
GitHub Actions runs on every push and PR:
- Tests on Go 1.21, 1.22, and 1.23
- Linting with golangci-lint (15 enabled linters)
- Cross-platform builds (Linux, macOS, Windows)
- Coverage reporting to Codecov

### Testing Considerations

When testing or extending this tool:

1. **Date Parsing**: Support multiple formats (YYYY/MM/DD, YYYY-MM-DD, MM/DD/YYYY, MM-DD-YYYY)
2. **Tag Extraction**: Tags must match `#\w+` pattern
3. **File Safety**: Always read before write operations
4. **Empty Lines**: Skip blank lines when reading
5. **Invalid Lines**: Warn but continue if a line can't be parsed
6. **Write Tests First**: All new functionality should have tests before implementation
7. **Table-Driven Tests**: Use table-driven patterns for multiple test cases
8. **Test Isolation**: Use t.TempDir() for file-based tests to avoid conflicts

## Extension Ideas

Future enhancements that would fit the design:

1. **Export Command**: Export filtered wins to markdown/JSON
2. **Stats Command**: Show tag frequency, wins per day/week
3. **Sync**: Sync wins file across machines (git-based?)
4. **Interactive Mode**: TUI for browsing/editing wins
5. **Shell Completions**: Better autocompletion for tags
6. **Undo**: Undo last add operation

## Important Notes for AI Assistants

### When Making Changes

1. **Preserve Simplicity**: Don't add complexity without clear user value
2. **Keep File Format Stable**: Changes to format require migration path
3. **Test Date Parsing**: Multiple formats make this error-prone
4. **Validate User Input**: But fail gracefully with helpful messages
5. **Maintain Backward Compatibility**: Users have existing wins files
6. **Write Tests First**: Add tests before implementing new functionality
7. **Run CI Locally**: Use `make ci` before pushing changes
8. **Check Linter**: All code must pass `make lint` with zero issues

### Code Style

- Use Go standard formatting (`gofmt` and `make fmt`)
- Prefer standard library over external dependencies
- Clear error messages to stderr
- Success messages to stdout
- Exit codes: 0 for success, 1 for errors
- All code must pass golangci-lint checks
- Maintain test coverage above 80% for new code

### Common Gotchas

1. **Cobra Command Registration**: Subcommands must be registered in `init()` functions
2. **File Paths**: Use `filepath.Join()` for cross-platform compatibility
3. **Date Parsing**: Time.Parse() format string is Go's weird reference date: `2006-01-02 15:04`
4. **Tag Extraction**: Remember to handle tags with or without leading `#`
5. **Shell Escaping**: Users must quote hashtags - document this clearly

## Building and Installing

### Using Makefile (Recommended)
```bash
# Build locally
make build

# Install to GOPATH/bin
make install

# Build for all platforms
make build-all

# Run tests and build
make ci

# See all available commands
make help
```

### Using Go directly
```bash
# Build locally
go build -o brag

# Install to GOPATH/bin
go install

# Cross-compile for other platforms
GOOS=linux GOARCH=amd64 go build -o brag-linux-amd64
GOOS=darwin GOARCH=arm64 go build -o brag-darwin-arm64
```

### Makefile Targets
The Makefile provides these common commands:
- `make test` - Run tests
- `make test-coverage` - Generate coverage report
- `make lint` - Run golangci-lint
- `make fmt` - Format all code
- `make build` - Build binary
- `make build-all` - Cross-compile for all platforms
- `make clean` - Remove build artifacts
- `make ci` - Run full CI pipeline locally

## Questions for Users

When users request changes, consider asking:

1. **File Format Changes**: "This would change the file format. Do you want to add migration logic for existing files?"
2. **New Dependencies**: "This would add a new dependency. Is the trade-off worth it for this feature?"
3. **Breaking Changes**: "This would break existing usage patterns. Should we version this or add a new command instead?"

## Related Commands

Users coming from other tools might search for:
- `git commit` (similar append-only log pattern)
- `todo` tools (but brag is for done items, not todos)
- `done` command (similar concept)
- `accomplishments` tracking

## Success Metrics

The tool is succeeding if:
1. Users run `brag` multiple times per day
2. Users actually review their wins periodically
3. Users find it helpful for performance reviews
4. The tool "gets out of the way" - minimal friction

## Philosophy

> "Track your wins without thinking about it. Review them when you need to feel accomplished or prepare for performance evaluations."

The goal is to make logging wins as easy as possible, with just enough structure (tags) to be useful later, but not so much that it becomes a chore.
