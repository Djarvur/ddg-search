# Perplexity Search Skill

Web search using Perplexity API with AI-powered results and citations.

## Installation

This skill requires the `perplexity-search` binary and a Perplexity API key.

### Install the binary

```bash
go install github.com/Djarvur/ddg-search/cmd/perplexity-search@latest
```

### Configure API Key

Export the API key in your shell environment:

```bash
export PERPLEXITY_API_KEY="your-api-key-here"
```

The binary reads the key from the environment only; it does not load `.env`
files itself. To keep the key in a `.env` file, source it first
(`set -a; . ./.env; set +a`) or use a runner such as `direnv` or `dotenv`.

### Get a Perplexity API Key

1. Visit https://www.perplexity.ai/settings/api
2. Sign up or log in to your Perplexity account
3. Generate a new API key
4. Copy the key and export it as `PERPLEXITY_API_KEY`

### Verify installation

```bash
which perplexity-search

# Test the binary
perplexity-search --help

# Test with a simple query (requires API key)
perplexity-search "What is Go programming language?"
```

## Tools

### perplexity-search

Search the web using Perplexity API for high-quality, AI-powered search results with citations.

**Parameters:**
- `query` (required): The search query string
- `max-results` (optional): Maximum number of sources listed with the answer (default: 5)
- `model` (optional): Perplexity model to use (default: sonar)

**Usage:**
```
perplexity-search "golang http client best practices"
```

**Output:** AI-generated answer, followed by a `## Sources` list of the
sources it was grounded in.

**Options:**
- `--max-results int`: Maximum number of sources listed with the answer (default: 5)
- `--model string`: Perplexity model to use (default: sonar)
- `--debug`: Enable debug logging to stderr

**Available Models:**
- `sonar`: Fast and low cost (default)
- `sonar-pro`: Higher quality, more expensive
- `sonar-reasoning`: Adds a reasoning pass
- `sonar-reasoning-pro`: Higher-quality reasoning
- `sonar-deep-research`: Long-running, exhaustive research

## Usage Examples

### Basic search
```
perplexity-search "rust async programming guide"
```

### Limit the number of sources listed
```
perplexity-search --max-results 10 "docker compose best practices"
```

### Search with specific model
```
perplexity-search --model sonar-pro "machine learning fundamentals"
```

### Search with debug mode
```
perplexity-search --debug "kubernetes deployment strategies"
```

## Troubleshooting

### "command not found: perplexity-search"

The binary is not in your PATH. Ensure `$(go env GOPATH)/bin` is in your PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### "PERPLEXITY_API_KEY is not set"

Export the key in the shell that runs the binary:
`export PERPLEXITY_API_KEY="your-key"`. A `.env` file is not read
automatically -- source it first if that is where you keep the key.

### "invalid API key"

Your API key is incorrect or expired. Verify:
1. The key is correctly copied from https://www.perplexity.ai/settings/api
2. The key hasn't been revoked
3. The key has valid credits remaining

### "rate limited by Perplexity API"

You've exceeded the API rate limit. Wait a few minutes before retrying, or upgrade your Perplexity plan for higher limits.

### "payment required - API quota exceeded"

Your API quota has been exhausted. Check your usage at https://www.perplexity.ai/settings/api and consider upgrading your plan.

### Network errors

If you encounter network errors:
1. Check your internet connection
2. Verify the Perplexity API is accessible (https://api.perplexity.ai)
3. Try again with `--debug` flag for more information

## Comparison with ddg-search

| Feature | perplexity-search | ddg-search |
|---------|------------------|-------------|
| API Key Required | Yes | No |
| Result Quality | AI-powered, higher quality | Standard search results |
| Citations | Yes, with sources | No |
| Rate Limits | Based on API tier | DuckDuckGo rate limits |
| Cost | Free tier available, paid for higher limits | Free |
| Speed | Slightly slower (API call) | Faster (HTML scraping) |

Choose `perplexity-search` for higher-quality, AI-summarized results with citations. Choose `ddg-search` for free, fast searches without API key requirements.
