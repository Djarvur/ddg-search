# Stage 7: Perplexity Integration - Manual Test Plan

## Prerequisites

1. Go 1.23+ installed
2. Network connectivity (for Perplexity API and DuckDuckGo search)
3. Valid Perplexity API access token (for testing Perplexity integration)
4. Optional: MCP Inspector or Claude Desktop for testing

## Known Limitations

### 1. Perplexity API Token Required
Perplexity integration requires a valid API access token. Without a token, the system will fall back to DuckDuckGo search.

**Impact**: Tests that require Perplexity will skip or fall back to DuckDuckGo if no token is provided.

### 2. Network Dependencies
Perplexity integration requires network access to:
- Perplexity API (api.perplexity.ai)
- DuckDuckGo search service (for fallback)

**Impact**: Tests will fail if network is unavailable or if external services are down.

### 3. Rate Limiting
Perplexity may rate limit excessive requests. The implementation does not retry on rate limit errors.

**Impact**: Some Perplexity search requests may fail with rate limit errors during testing.

---

## Test Cases

### Test Case 1: Perplexity Disabled - DuckDuckGo Fallback

**Objective**: Verify that when Perplexity is disabled, DuckDuckGo is used for search.

**Shell Command**:
```bash
# Create config with Perplexity disabled
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: false
  access_token: ""
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search requests use DuckDuckGo
- Results are in markdown format with numbered links
- No Perplexity API calls are made
- No fallback note in response

---

### Test Case 2: Perplexity Enabled Without Token - DuckDuckGo Direct

**Objective**: Verify that when Perplexity is enabled but no token is provided, DuckDuckGo is used directly without attempting Perplexity.

**Shell Command**:
```bash
# Create config with Perplexity enabled but no token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: ""
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search requests use DuckDuckGo directly (no Perplexity attempt)
- Results are in markdown format with numbered links
- No Perplexity API calls are made
- Response includes: "*Search powered by DuckDuckGo*"
- No fallback note in response (since Perplexity was never attempted)

---

### Test Case 3: Perplexity Enabled With Valid Token

**Objective**: Verify that Perplexity search works with a valid token.

**Shell Command**:
```bash
# Create config with Perplexity enabled and valid token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search requests use Perplexity API
- Results are formatted as markdown with:
  - Query as markdown heading
  - Perplexity answer
  - References section with numbered URLs
- No fallback note in response

---

### Test Case 4: Perplexity Invalid Token - DuckDuckGo Fallback

**Objective**: Verify that invalid Perplexity token triggers fallback to DuckDuckGo.

**Shell Command**:
```bash
# Create config with invalid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "invalid_token_12345"
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search requests fall back to DuckDuckGo
- Results are in markdown format with numbered links
- Response includes fallback note: "*Note: Perplexity search failed, falling back to DuckDuckGo.*"

---

### Test Case 5: Perplexity Rate Limit - DuckDuckGo Fallback

**Objective**: Verify that Perplexity rate limit errors trigger fallback to DuckDuckGo.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
EOF

# Send multiple rapid requests to trigger rate limit
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; for i in {1..10}; do printf '{"jsonrpc":"2.0","id":%d,"method":"tools/call","params":{"name":"search","arguments":{"query":"test query %d"}}}\n' "$((i+1))" "$i"; done) | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Initial requests use Perplexity API
- If rate limit is hit, subsequent requests fall back to DuckDuckGo
- Responses that fell back include fallback note
- No retries to Perplexity API

---

### Test Case 6: Perplexity Search With Parameters

**Objective**: Verify that search parameters are passed to Perplexity API.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
search:
  max_results: 5
  safe_search: false
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming","max_results":2,"safe_search":false}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search requests use Perplexity API
- max_results parameter is respected
- safe_search parameter is respected
- Results are formatted as markdown

---

### Test Case 7: Perplexity Result Formatting

**Objective**: Verify that Perplexity results are formatted correctly as markdown.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"what is golang"}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search results include:
  - Query as markdown heading: "# what is golang"
  - Perplexity answer (AI-generated response)
  - References section with numbered URLs
- Format is clean markdown

---

### Test Case 8: Perplexity Network Error - DuckDuckGo Fallback

**Objective**: Verify that network errors trigger fallback to DuckDuckGo.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
EOF

# Note: This test requires blocking Perplexity API access
# On macOS with pfctl:
sudo pfctl -e
sudo pfctl -f /dev/stdin << 'EOF'
block drop from any to api.perplexity.ai
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml

