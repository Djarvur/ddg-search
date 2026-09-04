## ADDED Requirements

### Requirement: Rate limit detection

The system SHALL detect when DuckDuckGo is rate-limiting requests.

#### Scenario: HTTP 429 triggers rate limit handling

- **WHEN** DuckDuckGo returns HTTP 429 (Too Many Requests)
- **THEN** system recognizes this as a rate-limit condition

#### Scenario: Server error triggers retry

- **WHEN** DuckDuckGo returns HTTP 5xx error
- **THEN** system recognizes this as a retryable condition

#### Scenario: HTTP 202 triggers rate limit handling

- **WHEN** DuckDuckGo returns HTTP 202 (Accepted)
- **THEN** system recognizes this as a rate-limit condition

#### Scenario: Rate limiting is not inferred from page text

- **WHEN** a response body contains words such as "captcha" or "blocked"
- **THEN** system does NOT treat the response as rate-limited on that basis alone
- **AND** rate limiting is detected from the HTTP status only

Scanning body text produced false positives when those words appeared in
legitimate result snippets. Re-enabling it needs a snippet-safe heuristic.

### Requirement: Automatic retry

The system SHALL automatically retry failed requests up to a configurable maximum.

#### Scenario: Retry on rate limit

- **WHEN** a rate-limit condition is detected
- **THEN** system waits for the configured delay
- **AND** system retries the request

#### Scenario: Maximum retries respected

- **WHEN** rate-limit condition persists after maximum retries
- **THEN** system returns an error
- **AND** error message indicates rate limiting was encountered

### Requirement: Exponential backoff

The system SHALL use exponential backoff with jitter for retry delays.

#### Scenario: Delay increases with each retry

- **WHEN** multiple retries are needed
- **THEN** each retry delay is longer than the previous
- **AND** delay is calculated as: `baseDelay * (multiplier ^ attempt) + jitter`

#### Scenario: Delay capped at maximum

- **WHEN** calculated delay exceeds configured maximum delay
- **THEN** system uses the maximum delay instead

### Requirement: Configurable retry parameters

The system SHALL expose configurable retry behavior.

#### Scenario: Max retries configurable

- **WHEN** user specifies maximum retry count
- **THEN** system uses that value as the retry limit

#### Scenario: Base delay configurable

- **WHEN** user specifies base delay
- **THEN** system uses that value as the initial retry delay

#### Scenario: Max delay configurable

- **WHEN** user specifies maximum delay
- **THEN** system caps retry delays at that value

### Requirement: Graceful failure

The system SHALL provide clear error messages when rate limiting cannot be overcome.

#### Scenario: Clear error on max retries exceeded

- **WHEN** maximum retries are exhausted
- **THEN** system outputs a clear error message
- **AND** error indicates rate limiting was the cause
