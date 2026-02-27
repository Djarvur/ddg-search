# TLS/mTLS Testing Guide

This guide explains how to create self-signed certificates for testing TLS and mTLS functionality locally.

## Prerequisites

- OpenSSL installed (available on macOS by default)
- Basic understanding of TLS certificates

## Part 1: Create a Self-Signed CA for TLS Testing

### Step 1: Create a directory structure

```bash
mkdir -p ~/local-tls/ca
mkdir -p ~/local-tls/server
mkdir -p ~/local-tls/client
cd ~/local-tls
```

### Step 2: Generate CA private key

```bash
openssl genrsa -out ca/ca-key.pem 4096
```

### Step 3: Generate CA certificate

```bash
openssl req -new -x509 -days 365 -key ca/ca-key.pem -sha256 -out ca/ca-cert.pem \
  -subj "/C=US/ST=State/L=City/O=LocalTest/CN=LocalTest-CA"
```

### Step 4: Add CA to macOS Keychain (Trusted)

```bash
# Add to system keychain (requires sudo)
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ca/ca-cert.pem

# Or add to user keychain (no sudo required)
security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain-db ca/ca-cert.pem
```

**Verification:**

```bash
# List trusted certificates
security find-certificate -c "LocalTest-CA" -p | openssl x509 -noout -text
```

**To remove the CA later:**

```bash
# Remove from system keychain
sudo security delete-certificate -c "LocalTest-CA" /Library/Keychains/System.keychain

# Or remove from user keychain
security delete-certificate -c "LocalTest-CA" ~/Library/Keychains/login.keychain-db
```

## Part 2: Generate Server Certificate Signed by CA

### Step 1: Generate server private key

```bash
openssl genrsa -out server/server-key.pem 4096
```

### Step 2: Create server certificate signing request (CSR)

```bash
openssl req -new -key server/server-key.pem -out server/server.csr \
  -subj "/C=US/ST=State/L=City/O=LocalTest/CN=localhost"
```

**Note:** For local testing, use `localhost` as the Common Name (CN). For production, use your actual domain.

### Step 3: Create server certificate configuration

Create a file `server/server-cert.conf`:

```ini
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = US
ST = State
L = City
O = LocalTest
CN = localhost

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = *.localhost
IP.1 = 127.0.0.1
IP.2 = ::1
```

### Step 4: Generate server certificate signed by CA

```bash
openssl x509 -req -in server/server.csr -CA ca/ca-cert.pem -CAkey ca/ca-key.pem \
  -CAcreateserial -out server/server-cert.pem -days 365 -sha256 \
  -extfile server/server-cert.conf -extensions v3_req
```

### Step 5: Verify server certificate

```bash
openssl x509 -in server/server-cert.pem -text -noout
```

## Part 3: Create a Separate CA for mTLS Testing

For mTLS testing, it's best practice to use a separate CA for client certificates.

### Step 1: Generate mTLS CA private key

```bash
openssl genrsa -out ca/mtls-ca-key.pem 4096
```

### Step 2: Generate mTLS CA certificate

```bash
openssl req -new -x509 -days 365 -key ca/mtls-ca-key.pem -sha256 -out ca/mtls-ca-cert.pem \
  -subj "/C=US/ST=State/L=City/O=LocalTest/CN=LocalTest-mTLS-CA"
```

### Step 3: Add mTLS CA to macOS Keychain (Trusted)

```bash
# Add to user keychain
security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain-db ca/mtls-ca-cert.pem
```

## Part 4: Generate mTLS Client Certificates

### Step 1: Generate client private key

```bash
openssl genrsa -out client/client-key.pem 4096
```

### Step 2: Create client certificate signing request (CSR)

```bash
openssl req -new -key client/client-key.pem -out client/client.csr \
  -subj "/C=US/ST=State/L=City/O=LocalTest/CN=test-client"
```

### Step 3: Create client certificate configuration

Create a file `client/client-cert.conf`:

```ini
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = US
ST = State
L = City
O = LocalTest
CN = test-client

[v3_req]
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth
```

### Step 4: Generate client certificate signed by mTLS CA

```bash
openssl x509 -req -in client/client.csr -CA ca/mtls-ca-cert.pem -CAkey ca/mtls-ca-key.pem \
  -CAcreateserial -out client/client-cert.pem -days 365 -sha256 \
  -extfile client/client-cert.conf -extensions v3_req
```

### Step 5: Verify client certificate

```bash
openssl x509 -in client/client-cert.pem -text -noout
```

## Part 5: Test TLS Configuration

### Configure your-server with TLS

Create or update `~/.config/local/config.yaml`:

```yaml
server:
  protocol: http
  bind_address: localhost:9100
  tls:
    enabled: true
    cert_file: ~/local-tls/server/server-cert.pem
    key_file: ~/local-tls/server/server-key.pem
    min_version: "1.2"
logging:
  level: debug
```

### Test TLS connection

```bash
# Start the server
go run ./cmd/your-server

# In another terminal, test with curl
curl -v https://localhost:9100/health

# Test with curl using the CA cert
curl -v --cacert ~/local-tls/ca/ca-cert.pem https://localhost:9100/health

# Test SSE endpoint
curl -v -N --cacert ~/local-tls/ca/ca-cert.pem https://localhost:9100/sse
```

## Part 6: Test mTLS Configuration

### Configure your-server with mTLS

Update `~/.config/local/config.yaml`:

```yaml
server:
  protocol: http
  bind_address: localhost:9100
  tls:
    enabled: true
    cert_file: ~/local-tls/server/server-cert.pem
    key_file: ~/local-tls/server/server-key.pem
    min_version: "1.2"
    mtls:
      enabled: true
      ca_file: ~/local-tls/ca/mtls-ca-cert.pem
logging:
  level: debug
```

