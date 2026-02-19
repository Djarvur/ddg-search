## Why

The ddg-search project currently provides web search via DuckDuckGo (no API key required) and page content fetching. However, DuckDuckGo's free search has limitations in terms of result quality and rate limiting. Perplexity API provides higher-quality search results with AI-powered summaries and better content understanding. Adding Perplexity search as an alternative CLI command gives users more options for web search with different trade-offs (API key required vs. free service).

## What Changes

- Create new CLI command `perplexity-search` at `cmd/perplexity-search/`
- Implement Perplexity API client in `internal/perplexity/` package
- Add environment variable support for `PERPLEXITY_API_KEY` via `.env` file
- Add corresponding agent skill at `skills/perplexity-search/` with `skill.yaml` and `SKILL.md`
- Update `README.md` with documentation for the new command

## Capabilities

### New Capabilities
- `perplexity-search`: CLI command for searching the web using Perplexity API
- `perplexity-search-skill`: Claude Code skill providing Perplexity search functionality through tool definitions

### Modified Capabilities
- None - this is a new command and skill

## Impact

- **New files**: 
  - `cmd/perplexity-search/main.go`
  - `internal/perplexity/client.go`
  - `internal/perplexity/search.go`
  - `skills/perplexity-search/skill.yaml`
  - `skills/perplexity-search/SKILL.md`
- **Modified files**: 
  - `README.md` (add documentation)
  - `.env.example` (already contains PERPLEXITY_API_KEY)
- **Dependencies**: 
  - Requires Perplexity API key (user-provided via environment variable)
  - May require additional Go dependencies for HTTP client and JSON handling
- **Systems**: 
  - Integrates with existing CLI structure using urfave/cli
  - Integrates with Claude Code's skill system
