## ADDED Requirements

### Requirement: TLS enablement
The system SHALL support enabling TLS for TCP transport via configuration.

#### Scenario: TLS disabled by default
- **WHEN** TLS is not configured
- **THEN** the server SHALL use plain HTTP
- **AND** connections SHALL be unencrypted

#### Scenario: TLS enabled via config
- **WHEN** tls.enabled is set to true in config
- **THEN** the server SHALL use HTTPS
- **AND** connections SHALL be encrypted

### Requirement: Separate key and certificate files
The system SHALL support loading TLS key and certificate from separate files.

#### Scenario: TLS with separate files
- **WHEN** tls.cert_file and tls.key_file are specified
- **THEN** the system SHALL load the certificate from cert_file
- **AND** load the key from key_file
- **AND** use them for TLS connections

#### Scenario: Missing cert file
- **WHEN** tls.cert_file is specified but the file does not exist
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate the certificate file is missing

#### Scenario: Missing key file
- **WHEN** tls.key_file is specified but the file does not exist
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate the key file is missing

### Requirement: Combined key and certificate file
The system SHALL support loading TLS key and certificate from a single combined file.

#### Scenario: TLS with combined file
- **WHEN** tls.combined_file is specified
- **THEN** the system SHALL load both certificate and key from the combined file
- **AND** use them for TLS connections

#### Scenario: Missing combined file
- **WHEN** tls.combined_file is specified but the file does not exist
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate the combined file is missing

### Requirement: TLS file format validation
The system SHALL validate that TLS files contain valid certificates and keys.

#### Scenario: Invalid certificate file
- **WHEN** the certificate file contains invalid data
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate the certificate is invalid

#### Scenario: Invalid key file
- **WHEN** the key file contains invalid data
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate the key is invalid

#### Scenario: Certificate and key mismatch
- **WHEN** the certificate and key do not match
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate the certificate and key do not match

### Requirement: mTLS with CA certificate
The system SHALL support mutual TLS (mTLS) with client certificate verification using a CA certificate.

#### Scenario: mTLS enabled with CA cert
- **WHEN** tls.ca_file is specified
- **THEN** the system SHALL load the CA certificate
- **AND** verify client certificates against the CA
- **AND** reject connections without valid client certificates

#### Scenario: mTLS with valid client certificate
- **WHEN** a client connects with a valid certificate signed by the CA
- **THEN** the system SHALL accept the connection
- **AND** proceed with the request

#### Scenario: mTLS with invalid client certificate
- **WHEN** a client connects with an invalid or missing certificate
- **THEN** the system SHALL reject the connection
- **AND** return a TLS handshake error

#### Scenario: Missing CA file
- **WHEN** tls.ca_file is specified but the file does not exist
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate the CA file is missing

### Requirement: TLS configuration priority
The system SHALL use combined_file if specified, otherwise use separate cert_file and key_file.

#### Scenario: Combined file takes precedence
- **WHEN** both combined_file and separate cert_file/key_file are specified
- **THEN** the system SHALL use combined_file
- **AND** ignore the separate files

#### Scenario: Separate files used when no combined file
- **WHEN** only cert_file and key_file are specified
- **THEN** the system SHALL use the separate files

### Requirement: TLS with stdio transport
The system SHALL NOT enable TLS when using stdio transport.

#### Scenario: TLS configuration ignored with stdio
- **WHEN** transport is set to "stdio"
- **AND** TLS configuration is present
- **THEN** the system SHALL ignore the TLS configuration
- **AND** use stdio without encryption

### Requirement: TLS error handling
The system SHALL handle TLS errors gracefully and provide informative error messages.

#### Scenario: TLS handshake failure
- **WHEN** TLS handshake fails
- **THEN** the system SHALL log an error
- **AND** the error SHALL include details about the failure

#### Scenario: Certificate expired
- **WHEN** the server certificate has expired
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate the certificate is expired

#### Scenario: Certificate not yet valid
- **WHEN** the server certificate is not yet valid (future start date)
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate the certificate is not yet valid
