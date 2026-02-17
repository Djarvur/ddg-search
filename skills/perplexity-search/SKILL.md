# Perplexity Search Skill

Web search using Perplexity API with AI-powered results and citations.

## Installation

This skill requires the `perplexity-search` binary and a Perplexity API key.

### Install the binary

```bash
go install github.com/Djarvur/ddg-search/cmd/perplexity-search@latest
```

### Configure API Key

Create a `.env` file in your project directory or set the environment variable:

```bash
# In .env file
PERPLEXITY_API_KEY="your-api-key-here"

# Or set in your shell
export PERPLEXITY_API_KEY="your-api-key-here"
```

### Get a Perplexity API Key

1. Visit https://www.perplexity.ai/settings/api
2. Sign up or log in to your Perplexity account
3. Generate a new API key
4. Copy the key and add it to your `.env` file or set it as an environment variable

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
- `max-results` (optional): Maximum number of results to return (default: 5)
- `model` (optional): Perplexity model to use (default: sonar-medium-online)

**Usage:**
```
perplexity-search "golang http client best practices"
```

**Output:** AI-generated answer with citations to sources.

**Options:**
- `--max-results int`: Maximum number of results (default: 5)
- `--model string`: Perplexity model to use (default: sonar-medium-online)
- `--debug`: Enable debug logging to stderr

**Available Models:**
- `sonar-medium-online`: Balanced performance and quality (default)
- `sonar-small-online`: Faster, lower cost
- `sonar-pro-online`: Higher quality, more expensive

## Usage Examples

### Basic search
```
perplexity-search "rust async programming guide"
```

### Search with custom max results
```
perplexity-search --max-results 10 "docker compose best practices"
```

### Search with specific model
```
perplexity-search --model sonar-pro-online "machine learning fundamentals"
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

### "PERPLEXITY_API_KEY environment variable not set"

You need to configure your Perplexity API key. Either:
1. Create a `.env` file with `PERPLEXITY_API_KEY="your-key"`
2. Set the environment variable: `export PERPLEXITY_API_KEY="your-key"`

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
