# Brag - Micro-Wins Tracker

A simple CLI tool to track your daily micro-wins and achievements. `brag` makes it easy to record your accomplishments, review them over time, and identify macro wins for year-end reports.

## What are Micro-Wins?

Micro-wins are the small daily accomplishments that often go unnoticed but add up to significant progress over time. They can be anything from fixing a bug, helping a colleague, learning something new, or shipping a feature. The concept stems from [Teresa Amabile's research on the progress principle](https://hbr.org/2011/05/the-power-of-small-wins) - tracking small wins boosts motivation and helps you see patterns in your work.

By logging micro-wins as they happen, you'll have a rich history to draw from during performance reviews, when updating your resume, or simply when you need a reminder of what you've accomplished.

## Features

- Quick win logging with hashtag support
- Post-dating wins with flexible date formats
- Review wins by date range
- Filter by tags (including macro wins)
- Edit and delete specific entries
- Search through your wins
- Simple plain text storage

## Installation

### Pre-built Binaries (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/dkrichards86/brag/releases).

**Linux (AMD64):**
```bash
curl -LO https://github.com/dkrichards86/brag/releases/latest/download/brag-linux-amd64
chmod +x brag-linux-amd64
sudo mv brag-linux-amd64 /usr/local/bin/brag
```

**macOS (Apple Silicon):**
```bash
curl -LO https://github.com/dkrichards86/brag/releases/latest/download/brag-darwin-arm64
chmod +x brag-darwin-arm64
sudo mv brag-darwin-arm64 /usr/local/bin/brag
```

**macOS (Intel):**
```bash
curl -LO https://github.com/dkrichards86/brag/releases/latest/download/brag-darwin-amd64
chmod +x brag-darwin-amd64
sudo mv brag-darwin-amd64 /usr/local/bin/brag
```

**Windows:**
Download `brag-windows-amd64.exe` from the [releases page](https://github.com/dkrichards86/brag/releases) and add it to your PATH.

### Using go install

```bash
go install github.com/dkrichards86/brag@latest
```

**Note:** This installs the binary to `$GOPATH/bin` (typically `~/go/bin`). Ensure this directory is in your `PATH`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

Add this line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) to make it permanent.

### From Source

```bash
git clone https://github.com/dkrichards86/brag.git
cd brag
make build
sudo mv brag /usr/local/bin/
```

## Usage

### Adding Wins

Add a win with the current timestamp. **Note:** Use quotes to include hashtags, as `#` starts a comment in most shells:

```bash
brag "Fixed authentication bug #work #backend"
```

You can also add wins without tags (no quotes needed):

```bash
brag Completed code review
```

Add a win with a custom date (date must be the last argument):

```bash
brag "Finished the project setup #project" 2026/01/08
```

Supported date formats:
- `YYYY/MM/DD` (e.g., `2026/01/08`)
- `YYYY-MM-DD` (e.g., `2026-01-08`)
- `MM/DD/YYYY` (e.g., `01/08/2026`)
- `MM-DD-YYYY` (e.g., `01-08-2026`)

### Reviewing Wins

Review wins from the last 7 days (default):

```bash
brag review
```

Review wins within a specific date range:

```bash
brag review --from 2026/01/01 --to 2026/01/09
```

Review only wins with tags (your "macro" wins):

```bash
brag review --tagged
```

Filter by a specific tag:

```bash
brag review --tag work
```

You can combine filters:

```bash
brag review --from 2026/01/01 --tag work
```

### Listing Wins

List all wins with line numbers:

```bash
brag list
```

List wins within a date range:

```bash
brag list --from 2026/01/01 --to 2026/01/09
```

### Editing Wins

Open the entire wins file in your `$EDITOR`:

```bash
brag edit
```

Edit a specific win by line number:

```bash
brag edit 3
```

### Deleting Wins

Delete a win by line number (with confirmation):

```bash
brag delete 2
```

### Searching Wins

Search for wins containing a keyword:

```bash
brag search authentication
```

Search for multiple words:

```bash
brag search "fixed bug"
```

## File Storage

Wins are stored in a plain text file at `~/.brag/wins.txt`.

The format is simple and human-readable:

```
2026-01-09 11:19 | Fixed authentication bug #work #backend
2026-01-09 11:19 | Helped teammate debug their issue #mentoring
2026-01-08 00:00 | Finished the project setup #project
2026-01-07 14:30 | Completed code review
```

You can manually edit this file if needed.

## Macro Wins

Any win with tags is considered a "macro win" - something worth highlighting. Tag important achievements to easily filter them for year-end reviews or performance evaluations:

```bash
brag "Shipped major feature to production #work #backend"
brag "Mentored junior developer on best practices #mentoring #leadership"
```

Then review all tagged wins:

```bash
brag review --tagged
```

Or review by specific tag:

```bash
brag review --tag leadership
```

## Tips

1. **Make it a habit**: Add wins throughout the day as you accomplish things
2. **Be specific**: Include enough detail to remember what you did
3. **Use consistent tags**: Develop a tagging system that works for you (e.g., `#work`, `#learning`, `#mentoring`)
4. **Review regularly**: Look back at your wins weekly or monthly to stay motivated
5. **Tag important wins**: Tag significant achievements so you can filter them later with `--tagged` for year-end reporting
6. **Use quotes for hashtags**: Remember to quote your message when including tags: `brag "message #tag"`

## License

MIT License - see LICENSE file for details.

## Development

### Building from Source

```bash
# Build for your platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint

# Run full CI pipeline locally
make ci
```

### Creating a Release

For maintainers, to create a new release:

```bash
make release VERSION=v1.0.0
```

This will create and push a git tag, which triggers GitHub Actions to build and publish binaries. See [RELEASING.md](RELEASING.md) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

1. Fork and clone the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Submit a pull request
