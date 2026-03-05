# Stage 4 Manual Test Plan

## Prerequisites

- Go 1.21+ installed
- Project cloned and dependencies installed (`go mod download`)
- Terminal access for running commands
- Text editor for creating config files

## Test Cases

### Test Case 1: Basic Startup

**Description:** Verify the application starts successfully with default configuration.

**Steps:**
1. Open a terminal
2. Navigate to the project directory
3. Run: `go run ./cmd/ddg-search-mcp`
4. Observe the output
5. Press Ctrl-C to stop the application

**Expected Result:**
- Application starts without errors
- Output shows "Starting ddg-search-mcp"
- Output shows "Configuration loaded" with default values
- Output shows "Server ready, waiting for signals..."
- Application stops cleanly on Ctrl-C with "Shutting down..." and "Shutdown complete"

### Test Case 2: Configuration File Loading

**Description:** Verify the application loads configuration from a YAML file.

**Steps:**
1. Create directory: `mkdir -p ~/.config/ddg-search`
2. Create config file: `~/.config/ddg-search/config.yaml` with content:
   ```yaml
   logging:
     level: debug
   search:
     max_results: 20
     safe_search: true
   ```
3. Run: `go run ./cmd/ddg-search-mcp`
4. Observe the output
5. Press Ctrl-C to stop

**Expected Result:**
- Application starts successfully
- Output shows "Configuration loaded" with custom values
- Output shows `level=debug` and `max_results=20`
- Application stops cleanly

### Test Case 3: Environment Variable Overrides

**Description:** Verify environment variables override config file values.

**Steps:**
1. Ensure config file exists from Test Case 2
2. Run: `DDG_SEARCH_LOGGING_LEVEL=debug DDG_SEARCH_SEARCH_MAX_RESULTS=15 go run ./cmd/ddg-search-mcp`
3. Observe the output
4. Press Ctrl-C to stop

**Expected Result:**
- Application starts successfully
- Output shows `level=debug` (env override)
- Output shows `max_results=15` (env override)
- Application stops cleanly

### Test Case 4: CLI Parameter Overrides

**Description:** Verify CLI parameters override environment variables and config file.

**Steps:**
1. Run: `go run ./cmd/ddg-search-mcp --log-level=debug`
2. Observe the output
3. Press Ctrl-C to stop

**Expected Result:**
- Application starts successfully
- Output shows `level=debug` (CLI override)
- Application stops cleanly

### Test Case 5: Configuration Priority

**Description:** Verify the priority order: CLI > env > config file > defaults.

**Steps:**
1. Create config file with `level: info`
2. Run: `DDG_SEARCH_LOGGING_LEVEL=debug go run ./cmd/ddg-search-mcp --log-level=warn`
3. Observe the output
4. Press Ctrl-C to stop

**Expected Result:**
- Application starts successfully
- Output shows `level=warn` (CLI flag has highest priority)
- Application stops cleanly

### Test Case 6: Invalid Configuration

**Description:** Verify the application fails gracefully with invalid configuration.

**Steps:**
1. Create config file with invalid content:
   ```yaml
   logging:
     level: invalid
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output

**Expected Result:**
- Application fails to start
- Error message indicates "invalid configuration"
- Error message indicates the specific validation failure

### Test Case 7: Missing Config File

**Description:** Verify the application uses defaults when config file doesn't exist.

**Steps:**
1. Remove or rename config file: `mv ~/.config/ddg-search/config.yaml ~/.config/ddg-search/config.yaml.bak`
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output
4. Press Ctrl-C to stop
5. Restore config file: `mv ~/.config/ddg-search-mcp/config.yaml.bak ~/.config/ddg-search/config.yaml`

**Expected Result:**
- Application starts successfully
- Output shows default configuration values
- No error about missing config file
- Application stops cleanly

### Test Case 8: Invalid YAML

**Description:** Verify the application fails gracefully with invalid YAML syntax.

**Steps:**
1. Create config file with invalid YAML:
   ```yaml
   logging:
     level: info
   invalid: yaml: content
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output

**Expected Result:**
- Application fails to start
- Error message indicates YAML parsing error
- Application exits with non-zero status

### Test Case 9: Log Levels

**Description:** Verify all log levels work correctly.

