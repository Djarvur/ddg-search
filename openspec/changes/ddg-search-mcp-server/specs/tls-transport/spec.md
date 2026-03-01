# TLS Transport Specification

## Overview

This specification defines the TLS and HTTP transport capabilities of the MCP server, including certificate management and mTLS support.

## ADDED Requirements

### Requirement: TLS Support
The HTTP transport SHALL support TLS encryption.

#### Scenario: Enable TLS with certificate files
- **WHEN** tls.enabled is true and tls.cert_file and tls.key_file are specified
- **THEN** the server SHALL serve over HTTPS using the specified certificate

#### Scenario: Enable TLS with combined certificate
- **WHEN** tls.enabled is true and tls.combined is specified
- **THEN** the server SHALL serve over HTTPS using the combined certificate file

#### Scenario: TLS disabled by default
- **WHEN** tls.enabled is not set or false
- **THEN** the server SHALL serve over plain HTTP

### Requirement: mTLS Client Authentication
The server SHALL support mTLS with configurable client authentication modes.

#### Scenario: No client certificate required
- **WHEN** tls.client_auth is "none"
- **THEN** the server SHALL accept connections without requesting client certificates

#### Scenario: Request client certificate
- **WHEN** tls.client_auth is "request"
- **THEN** the server SHALL request but not require a client certificate

#### Scenario: Require client certificate
- **WHEN** tls.client_auth is "require"
- **THEN** the server SHALL reject connections without a valid client certificate

### Requirement: CA Certificate Support
The server SHALL support custom CA certificates for mTLS.

#### Scenario: CA certificate configured
- **WHEN** tls.ca_cert is specified
- **THEN** the server SHALL use that CA to verify client certificates

### Requirement: HTTP/2 Support
HTTP/2 SHALL be automatically enabled when TLS is enabled.

#### Scenario: HTTP/2 negotiation
- **WHEN** TLS is enabled and client supports HTTP/2
- **THEN** the connection SHALL use HTTP/2 protocol

### Requirement: Certificate Reload on SIGHUP
The server SHALL reload TLS certificates when receiving SIGHUP.

#### Scenario: SIGHUP reloads certificates
- **WHEN** SIGHUP is sent and certificate paths have changed
- **THEN** the server SHALL reload certificates from the new paths

#### Scenario: SIGHUP with unchanged certificates
- **WHEN** SIGHUP is sent but certificate paths are unchanged
- **THEN** the server SHALL continue using existing certificates
