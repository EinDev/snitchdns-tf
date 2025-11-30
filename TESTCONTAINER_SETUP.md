# SnitchDNS Testcontainer Setup - Complete ✅

## Status: Fully Functional

All tests are passing! The testcontainer infrastructure is ready for TDD development of the Terraform provider.

## Test Results

```
=== RUN   TestSnitchDNSContainer
✅ Container started successfully
✅ API Key: JSvJuAR6xx... (auto-generated)
✅ HTTP Host: http://localhost:32790

=== RUN   TestSnitchDNSContainer/API_Authentication
✅ API authentication working correctly

=== RUN   TestSnitchDNSContainer/API_Zones_Endpoint
✅ Zones endpoint returning valid response
   Response: map[data:[] page:1 pages:0 per_page:50 total:0]

=== RUN   TestSnitchDNSContainer/API_Record_Types
✅ Record types endpoint working correctly
   Found 18 record types: [A AAAA AFSDB CAA CNAME DNAME HINFO MX NAPTR NS PTR RP SOA SPF SRV SSHFP TSIG TXT]

--- PASS: TestSnitchDNSContainer (100.82s)
    --- PASS: TestSnitchDNSContainer/API_Authentication (0.02s)
    --- PASS: TestSnitchDNSContainer/API_Zones_Endpoint (0.01s)
    --- PASS: TestSnitchDNSContainer/API_Record_Types (0.00s)
PASS
ok  	snitchdns-tf/internal/testcontainer	100.843s
```

## Architecture

### Container Setup
- **Base Image**: `debian:bookworm-slim` (lightweight, modern)
- **Python Environment**: Virtual environment with Flask-Migrate<3.0
- **Database**: SQLite (simple, no external dependencies)
- **Web Server**: Flask development server on port 80 (no Apache needed)
- **DNS Server**: SnitchDNS daemon on UDP port 2024
- **Default User**: `testadmin` / `password123` (auto-created)
- **API Key**: Auto-generated on startup, extracted by Go code

### Key Design Decisions

1. **Direct Flask Exposure**: Bypassed Apache proxy for simplicity and reliability
2. **Flask-Migrate<3.0**: Pinned to avoid compatibility issues with newer versions
3. **Debian over Ubuntu**: Smaller image, faster builds
4. **Port 80 for HTTP**: Standard port, easy to test
5. **Image Caching**: KeepImage=true for faster subsequent test runs

## File Structure

```
.
├── testcontainer/
│   ├── Dockerfile              # Container definition
│   └── entrypoint.sh          # Startup script, creates user & API key
├── internal/testcontainer/
│   ├── snitchdns.go           # Go wrapper for container management
│   └── snitchdns_test.go      # Integration tests
├── .github/workflows/
│   └── test.yml               # CI/CD pipeline
├── Makefile                   # Development commands
├── README.md                  # Project documentation
├── API_SPEC.md               # SnitchDNS API specification
└── TESTCONTAINER_SETUP.md    # This file
```

## Usage

### Run Tests Locally

```bash
# Run all tests
make test

# Run only integration tests
make test-integration

# Run with verbose output
go test -v ./internal/testcontainer/

# Run in short mode (skips integration tests)
go test -short ./...
```

### Use in Your Tests

```go
package main

import (
    "context"
    "testing"
    "snitchdns-tf/internal/testcontainer"
)

func TestMyFeature(t *testing.T) {
    ctx := context.Background()

    // Start container
    container, err := testcontainer.NewSnitchDNSContainer(ctx,
        testcontainer.SnitchDNSContainerRequest{
            ExposePorts: true,
        })
    if err != nil {
        t.Fatal(err)
    }
    defer container.Terminate(ctx)

    // Use the container
    apiEndpoint := container.GetAPIEndpoint()  // http://localhost:XXXXX/api/v1
    apiKey := container.APIKey                 // Auto-generated key

    // Make API calls to test your Terraform provider...
}
```

## Performance

- **Build Time**: ~80 seconds (first time), ~2 seconds (cached)
- **Startup Time**: ~7 seconds
- **Total Test Time**: ~100 seconds (includes build + startup + tests)

## CI/CD Integration

The GitHub Actions workflow automatically:
1. Sets up Go 1.24
2. Installs dependencies
3. Builds the testcontainer
4. Runs all tests
5. Uploads coverage reports

## Next Steps for TDD

Now you can start developing the Terraform provider using TDD:

1. **Write a failing test** for a Terraform resource (e.g., `resource_zone_test.go`)
2. **Run the test** - it should fail
3. **Implement the resource** (e.g., `resource_zone.go`)
4. **Run the test again** - it should pass
5. **Refactor** while keeping tests green

## Troubleshooting

### Container fails to start

Check Docker daemon:
```bash
docker ps
docker info
```

### Tests timeout

Increase timeout:
```bash
go test -timeout 10m ./internal/testcontainer/
```

### View container logs

```go
logs, err := container.Logs(ctx)
fmt.Println(logs)
```

### Manual container testing

Build and run manually:
```bash
cd testcontainer
docker build -t snitchdns-test .
docker run -p 8080:80 -p 2024:2024/udp snitchdns-test
```

Test API:
```bash
# Get the API key from logs
docker logs <container_id> | grep API_KEY

# Test endpoint
curl -H "X-SnitchDNS-Auth: YOUR_API_KEY" http://localhost:8080/api/v1/zones
```

## Features Verified

✅ Container builds successfully
✅ Database initializes
✅ Test user created
✅ API key generated
✅ Flask server starts
✅ DNS daemon starts
✅ API authentication works
✅ API endpoints return correct data
✅ Healthcheck passes
✅ Container cleanup works

## Known Limitations

- Single-threaded Flask dev server (not for production)
- SQLite database (lost on container restart)
- No HTTPS (HTTP only for testing)
- Test data not persisted

These limitations are acceptable for a test container and actually beneficial (clean state for each test).