**Steps:**
1. Run with debug level: `go run ./cmd/ddg-search-mcp --log-level=debug`
2. Observe verbose output
3. Press Ctrl-C to stop
4. Run with info level: `go run ./cmd/ddg-search-mcp --log-level=info`
5. Observe normal output
6. Press Ctrl-C to stop
7. Run with warn level: `go run ./cmd/ddg-search-mcp --log-level=warn`
8. Observe minimal output
9. Press Ctrl-C to stop
10. Run with error level: `go run ./cmd/ddg-search-mcp --log-level=error`
11. Observe only errors
12. Press Ctrl-C to stop

**Expected Result:**
- All log levels work correctly
- Debug level shows most verbose output
- Error level shows only errors
- Application stops cleanly in all cases

### Test Case 10: Version Flag

**Description:** Verify the --version flag works correctly.

**Steps:**
1. Run: `go run ./cmd/ddg-search-mcp --version`
2. Observe the output

**Expected Result:**
- Application prints version information
- Output includes "ddg-search-mcp"
- Application exits immediately

### Test Case 11: Help Flag

**Description:** Verify the --help flag works correctly.

**Steps:**
1. Run: `go run ./cmd/ddg-search-mcp --help`
2. Observe the output

**Expected Result:**
- Application prints help information
- Output includes description
- Output includes usage information
- Output shows available flags (--log-level)
- Application exits immediately

### Test Case 12: Graceful Shutdown (SIGTERM)

**Description:** Verify the application shuts down gracefully on SIGTERM.

**Steps:**
1. Run: `go run ./cmd/ddg-search-mcp`
2. In another terminal, find the process ID: `ps aux | grep ddg-search-mcp`
3. Send SIGTERM: `kill -TERM <pid>`
4. Observe the output in the first terminal

**Expected Result:**
- Application receives SIGTERM
- Output shows "Received shutdown signal"
- Output shows "Shutting down..."
- Output shows "Shutdown complete"
- Application exits cleanly

### Test Case 13: Graceful Shutdown (SIGINT)

**Description:** Verify the application shuts down gracefully on SIGINT (Ctrl-C).

**Steps:**
1. Run: `go run ./cmd/ddg-search-mcp`
2. Press Ctrl-C
3. Observe the output

**Expected Result:**
- Application receives SIGINT
- Output shows "Received shutdown signal"
- Output shows "Shutting down..."
- Output shows "Shutdown complete"
- Application exits cleanly

### Test Case 14: Config Reload (SIGHUP)

**Description:** Verify the application reloads configuration on SIGHUP signal.

**Steps:**
1. Create config file with `level: info`
2. Run: `go run ./cmd/ddg-search-mcp`
3. In another terminal, find the process ID
4. Update config file to `level: debug`
5. Send SIGHUP: `kill -HUP <pid>`
6. Observe the output in the first terminal
7. Press Ctrl-C to stop

**Expected Result:**
- Application receives SIGHUP
- Output shows "Received SIGHUP signal, reloading configuration"
- Output shows "Configuration reloaded successfully"
- Output shows new configuration values
- Application continues running

### Test Case 15: Config Reload with Invalid Config

**Description:** Verify the application keeps previous config on reload failure.

**Steps:**
1. Create config file with valid content
2. Run: `go run ./cmd/ddg-search-mcp`
3. In another terminal, find the process ID
4. Update config file with invalid content
5. Send SIGHUP: `kill -HUP <pid>`
6. Observe the output in the first terminal
7. Press Ctrl-C to stop

**Expected Result:**
- Application receives SIGHUP
- Output shows "Configuration reload failed"
- Output shows "keeping previous configuration"
- Application continues running with previous config
- Application stops cleanly on Ctrl-C

### Test Case 16: Full Configuration

**Description:** Verify the application works with a complete configuration file.

**Steps:**
1. Create complete config file:
   ```yaml
   server:
     protocol: stdio
     bind_address: localhost:9100
     tls:
       enabled: false
       min_version: "1.2"
       mtls:
         enabled: false
   logging:
     level: info
   search:
     max_results: 10
     safe_search: false
   perplexity:
     enabled: false
     access_token: ""
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output
4. Press Ctrl-C to stop

**Expected Result:**
- Application starts successfully
- Output shows all configuration values
- Configuration string shows all sections
- Application stops cleanly

## Cleanup

After testing:
1. Remove test config file: `rm ~/.config/ddg-search/config.yaml`
2. Remove test directory if empty: `rmdir ~/.config/ddg-search`
3. Verify no test processes are running: `ps aux | grep ddg-search-mcp`
