# Testing Guide

This document explains how to test the cert-manager-webhook-infoblox-wapi webhook.

## Prerequisites

- Go 1.23 or later
- `make` installed
- (Optional) Access to an Infoblox GRID for integration testing

## Quick Start

Run unit tests (no Infoblox required):

```bash
make test
```

This will:
1. Download and install `setup-envtest` tool
2. Download Kubernetes test binaries (etcd, kube-apiserver, kubectl)
3. Run the test suite

## Test Types

### Unit Tests (No Infoblox Required)

The default `make test` command runs tests that validate the webhook structure without requiring actual Infoblox credentials. Tests are skipped gracefully if `TEST_ZONE_NAME` is not set.

```bash
make test
```

### Integration Tests (Requires Infoblox)

To run full integration tests against a real Infoblox GRID:

```bash
export TEST_ZONE_NAME="example.com."  # Must end with a dot
make test
```

**Note:** Integration tests require:
- A configured Infoblox GRID accessible from your test environment
- Valid credentials in `testdata/infoblox-wapi/credentials.yaml`
- Proper configuration in `testdata/infoblox-wapi/config.json`

### Coverage Tests

Generate a coverage report:

```bash
make test-coverage
```

This creates `coverage.out` which can be viewed with:

```bash
go tool cover -html=coverage.out
```

### Race Detection Tests

Run tests with the Go race detector:

```bash
make test-race
```

## Test Configuration

### Test Data Structure

The `testdata/infoblox-wapi/` directory contains test configuration:

```
testdata/infoblox-wapi/
├── config.json           # Webhook configuration (required)
├── config.json.sample    # Template
├── credentials.yaml      # Infoblox credentials (gitignored)
└── credentials.yaml.sample # Template
```

### Setting Up Integration Tests

1. **Copy sample files:**
   ```bash
   cp testdata/infoblox-wapi/config.json.sample testdata/infoblox-wapi/config.json
   cp testdata/infoblox-wapi/credentials.yaml.sample testdata/infoblox-wapi/credentials.yaml
   ```

2. **Configure Infoblox connection** in `config.json`:
   ```json
   {
     "host": "infoblox.example.com",
     "port": 443,
     "version": "2.5",
     "view": "default",
     "sslVerify": true
   }
   ```

3. **Add credentials** in `credentials.yaml`:
   ```yaml
   username: admin
   password: your-password
   ```

4. **Run tests:**
   ```bash
   TEST_ZONE_NAME="your-zone.example.com." make test
   ```

## Continuous Integration

The CI workflow (`.github/workflows/ci.yml`) automatically:
- Runs unit tests on every PR and push
- Skips integration tests (no credentials available)
- Generates coverage reports
- Runs race detector
- Uploads coverage to Codecov

### Local CI Simulation

To simulate CI locally:

```bash
# Run all checks
make test
make test-race
make test-coverage
go vet ./...
gofmt -l .
```

Or use the linter (after installing golangci-lint):

```bash
golangci-lint run
```

## Test Environment

The test suite uses [envtest](https://book.kubebuilder.io/reference/envtest.html) which provides:
- A real Kubernetes API server
- etcd for storage
- No kubelet (pods won't run)
- No controller manager
- Perfect for testing CRDs and webhooks

### Envtest Versions

By default, tests use Kubernetes 1.31.x. Override with:

```bash
make test ENVTEST_K8S_VERSION=1.30.x
```

Available versions are listed at:
https://raw.githubusercontent.com/kubernetes-sigs/controller-tools/HEAD/envtest-releases.yaml

## Troubleshooting

### Tests fail with "etcd not found"

The Makefile should automatically download test assets. If it fails:

```bash
# Clean and retry
make clean
make test
```

### Tests hang

Check if etcd or kube-apiserver are already running:

```bash
pkill -9 etcd kube-apiserver
make test
```

### Integration tests fail

1. Verify Infoblox connectivity:
   ```bash
   curl -k https://your-infoblox-host/wapi/v2.5/
   ```

2. Check credentials are correct in `testdata/infoblox-wapi/credentials.yaml`

3. Verify DNS zone exists in Infoblox

4. Ensure TEST_ZONE_NAME ends with a dot: `example.com.`

### Permission denied errors

The test binaries need execute permissions:

```bash
chmod +x bin/setup-envtest
```

## Writing Tests

When adding new tests, follow the cert-manager webhook test patterns:

```go
func TestMyFeature(t *testing.T) {
    // Arrange
    solver := &customDNSProviderSolver{}
    
    // Act
    result := solver.MyMethod()
    
    // Assert
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

For DNS challenge tests, use the `acmetest.NewFixture` helper (see `main_test.go`).

## Related Documentation

- [cert-manager Webhook Testing](https://cert-manager.io/docs/contributing/dns-providers/)
- [Kubebuilder Envtest](https://book.kubebuilder.io/reference/envtest.html)
- [cert-manager Webhook Example](https://github.com/cert-manager/webhook-example)
