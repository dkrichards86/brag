# Testing Guide

This document describes the testing and CI setup for the brag project.

## Quick Start

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Format code
make fmt

# Run full CI pipeline locally
make ci
```

## Test Coverage

Current test coverage:
- `internal/models`: 100% coverage
- `internal/storage`: 82.8% coverage
- Overall: High coverage of core business logic

## Running Tests

### Basic Test Run

```bash
go test ./...
```

### With Race Detection

```bash
go test -race ./...
```

### With Coverage Report

```bash
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out -o coverage.html
```

Or simply:

```bash
make test-coverage
```

This will generate an HTML coverage report at `coverage.html`.

## Linting

The project uses [golangci-lint](https://golangci-lint.run/) for comprehensive code analysis.

### Installing golangci-lint

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Running the Linter

```bash
golangci-lint run --timeout=5m
```

Or:

```bash
make lint
```

### Enabled Linters

The following linters are enabled (see [.golangci.yml](.golangci.yml)):

- **errcheck**: Check for unchecked errors
- **gosimple**: Simplify code
- **govet**: Vet examines Go source code
- **ineffassign**: Detect ineffectual assignments
- **staticcheck**: Advanced Go linter
- **unused**: Check for unused code
- **gofmt**: Check formatting
- **misspell**: Check for misspellings
- **goconst**: Find repeated strings that could be constants
- **gocyclo**: Check cyclomatic complexity
- **gosec**: Security checks
- **revive**: Fast, configurable, extensible linter
- **copyloopvar**: Check for loop variable issues (Go 1.22+)
- **unconvert**: Remove unnecessary type conversions

## Continuous Integration

The project uses GitHub Actions for CI. The workflow is defined in [.github/workflows/ci.yml](.github/workflows/ci.yml).

### CI Pipeline

The CI pipeline runs on every push to `main`/`master` and on all pull requests:

1. **Test Job**
   - Tests against Go 1.21, 1.22, and 1.23
   - Runs tests with race detection
   - Generates coverage reports
   - Uploads coverage to Codecov (on Go 1.23 only)

2. **Lint Job**
   - Runs golangci-lint with all configured linters
   - Uses the latest Go version (1.23)

3. **Build Job**
   - Cross-compiles for multiple platforms:
     - Linux (amd64, arm64)
     - macOS (amd64, arm64)
     - Windows (amd64)
   - Uploads build artifacts

### Viewing CI Results

1. Go to the [Actions tab](../../actions) in GitHub
2. Select the workflow run you want to view
3. Review the test results, lint output, and build artifacts

## Writing Tests

### Test Structure

Tests are organized alongside the code they test:

```
internal/
├── models/
│   ├── win.go
│   └── win_test.go
└── storage/
    ├── storage.go
    └── storage_test.go
```

### Test Conventions

1. **File Naming**: Test files are named `*_test.go`
2. **Function Naming**: Test functions start with `Test` (e.g., `TestParseWin`)
3. **Table-Driven Tests**: Use table-driven tests for multiple test cases
4. **Helper Functions**: Use `t.Helper()` in test helper functions
5. **Test Cleanup**: Use `t.Cleanup()` or `defer` for cleanup

### Example Test

```go
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
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseWin(tt.line)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseWin() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Assertions...
        })
    }
}
```

## Makefile Targets

The [Makefile](Makefile) provides convenient shortcuts:

- `make test` - Run tests
- `make test-verbose` - Run tests with verbose output
- `make test-coverage` - Run tests with coverage report
- `make lint` - Run linter
- `make fmt` - Format code
- `make vet` - Run go vet
- `make build` - Build the binary
- `make build-all` - Build for all platforms
- `make clean` - Clean build artifacts
- `make install` - Install to GOPATH/bin
- `make ci` - Run full CI pipeline locally
- `make help` - Show all available targets

## Pre-commit Checklist

Before committing code, ensure:

1. ✅ All tests pass: `make test`
2. ✅ Linter passes: `make lint`
3. ✅ Code is formatted: `make fmt`
4. ✅ New code has tests
5. ✅ Coverage remains high

Or simply run:

```bash
make ci
```

This runs the full CI pipeline locally.

## Troubleshooting

### golangci-lint Not Found

Install it:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Make sure `$GOPATH/bin` is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Tests Failing with Race Detector

Race conditions are caught by the `-race` flag. Fix them before committing:

1. Review the race detector output
2. Add proper synchronization (mutexes, channels, etc.)
3. Re-run tests until they pass

### Coverage Dropping

If coverage drops below acceptable levels:

1. Identify untested code paths
2. Add tests for new functionality
3. Consider edge cases and error paths
4. Aim for at least 80% coverage

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Table Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