### Test mTLS connection

```bash
# Start the server
go run ./cmd/your-server

# Test without client certificate (should fail)
curl -v --cacert ~/local-tls/ca/ca-cert.pem https://localhost:9100/health

# Test with client certificate (should succeed)
curl -v --cacert ~/local-tls/ca/ca-cert.pem \
  --cert ~/local-tls/client/client-cert.pem \
  --key ~/local-tls/client/client-key.pem \
  https://localhost:9100/health

# Test SSE endpoint with mTLS
curl -v -N --cacert ~/local-tls/ca/ca-cert.pem \
  --cert ~/local-tls/client/client-cert.pem \
  --key ~/local-tls/client/client-key.pem \
  https://localhost:9100/sse
```

## Part 7: Generate Multiple Client Certificates

For testing different clients, you can generate multiple client certificates:

```bash
# Generate client 2
openssl genrsa -out client/client2-key.pem 4096
openssl req -new -key client/client2-key.pem -out client/client2.csr \
  -subj "/C=US/ST=State/L=City/O=LocalTest/CN=test-client-2"
openssl x509 -req -in client/client2.csr -CA ca/mtls-ca-cert.pem -CAkey ca/mtls-ca-key.pem \
  -CAcreateserial -out client/client2-cert.pem -days 365 -sha256 \
  -extfile client/client-cert.conf -extensions v3_req

# Generate client 3
openssl genrsa -out client/client3-key.pem 4096
openssl req -new -key client/client3-key.pem -out client/client3.csr \
  -subj "/C=US/ST=State/L=City/O=LocalTest/CN=test-client-3"
openssl x509 -req -in client/client3.csr -CA ca/mtls-ca-cert.pem -CAkey ca/mtls-ca-key.pem \
  -CAcreateserial -out client/client3-cert.pem -days 365 -sha256 \
  -extfile client/client-cert.conf -extensions v3_req
```

## Part 8: Certificate Expiration and Renewal

### Check certificate expiration

```bash
# Check server certificate
openssl x509 -in server/server-cert.pem -noout -dates

# Check client certificate
openssl x509 -in client/client-cert.pem -noout -dates

# Check CA certificate
openssl x509 -in ca/ca-cert.pem -noout -dates
```

### Renew certificates

To renew certificates, simply regenerate them following the same steps above. The CA can sign new certificates even after the original certificates expire.

## Part 9: Troubleshooting

### Common Issues

**Issue: "certificate verify failed"**

```bash
# Verify the certificate chain
openssl verify -CAfile ca/ca-cert.pem server/server-cert.pem

# Check if CA is in keychain
security find-certificate -c "LocalTest-CA"
```

**Issue: "hostname doesn't match"**

Ensure the Common Name (CN) in the certificate matches the hostname you're connecting to. For local testing, use `localhost`.

**Issue: "connection refused"**

Check that the server is running and listening on the correct port:

```bash
lsof -i :9100
```

**Issue: mTLS connection rejected**

Verify that:
1. The client certificate is signed by the correct CA
2. The CA is trusted by the server
3. The client certificate is not expired

```bash
# Verify client certificate chain
openssl verify -CAfile ca/mtls-ca-cert.pem client/client-cert.pem
```

## Part 10: Cleanup

### Remove certificates from macOS Keychain

```bash
# Remove TLS CA
security delete-certificate -c "LocalTest-CA" ~/Library/Keychains/login.keychain-db

# Remove mTLS CA
security delete-certificate -c "LocalTest-mTLS-CA" ~/Library/Keychains/login.keychain-db
```

### Delete certificate files

```bash
rm -rf ~/local-tls
```

## Quick Reference

### File Locations

```
~/local-tls/
├── ca/
│   ├── ca-key.pem              # TLS CA private key
│   ├── ca-cert.pem             # TLS CA certificate
│   ├── mtls-ca-key.pem         # mTLS CA private key
│   └── mtls-ca-cert.pem        # mTLS CA certificate
├── server/
│   ├── server-key.pem          # Server private key
│   ├── server.csr              # Server CSR
│   ├── server-cert.pem         # Server certificate
│   └── server-cert.conf        # Server certificate config
└── client/
    ├── client-key.pem          # Client private key
    ├── client.csr              # Client CSR
    ├── client-cert.pem         # Client certificate
    └── client-cert.conf        # Client certificate config
```

### Configuration Paths

```yaml
# TLS configuration
server:
  tls:
    cert_file: ~/local-tls/server/server-cert.pem
    key_file: ~/local-tls/server/server-key.pem
    mtls:
      ca_file: ~/local-tls/ca/mtls-ca-cert.pem
```

### Test Commands

```bash
# Test TLS
curl -v --cacert ~/local-tls/ca/ca-cert.pem https://localhost:9100/health

# Test mTLS
curl -v --cacert ~/local-tls/ca/ca-cert.pem \
  --cert ~/local-tls/client/client-cert.pem \
  --key ~/local-tls/client/client-key.pem \
  https://localhost:9100/health
```

## Security Notes

⚠️ **Important:** These certificates are for testing purposes only. Do not use them in production.

- The CA private keys are not protected with passphrases
- Certificates have a 365-day validity period
- No certificate revocation mechanism is implemented
- Use proper certificate management for production deployments

## Additional Resources

- [OpenSSL Documentation](https://www.openssl.org/docs/)
- [macOS Security Command Reference](https://ss64.com/osx/security.html)
- [TLS Configuration Best Practices](https://wiki.mozilla.org/Security/Server_Side_TLS)
