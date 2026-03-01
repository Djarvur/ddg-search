# Inconsistencies and Assumptions for add-mcp-server Change

This document lists all inconsistencies found in the OpenSpec change artifacts and the reasonable assumptions made to resolve them for implementation.

---

## Binary & Naming

### Inconsistency 1: Directory name vs binary name
- **Issue**: User asked about `ddg-search-mcp-server` but directory is `add-mcp-server` and binary is named `ddg-search-mcp`
- **Assumption**: Use `ddg-search-mcp` as the binary name (as specified in proposal.md and design.md)

---

## Configuration

### Inconsistency 2: No config file example provided
- **Issue**: Configuration structure is described but no YAML example is provided
- **Assumption**: Config file structure:
  ```yaml
  server:
    transport: stdio  # or "tcp"
    port: 9100
    tls:
      enabled: false
      cert_file: ""
      key_file: ""
      combined_file: ""
      ca_file: ""
  search:
    max_results: 10
    region: us-en
    time: ""
    safe_search: false
  perplexity:
    api_key: ""
    model: sonar-medium-online
    enabled: true
  fetch:
    timeout: 30s
    user_agent: ddg-search-mcp/1.0
  output:
    format: text  # or "json"
  logging:
    level: info  # debug, info, warn, error
  ```

### Inconsistency 3: Environment variable names not specified
- **Assumption**: Environment variables follow `DDG_MCP_` prefix:
  - `DDG_MCP_TRANSPORT` - transport type
  - `DDG_MCP_PORT` - TCP port
  - `DDG_MCP_TLS_ENABLED` - enable TLS
  - `DDG_MCP_TLS_CERT_FILE` - TLS cert file path
  - `DDG_MCP_TLS_KEY_FILE` - TLS key file path
  - `DDG_MCP_TLS_COMBINED_FILE` - combined TLS file path
  - `DDG_MCP_TLS_CA_FILE` - CA cert file path for mTLS
  - `DDG_MCP_SEARCH_MAX_RESULTS` - max search results
  - `DDG_MCP_SEARCH_REGION` - search region
  - `DDG_MCP_SEARCH_TIME` - time filter
  - `DDG_MCP_SEARCH_SAFE_SEARCH` - enable safe search
  - `DDG_MCP_PERPLEXITY_API_KEY` - Perplexity API key
  - `DDG_MCP_PERPLEXITY_MODEL` - Perplexity model
  - `DDG_MCP_PERPLEXITY_ENABLED` - enable Perplexity
  - `DDG_MCP_FETCH_TIMEOUT` - fetch timeout
  - `DDG_MCP_FETCH_USER_AGENT` - fetch user agent
  - `DDG_MCP_OUTPUT_FORMAT` - output format
  - `DDG_MCP_LOG_LEVEL` - log level

### Inconsistency 4: CLI flag names not specified
- **Assumption**: CLI flags follow kebab-case:
  - `--transport` - transport type
  - `--port` - TCP port
  - `--tls-enabled` - enable TLS
  - `--tls-cert-file` - TLS cert file path
  - `--tls-key-file` - TLS key file path
  - `--tls-combined-file` - combined TLS file path
  - `--tls-ca-file` - CA cert file path for mTLS
  - `--search-max-results` - max search results
  - `--search-region` - search region
  - `--search-time` - time filter
  - `--search-safe-search` - enable safe search
  - `--perplexity-api-key` - Perplexity API key
  - `--perplexity-model` - Perplexity model
  - `--perplexity-enabled` - enable Perplexity
  - `--fetch-timeout` - fetch timeout
  - `--fetch-user-agent` - fetch user agent
  - `--output-format` - output format
  - `--log-level` - log level
  - `--config` - config file path

### Inconsistency 5: Config file permissions validation
- **Issue**: Risk mitigation mentions checking file permissions but no spec requirement
- **Assumption**: Log a warning if config file has permissions more permissive than 0600, but don't fail

---

## Search Tool

### Inconsistency 6: Valid values for `region` parameter
- **Issue**: Spec mentions "us-en", "uk-en" as examples but doesn't list all valid values
- **Assumption**: Accept any region string, pass through to DuckDuckGo/Perplexity without validation

### Inconsistency 7: Valid values for `time` parameter
- **Issue**: Spec mentions "d", "w", "m", "y" but doesn't specify exact format
- **Assumption**: Accept single character values: "d" (day), "w" (week), "m" (month), "y" (year)

### Inconsistency 8: Default values for search parameters
- **Issue**: Spec mentions defaults exist but doesn't specify all of them
- **Assumption**:
  - `max_results`: 10
  - `region`: "us-en"
  - `time`: "" (no filter)
  - `safe_search`: false

### Inconsistency 9: Maximum value for `max_results`
- **Issue**: No maximum specified
- **Assumption**: Maximum of 100 results to prevent abuse

### Inconsistency 10: Valid values for `safe_search`
- **Issue**: Not specified if boolean or levels
- **Assumption**: Boolean (true/false)

