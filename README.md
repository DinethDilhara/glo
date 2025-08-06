# glo

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/DinethDilhara/glo)](https://github.com/DinethDilhara/glo/releases)

A powerful CLI tool to explore Git with beautiful formatting and advanced filtering capabilities. Transform your git log experience with colorized output, JSON exports, and Markdown documentation generation.

## Features

- **Beautiful colored terminal output** with syntax highlighting
- **Multiple export formats**: JSON, Markdown, and colorized text
- **Advanced filtering**: by author, date range, and commit messages
- **Flexible output options**: summary views, tables, and detailed listings
- **Fast and lightweight** with zero external dependencies
- **Easy to use** with intuitive command-line interface

## Installation

### Quick Install (Recommended)

#### macOS / Linux

```bash
# Download and install the latest release
curl -L https://github.com/DinethDilhara/glo/releases/latest/download/glo_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv glo /usr/local/bin/
```

#### Windows (PowerShell)

```powershell
# Download the latest Windows release
Invoke-WebRequest -Uri "https://github.com/DinethDilhara/glo/releases/latest/download/glo_Windows_x86_64.zip" -OutFile "glo.zip"
Expand-Archive glo.zip -DestinationPath .
# Add glo.exe to your PATH
```

### Manual Download

Download the latest binary from [GitHub Releases](https://github.com/DinethDilhara/glo/releases) for your platform:

| Platform    | Architecture  | Download                                                                                                           |
| ----------- | ------------- | ------------------------------------------------------------------------------------------------------------------ |
| **Linux**   | x86_64        | [glo_Linux_x86_64.tar.gz](https://github.com/DinethDilhara/glo/releases/latest/download/glo_Linux_x86_64.tar.gz)   |
| **Linux**   | ARM64         | [glo_Linux_arm64.tar.gz](https://github.com/DinethDilhara/glo/releases/latest/download/glo_Linux_arm64.tar.gz)     |
| **macOS**   | Intel         | [glo_Darwin_x86_64.tar.gz](https://github.com/DinethDilhara/glo/releases/latest/download/glo_Darwin_x86_64.tar.gz) |
| **macOS**   | Apple Silicon | [glo_Darwin_arm64.tar.gz](https://github.com/DinethDilhara/glo/releases/latest/download/glo_Darwin_arm64.tar.gz)   |
| **Windows** | x86_64        | [glo_Windows_x86_64.zip](https://github.com/DinethDilhara/glo/releases/latest/download/glo_Windows_x86_64.zip)     |

## Installation

### Homebrew (Recommended)

```bash
brew install dinethdhilhara/tap/glo
```

### Download Binary

Download the latest binary from [GitHub Releases](https://github.com/DinethDilhara/glo/releases).

### From Source

```bash
git clone https://github.com/DinethDilhara/glo.git
cd glo
go build -o glo .
```

### Homebrew (Coming Soon)

```bash
# Will be available soon
brew install dinethdhilhara/tap/glo
```

### Shell Completions

Enable tab completion for your shell:

```bash
# Install completions for your shell
cd completions/
chmod +x install.sh
./install.sh

# Or manually for bash
source completions/glo_completion.sh

# Or manually for zsh
cp completions/_glo ~/.zsh/completions/
```

**Supported shells:** Bash, Zsh, Fish, PowerShell

## Quick Start

```bash
# Show recent commits with beautiful colors
glo log

# Filter by author
glo log --author="John Doe"

# Export to JSON
glo log --format=json

# Generate Markdown documentation
glo log --format=markdown

# Show commits from last week
glo log --since="1 week ago"
```

## Usage

### Basic Commands

```bash
glo log                    # Show recent commits with colors
glo log --help            # Show all available options
glo --help                # Show global help
```

### Filtering Options

| Flag        | Short | Description              | Example                   |
| ----------- | ----- | ------------------------ | ------------------------- |
| `--author`  | `-a`  | Filter by author name    | `glo log -a "John Doe"`   |
| `--since`   | `-s`  | Show commits since date  | `glo log -s "2024-01-01"` |
| `--until`   | `-u`  | Show commits until date  | `glo log -u "2024-12-31"` |
| `--message` | `-m`  | Filter by commit message | `glo log -m "fix"`        |
| `--limit`   | `-l`  | Limit number of commits  | `glo log -l 10`           |

### Output Formats

| Format     | Description                         | Example                     |
| ---------- | ----------------------------------- | --------------------------- |
| `color`    | Colorized terminal output (default) | `glo log --format=color`    |
| `json`     | JSON format for programmatic use    | `glo log --format=json`     |
| `markdown` | Markdown format for documentation   | `glo log --format=markdown` |

### Additional Options

| Flag        | Description                  | Example                             |
| ----------- | ---------------------------- | ----------------------------------- |
| `--table`   | Output markdown as table     | `glo log --format=markdown --table` |
| `--summary` | Show summary with statistics | `glo log --summary`                 |
| `--verbose` | Enable verbose output        | `glo log --verbose`                 |

## Examples

### Basic Usage

```bash
# Show last 5 commits with colors
glo log --limit=5

# Find all commits by specific author
glo log --author="Alice Smith"

# Show commits from last month
glo log --since="1 month ago"
```

### Advanced Filtering

```bash
# Combine multiple filters
glo log --author="John" --since="2024-01-01" --message="bug"

# Show commits in date range
glo log --since="2024-01-01" --until="2024-06-30"

# Search for specific keywords in commit messages
glo log --message="feature" --limit=20
```

### Export Examples

```bash
# Export recent commits to JSON
glo log --format=json --limit=10 > commits.json

# Generate markdown documentation
glo log --format=markdown --summary > CHANGELOG.md

# Create markdown table of commits
glo log --format=markdown --table --limit=20 > commits-table.md
```

### Summary and Statistics

```bash
# Show repository summary with author statistics
glo log --summary

# Show summary in JSON format
glo log --summary --format=json

# Generate markdown summary report
glo log --summary --format=markdown
```

## Output Examples

### Colorized Terminal Output

```
6bb59b7d Dineth De Alwis 2025-07-28 12:30:58 +0530 chore(core): initialize Go module and entrypoint
60efde4a Dineth De Alwis 2025-07-28 12:30:37 +0530 chore(release): configure GoReleaser for multi-arch builds
d5580101 Dineth De Alwis 2025-07-28 12:29:59 +0530 chore: add .gitignore for build and OS artifacts
```

### JSON Export

```json
[
  {
    "hash": "6bb59b7d6dd409d01606095af4c1c623dd07ce74",
    "author": "Dineth De Alwis",
    "date": "2025-07-28 12:30:58 +0530",
    "message": "chore(core): initialize Go module and entrypoint"
  }
]
```

### Markdown Output

```markdown
# Git Commit History

_Generated on 2025-07-28 15:22:55_

**Total Commits:** 4

---

### 1. chore(core): initialize Go module and entrypoint

- **Hash:** `6bb59b7d`
- **Author:** Dineth De Alwis
- **Date:** 2025-07-28 12:30:58 +0530
```

## Architecture

```
glo/
├── .github/               # GitHub workflows and templates
│   ├── workflows/        # CI/CD automation
│   │   ├── build.yml     # Build, test, and quality checks
│   │   └── release.yml   # Automated releases with GoReleaser
│   └── pull_request_template.md
├── build/                # Build artifacts and releases
├── cmd/                  # CLI commands and flags
│   ├── root.go          # Root command and global flags
│   └── log.go           # Log command implementation
├── completions/          # Shell completion scripts
│   ├── glo_completion.sh # Bash completion
│   ├── _glo             # Zsh completion
│   ├── glo_completion.ps1 # PowerShell completion
│   ├── install.sh       # Completion installation script
│   └── README.md        # Completion documentation
├── internal/
│   ├── formatter/       # Output formatters
│   │   ├── color.go    # ANSI color formatting
│   │   ├── json.go     # JSON output formatting
│   │   └── markdown.go # Markdown formatting
│   ├── gitexec/        # Git command execution
│   │   └── executor.go # Git operations wrapper
│   ├── models/         # Data structures
│   │   └── commit.go   # Commit model definition
│   └── parser/         # Data parsing utilities
│       └── parser.go   # Git output parsing
├── scripts/             # Development and release scripts
│   ├── pre-commit      # Git pre-commit hook
│   ├── release.sh      # Release automation script
│   └── setup-dev.sh    # Development environment setup
├── tests/              # Test suites
│   └── integration/    # Integration tests
│       ├── benchmark_test.go # Performance benchmarks
│       ├── cli_test.go      # CLI integration tests
│       └── git_test.go      # Git operations tests
├── .gitignore          # Git ignore patterns
├── .golangci.yml       # Go linting configuration
├── .goreleaser.yml     # GoReleaser configuration
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── LICENSE             # MIT license
├── Makefile            # Build automation
├── README.md           # Project documentation
└── main.go            # Application entry point
```

## Development

### Prerequisites

- Go 1.19 or higher
- Git installed on your system

### Building from Source

```bash
git clone https://github.com/DinethDilhara/glo.git
cd glo
go mod tidy
go build -o glo .
```

### Running Tests

```bash
go test ./...
```

### Building for Multiple Platforms

```bash
# Build for current platform
go build -o glo .

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o glo-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o glo-windows-amd64.exe .
GOOS=darwin GOARCH=amd64 go build -o glo-darwin-amd64 .
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Inspired by the need for better git log visualization tools
- Thanks to all contributors and users

## Project Status

- Core functionality implemented
- Multiple output formats supported
- Advanced filtering capabilities
- Cross-platform compatibility
- Additional features in development

---

Made with ❤️ by [Dineth De Alwis](https://github.com/DinethDilhara)
