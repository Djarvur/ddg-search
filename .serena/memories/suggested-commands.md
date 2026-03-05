# Suggested Commands for ddg-search

## Development Commands

### Building
```bash
# Build all binaries
go build -o bin/ddg-search ./cmd/ddg-search
go build -o bin/page-dump ./cmd/page-dump
go build -o bin/perplexity-search ./cmd/perplexity-search

# Build specific binary
go build -o bin/ddg-search ./cmd/ddg-search
```

### Testing
```bash
# Run all tests with coverage and race detector
go test -cover -race ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test -v ./internal/search
go test -v ./internal/perplexity
go test -v ./internal/dump
go test -v ./internal/config

# Run tests with coverage profile
go test -coverprofile=coverage.out -covermode=atomic -race ./...

# View coverage report
go tool cover -html=coverage.out
go tool cover -func=coverage.out
```

### Linting
```bash
# Run golangci-lint with auto-fix
golangci-lint run --new-from-merge-base main --whole-files --fix ./...

# Run golangci-lint without auto-fix
golangci-lint run ./...

# Run specific linters
golangci-lint run --disable-all --enable=gofmt,goimports ./...
```

### Using mise (task runner)
```bash
# Run linting
mise run lint

# Run tests
mise run test
```

## Running the Tools

### ddg-search
```bash
# Basic search
ddg-search golang

# With options
ddg-search --max-results 5 --site github.com docker compose
ddg-search --region uk-en premier league
ddg-search --time w news today
ddg-search --max-retries 5 --retry-delay 2s --max-retry-delay 60s slow query
ddg-search --debug golang
```

### page-dump
```bash
# Basic usage
page-dump https://example.com

# With options
page-dump --timeout 60s --user-agent "my-agent/1.0" https://example.com
```

### perplexity-search
```bash
# Set API key first
export PERPLEXITY_API_KEY="your-api-key"

# Basic search
perplexity-search "What is Go programming language?"

# With options
perplexity-search --max-results 3 golang tutorial
perplexity-search --model sonar-pro-online "machine learning fundamentals"
perplexity-search --debug "kubernetes deployment strategies"
```

## System Commands (Darwin/macOS)

```bash
# List files
ls -la

# Find files
find . -name "*.go"
find . -type f -name "*.go" | head -20

# Search in files
grep -r "Searcher" ./internal/
grep -r "func New" ./internal/

# Check git status
git status
git diff
git log --oneline -10

# Change directory
cd /path/to/directory

# Show file contents
cat file.txt
head -n 20 file.txt
tail -n 20 file.txt

# Process management
ps aux | grep ddg-search
kill <pid>
```

## CI/CD Commands

The project uses GitHub Actions for CI/CD:
- **Lint**: Runs golangci-lint v2.9.0
- **Test**: Runs tests with coverage, race detector, and enforces 50% coverage threshold
- **Build**: Builds all binaries and verifies they work with `--help`

## Installation

```bash
# Install from source
go install github.com/Djarvur/ddg-search@latest

# Or build from source
git clone https://github.com/Djarvur/ddg-search
cd ddg-search
make build
```
