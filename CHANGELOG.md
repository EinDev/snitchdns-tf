# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of the SnitchDNS Terraform Provider
- Zone resource (`snitchdns_zone`) for managing DNS zones
  - Support for domain, active, catch_all, forwarding, and regex options
  - Tag support for zone organization
  - Import functionality
- Record resource (`snitchdns_record`) for managing DNS records
  - Support for all standard DNS record types (A, AAAA, CNAME, MX, TXT, NS, SRV, CAA, etc.)
  - Conditional response support for canary deployments
  - Import functionality
- Provider configuration via HCL or environment variables
  - `SNITCHDNS_API_URL` environment variable support
  - `SNITCHDNS_API_KEY` environment variable support
- Comprehensive schema validation
  - Domain length validation (1-255 characters)
  - TTL range validation (1 to 2,147,483,647)
  - DNS class validation (IN, CH, HS)
  - DNS record type validation (18 supported types)
- Defensive read operations
  - Automatic detection of externally deleted resources
  - Clean removal from state with warning logs
- Plan modifiers for immutable fields
  - `zone_id` and `type` in records require replacement when changed
- Structured logging with tflog
  - Debug logs for operations
  - Warning logs for state changes
- Full acceptance test coverage
  - Zone CRUD operations
  - Record CRUD operations
  - Import state tests
  - Tag management tests
  - Conditional record tests
- Comprehensive documentation
  - Provider configuration guide
  - Zone resource documentation with examples
  - Record resource documentation with examples
  - Basic, complete, and advanced usage examples
- Security features
  - API key marked as sensitive
  - No sensitive data in logs

### Changed
N/A - Initial release

### Deprecated
N/A - Initial release

### Removed
N/A - Initial release

### Fixed
N/A - Initial release

### Security
- API keys are marked as sensitive and not exposed in logs
- Input validation prevents injection attacks
- Proper error handling without leaking sensitive information

## [1.0.0] - TBD

Initial release - see Unreleased section above for details.

---

## Release Notes

### Version 1.0.0

This is the initial release of the SnitchDNS Terraform Provider, providing full support for managing SnitchDNS zones and records through Terraform.

**Features:**
- Complete zone management with advanced options (catch-all, forwarding, regex)
- Full DNS record support for all standard record types
- Conditional responses for advanced deployment strategies
- Environment variable configuration support
- Comprehensive validation and error handling
- Import support for existing resources
- Production-ready with extensive testing

**Requirements:**
- Terraform >= 1.0
- Go >= 1.21 (for development)
- SnitchDNS server with API access

**Getting Started:**
1. Install the provider from the Terraform Registry
2. Configure with your SnitchDNS API URL and key
3. Start managing your DNS infrastructure as code

See the [documentation](docs/) and [examples](examples/) for detailed usage information.