# Cleanup
sudo pfctl -f /etc/pf.conf
```

**Expected Results**:
- Server starts successfully
- Search requests fall back to DuckDuckGo
- Results are in markdown format with numbered links
- Response includes fallback note

---

### Test Case 9: Perplexity Empty Results

**Objective**: Verify that Perplexity search with no results is handled correctly.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"xyzabc123nonexistent"}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search completes successfully
- Results may be empty or minimal
- No errors in server logs

---

### Test Case 10: Perplexity With Safe Search

**Objective**: Verify that safe_search parameter is passed to Perplexity API.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
search:
  safe_search: true
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"test query","safe_search":true}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search requests use Perplexity API
- Results are filtered according to safe search settings
- No errors in server logs

---

### Test Case 11: Perplexity Timeout - DuckDuckGo Fallback

**Objective**: Verify that Perplexity timeout errors trigger fallback to DuckDuckGo.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
EOF

# Note: This test requires simulating a timeout
# On macOS with pfctl, block responses from Perplexity:
sudo pfctl -e
sudo pfctl -f /dev/stdin << 'EOF'
block return out from any to api.perplexity.ai
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}\n') | timeout 5 go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml

# Cleanup
sudo pfctl -f /etc/pf.conf
```

**Expected Results**:
- Server starts successfully
- Search requests fall back to DuckDuckGo
- Results are in markdown format with numbered links
- Response includes fallback note

---

### Test Case 12: Perplexity Transient Error - No Retry

**Objective**: Verify that Perplexity transient errors do not trigger retries.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
EOF

# Note: This test requires simulating a transient error (5xx response)
# On macOS with pfctl, block responses from Perplexity:
sudo pfctl -e
sudo pfctl -f /dev/stdin << 'EOF'
block return out from any to api.perplexity.ai
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}\n') | timeout 5 go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml

# Cleanup
sudo pfctl -f /etc/pf.conf
```

**Expected Results**:
- Server starts successfully
- Search requests fall back to DuckDuckGo immediately (no retry to Perplexity)
- Results are in markdown format with numbered links
- Response includes fallback note

---

### Test Case 13: Perplexity Quota Exceeded - DuckDuckGo Fallback

**Objective**: Verify that Perplexity quota exceeded errors trigger fallback to DuckDuckGo.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
EOF

# Note: This test requires a Perplexity token that has exceeded quota
# If you have such a token, use it; otherwise, this test will fall back to DuckDuckGo

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- If quota is exceeded, search falls back to DuckDuckGo
- Results are in markdown format with numbered links
- Response includes fallback note if fallback occurred

---

### Test Case 14: Perplexity With Max Results Parameter

**Objective**: Verify that max_results parameter is passed to Perplexity API.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
search:
  max_results: 3
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming","max_results":3}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search requests use Perplexity API
- Results are limited to 3 results
- Results are formatted as markdown

---

### Test Case 15: Perplexity With Safe Search Enabled

**Objective**: Verify that safe_search parameter is passed to Perplexity API.

**Shell Command**:
```bash
# Create config with valid Perplexity token
cat > /tmp/test-config.yaml << 'EOF'
perplexity:
  enabled: true
  access_token: "YOUR_PERPLEXITY_TOKEN_HERE"
search:
  safe_search: true
EOF

(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"test query","safe_search":true}}}\n') | go run ./cmd/ddg-search-mcp --config /tmp/test-config.yaml
```

**Expected Results**:
- Server starts successfully
- Search requests use Perplexity API
- Results are filtered according to safe search settings
- No errors in server logs

---

## Cleanup

After testing, remove temporary config files:
```bash
rm -f /tmp/test-config.yaml
```

---

## Notes

1. **Perplexity Token**: Replace `YOUR_PERPLEXITY_TOKEN_HERE` with your actual Perplexity API access token for tests that require it.

2. **Network Access**: Ensure your network allows access to `api.perplexity.ai` and DuckDuckGo.

3. **Rate Limits**: Be mindful of Perplexity API rate limits during testing. The implementation does not retry on rate limit errors.

4. **Fallback Behavior**: The system automatically falls back to DuckDuckGo on any Perplexity error (authentication, rate limit, network, timeout, quota exceeded, etc.) and includes a fallback note in the response.

5. **No Retry Behavior**: The implementation does not retry Perplexity requests on any failure. It immediately falls back to DuckDuckGo.

6. **Test Cases 8, 11, 12**: These tests require macOS with pfctl. On Linux, you can use iptables or temporarily disconnect from the internet to simulate network errors.

7. **Test Case 13**: Requires a Perplexity token that has exceeded its quota. If you don't have such a token, this test will fall back to DuckDuckGo.
