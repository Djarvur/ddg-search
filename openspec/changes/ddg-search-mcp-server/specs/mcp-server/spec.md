# MCP Server Specification

## Overview

This specification defines the MCP (Model Context Protocol) server implementation that exposes web search and page fetching capabilities to Claude Code and other MCP clients.

## ADDED Requirements

### Requirement: MCP Server Supports Stdio Transport
The MCP server SHALL support stdio transport as the default transport mode for local development and integration with Claude Code.

#### Scenario: Server starts in stdio mode
- **WHEN** the server is started without specifying transport or with `--transport stdio`
- **THEN** it shall listen for MCP requests on standard input and respond on standard output

#### Scenario: Server handles stdio initialization
- **WHEN** a client sends an initialize request via stdio
- **THEN** the server SHALL respond with server name, version, and capabilities

### Requirement: MCP Server Supports HTTP Transport
The MCP server SHALL support HTTP transport with StreamableHTTP for remote connections.

#### Scenario: Server starts in HTTP mode
- **WHEN** the server is started with `--transport http`
- **THEN** it shall listen for HTTP requests on the configured host and port

#### Scenario: Server binds to configured address
- **WHEN** HTTP transport is enabled with `--host localhost --port 9100`
- **THEN** the server SHALL listen on localhost:9100

### Requirement: Server Provides Tool List
The MCP server SHALL expose available tools to connected clients.

#### Scenario: Client requests tool list
- **WHEN** a client sends a tools/list request
- **THEN** the server SHALL return a list containing search and fetch tools with their definitions

#### Scenario: Tool definitions include input schemas
- **WHEN** the server responds to tools/list
- **THEN** each tool SHALL include a JSON schema defining its input parameters

### Requirement: Server Handles Tool Invocation
The MCP server SHALL execute tools when called by clients.

#### Scenario: Client calls search tool
- **WHEN** a client sends a tools/call request for search with valid query
- **THEN** the server SHALL execute the search and return results in the specified format

#### Scenario: Client calls fetch tool
- **WHEN** a client sends a tools/call request for fetch with valid URL
- **THEN** the server SHALL fetch the page and return content in the specified format

### Requirement: Server Returns Proper Error Responses
The server SHALL return MCP-compliant error responses for invalid requests.

#### Scenario: Client sends malformed request
- **WHEN** a client sends a request with invalid JSON
- **THEN** the server SHALL return a JSON-RPC error response with code -32600

#### Scenario: Client requests unknown tool
- **WHEN** a client calls a tool that does not exist
- **THEN** the server SHALL return error code -32601 (method not found)
