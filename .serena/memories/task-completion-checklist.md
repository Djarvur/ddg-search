# Task Completion Checklist for ddg-search

## Before Completing a Task

When you finish implementing a feature or making changes, follow this checklist:

### 1. Code Quality
- [ ] Run linter: `mise run lint` or `golangci-lint run --fix ./...`
- [ ] Fix any linting issues
- [ ] Ensure code follows project conventions (see `code-style-and-conventions.md`)

### 2. Testing
- [ ] Run all tests: `mise run test` or `go test -cover -race ./...`
- [ ] Ensure all tests pass
- [ ] Check coverage is above 50% threshold
- [ ] Add tests for new functionality if applicable
- [ ] Test both success and error paths

### 3. Building
- [ ] Build all binaries:
  ```bash
  go build -o bin/ddg-search ./cmd/ddg-search
  go build -o bin/page-dump ./cmd/page-dump
  go build -o bin/perplexity-search ./cmd/perplexity-search
  ```
- [ ] Verify binaries work: `./bin/ddg-search --help`, `./bin/page-dump --help`, `./bin/perplexity-search --help`

### 4. Manual Testing (if applicable)
- [ ] Test the CLI tools manually with various options
- [ ] Verify error handling works correctly
- [ ] Test rate limiting behavior (if relevant)
- [ ] Test with different inputs and edge cases

### 5. Documentation
- [ ] Update README.md if user-facing changes were made
- [ ] Update or add godoc comments for exported symbols
- [ ] Update OpenSpec specs if using OpenSpec workflow
- [ ] Update skill documentation if skills were modified

### 6. Git
- [ ] Review changes with `git diff`
- [ ] Stage relevant files: `git add <files>`
- [ ] Commit with descriptive message following conventional commits
- [ ] Push changes if working on a branch

## Common Commands Summary

```bash
# Lint
mise run lint
# or
golangci-lint run --fix ./...

# Test
mise run test
# or
go test -cover -race ./...

# Build
go build -o bin/ddg-search ./cmd/ddg-search
go build -o bin/page-dump ./cmd/page-dump
go build -o bin/perplexity-search ./cmd/perplexity-search

# Verify
./bin/ddg-search --help
./bin/page-dump --help
./bin/perplexity-search --help
```

## CI/CD Checks

The CI pipeline will automatically run:
1. **Lint**: golangci-lint v2.9.0
2. **Test**: Tests with coverage, race detector, and 50% coverage threshold
3. **Build**: Build all binaries and verify they work

Ensure all these checks pass before considering a task complete.

## OpenSpec Workflow

If working with OpenSpec changes:
1. Create or update change artifacts (proposal.md, spec.md, design.md, tasks.md)
2. Implement the changes
3. Verify implementation matches specs
4. Archive the change when complete

See `.roo/skills/openspec-*` for detailed OpenSpec workflow instructions.