### Inconsistency 11: `site` parameter compatibility
- **Issue**: Not specified if `site` works with both providers
- **Assumption**: `site` parameter works with DuckDuckGo only (Perplexity doesn't support site filtering in the same way)

### Inconsistency 12: JSON output format structure
- **Issue**: Spec says "include all result fields" but doesn't define exact structure
- **Assumption**:
  - DuckDuckGo:
    ```json
    {
      "results": [
        {
          "title": "...",
          "url": "...",
          "snippet": "..."
        }
      ],
      "provider": "duckduckgo"
    }
    ```
  - Perplexity:
    ```json
    {
      "answer": "...",
      "citations": ["..."],
      "provider": "perplexity"
    }
    ```

### Inconsistency 13: `model` parameter when Perplexity disabled
- **Issue**: Not specified what happens if `model` is provided but Perplexity is disabled
- **Assumption**: Ignore `model` parameter if Perplexity is disabled or no API key is configured

---

## Fetch Tool

### Inconsistency 14: Default timeout value
- **Issue**: Not specified
- **Assumption**: 30 seconds

### Inconsistency 15: Timeout parameter format
- **Issue**: Not specified if "30s", "30", or "30000ms"
- **Assumption**: Accept Go duration format (e.g., "30s", "1m", "500ms")

### Inconsistency 16: Default user agent
- **Issue**: Not specified
- **Assumption**: "ddg-search-mcp/1.0"

### Inconsistency 17: JSON output metadata
- **Issue**: Spec says "include metadata" but doesn't specify what
- **Assumption**:
  ```json
  {
    "url": "...",
    "status_code": 200,
    "content_type": "text/html",
    "markdown": "..."
  }
  ```

### Inconsistency 18: Redirect limit
- **Issue**: Not specified
- **Assumption**: 10 redirects (Go HTTP client default)

---

## Perplexity

### Inconsistency 19: Model validation
- **Issue**: "Expose all models" but no validation specified
- **Assumption**: Accept any model name string, pass to Perplexity API without validation

### Inconsistency 20: API key format validation
- **Issue**: Not specified
- **Assumption**: No format validation, just check if non-empty string

---

## Logging & Request Handling

### Inconsistency 21: Request ID generation
- **Issue**: Open question says "Yes, generate UUID" but not specified when
- **Assumption**: Generate UUID for every request, include in logs

### Inconsistency 22: Fallback indication location
- **Issue**: Spec says "response SHALL indicate" but not clear if in response body or log
- **Assumption**: Both - log warning and include in response metadata

### Inconsistency 23: Server name for logging
- **Issue**: Not specified
- **Assumption**: "ddg-search-mcp"

### Inconsistency 24: Server version source
- **Issue**: Not specified
- **Assumption**: Use build tag or hardcoded "dev" if not set

---

## MCP Protocol

### Inconsistency 25: MCP protocol version
- **Issue**: Not specified
- **Assumption**: Use latest version supported by mark3labs/mcp-go library

### Inconsistency 26: Tool descriptions
- **Issue**: Not specified
- **Assumption**:
  - `search`: "Perform web search with Perplexity (if configured) falling back to DuckDuckGo"
  - `fetch`: "Fetch a web page and convert to markdown"

### Inconsistency 27: Input schema format
- **Issue**: Not specified
- **Assumption**: Use JSON Schema as required by MCP protocol

### Inconsistency 28: Error response format
- **Issue**: Not specified
- **Assumption**: Use MCP error response format with code and message fields

---

## Signal Handling & Shutdown

### Inconsistency 29: SIGTERM handling
- **Issue**: Spec mentions SIGINT but not SIGTERM
- **Assumption**: Handle SIGTERM the same way as SIGINT (immediate exit)

### Inconsistency 30: Existing connections on shutdown
- **Issue**: Not specified what happens to in-flight requests
- **Assumption**: Close connections immediately, in-flight requests may fail

### Inconsistency 31: Config reload failure notification
- **Issue**: Not specified if clients should be notified
- **Assumption**: Log error but don't notify clients (no mechanism in MCP protocol)

---

## TLS

### Inconsistency 32: Certificate validation
- **Issue**: Spec says "validate that TLS files contain valid certificates" but not what validation
- **Assumption**: Check if files can be loaded by Go's TLS library, don't check expiration or validity period

### Inconsistency 33: mTLS client without certificate
- **Issue**: Not specified behavior
- **Assumption**: Reject connection immediately with TLS handshake error

---

## Additional Inconsistencies Found

### Inconsistency 34: Logging level behavior
- **Issue**: design.md says "Info level: Log info/warn/error, successful requests at debug (not shown)" but spec says "log successful tool calls at debug level (not shown)" - these are consistent but could be clearer
- **Assumption**: Successful requests are logged at debug level, so they only appear when log level is "debug"

### Inconsistency 35: Bad request logging level
- **Issue**: design.md says "Bad requests always logged at error level" but spec says "log bad requests at error level" - consistent
- **Assumption**: Bad requests are always logged at error level regardless of configured log level

### Inconsistency 36: Config reload atomicity
- **Issue**: design.md mentions "atomic config swap" but no spec requirement
- **Assumption**: Implement atomic config swap using pointer replacement

### Inconsistency 37: Test coverage target
- **Issue**: design.md says "50%+" but tasks.md says "Verify test coverage is 50%+" - consistent
- **Assumption**: Target 50%+ code coverage

### Inconsistency 38: E2E test with real MCP client
- **Issue**: Not specified which MCP client to use
- **Assumption**: Use mark3labs/mcp-go client for E2E tests

---

## Summary

Total inconsistencies found: 38

All have been resolved with reasonable assumptions that:
1. Follow Go conventions
2. Align with existing codebase patterns
3. Are consistent with MCP protocol standards
4. Provide sensible defaults
5. Allow for future extensibility

These assumptions should be reviewed and adjusted if they don't match the intended behavior.
