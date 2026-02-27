## ADDED Requirements

### Requirement: TLS configuration
The system SHALL support TLS configuration for HTTP transport.

#### Scenario: TLS enabled
- **WHEN** TLS is enabled in configuration
- **AND** HTTP transport is configured
- **THEN** the system starts an HTTPS server
- **AND** the system uses the configured TLS certificate and key

#### Scenario: TLS disabled
- **WHEN** TLS is disabled in configuration
- **THEN** the system starts an HTTP server without TLS
- **AND** the system does not require TLS certificates

#### Scenario: Default TLS setting
- **WHEN** no TLS setting is specified in configuration
- **THEN** TLS is disabled by default

### Requirement: TLS certificate and key
The system SHALL require TLS certificate and key files when TLS is enabled.

#### Scenario: Valid certificate and key
- **WHEN** TLS is enabled
- **AND** valid certificate and key files are provided
- **THEN** the system loads the certificate and key
- **AND** the server starts with TLS enabled

#### Scenario: Missing certificate file
- **WHEN** TLS is enabled
- **AND** the certificate file is not found
- **THEN** the system logs an error
- **AND** the application exits with a non-zero status

#### Scenario: Missing key file
- **WHEN** TLS is enabled
- **AND** the key file is not found
- **THEN** the system logs an error
- **AND** the application exits with a non-zero status

#### Scenario: Invalid certificate or key
- **WHEN** TLS is enabled
- **AND** the certificate or key file is invalid
- **THEN** the system logs an error
- **AND** the application exits with a non-zero status

### Requirement: mTLS configuration
The system SHALL support mutual TLS (mTLS) for client authentication.

#### Scenario: mTLS enabled
- **WHEN** mTLS is enabled in configuration
- **AND** TLS is enabled
- **THEN** the system requires client certificates
- **AND** the system validates client certificates

#### Scenario: mTLS disabled
- **WHEN** mTLS is disabled in configuration
- **THEN** the system does not require client certificates
- **AND** the server accepts connections without client certificates

#### Scenario: mTLS with CA certificate
- **WHEN** mTLS is enabled
- **AND** a CA certificate is provided
- **THEN** the system validates client certificates against the CA
- **AND** the system rejects invalid client certificates

### Requirement: Client certificate validation
The system SHALL validate client certificates when mTLS is enabled.

#### Scenario: Valid client certificate
- **WHEN** a client connects with a valid certificate
- **AND** mTLS is enabled
- **THEN** the system accepts the connection
- **AND** the system logs the client certificate details

#### Scenario: Invalid client certificate
- **WHEN** a client connects with an invalid certificate
- **AND** mTLS is enabled
- **THEN** the system rejects the connection
- **AND** the system logs the rejection reason

#### Scenario: No client certificate
- **WHEN** a client connects without a certificate
- **AND** mTLS is enabled
- **THEN** the system rejects the connection
- **AND** the system logs that no certificate was provided

### Requirement: TLS configuration parameters
The system SHALL support configurable TLS parameters.

#### Scenario: Certificate path configuration
- **WHEN** TLS is enabled
- **AND** a certificate path is specified in configuration
- **THEN** the system uses the specified certificate file

#### Scenario: Key path configuration
- **WHEN** TLS is enabled
- **AND** a key path is specified in configuration
- **THEN** the system uses the specified key file

#### Scenario: CA certificate path configuration
- **WHEN** mTLS is enabled
- **AND** a CA certificate path is specified in configuration
- **THEN** the system uses the specified CA certificate for validation

#### Scenario: Min TLS version configuration
- **WHEN** a minimum TLS version is specified in configuration
- **THEN** the system enforces the minimum TLS version
- **AND** the system rejects connections with lower TLS versions

### Requirement: TLS connection logging
The system SHALL log TLS connection details at debug level.

#### Scenario: TLS connection established
- **WHEN** a TLS connection is established
- **THEN** the system logs the connection at debug level
- **AND** the log includes TLS version and cipher suite

#### Scenario: mTLS connection established
- **WHEN** an mTLS connection is established
- **THEN** the system logs the connection at debug level
- **AND** the log includes client certificate subject

#### Scenario: TLS handshake failure
- **WHEN** a TLS handshake fails
- **THEN** the system logs the failure at debug level
- **AND** the log includes the failure reason

### Requirement: TLS with HTTP SSE
The system SHALL support TLS with HTTP SSE transport.

#### Scenario: HTTPS SSE connection
- **WHEN** TLS is enabled
- **AND** a client connects to the SSE endpoint via HTTPS
- **THEN** the system establishes a secure SSE connection
- **AND** the system sends SSE messages over TLS

#### Scenario: mTLS with SSE
- **WHEN** mTLS is enabled
- **AND** a client connects to the SSE endpoint with a valid certificate
- **THEN** the system establishes a secure SSE connection
- **AND** the system validates the client certificate

### Requirement: TLS reload
The system SHALL support reloading TLS configuration on HUP signal.

#### Scenario: TLS certificate reloaded
- **WHEN** the server receives SIGHUP signal
- **AND** TLS is enabled
- **AND** the certificate file has changed
- **THEN** the system reloads the certificate
- **AND** the system logs the reload

#### Scenario: TLS key reloaded
- **WHEN** the server receives SIGHUP signal
- **AND** TLS is enabled
- **AND** the key file has changed
- **THEN** the system reloads the key
- **AND** the system logs the reload

#### Scenario: Invalid TLS configuration on reload
- **WHEN** the server receives SIGHUP signal
- **AND** the new TLS configuration is invalid
- **THEN** the system logs an error
- **AND** the system continues with the previous TLS configuration
