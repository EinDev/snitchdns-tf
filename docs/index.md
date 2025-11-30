---
page_title: "SnitchDNS Provider"
subcategory: ""
description: |-
  Terraform provider for managing SnitchDNS resources including DNS zones and records.
---

# SnitchDNS Provider

The SnitchDNS provider allows you to manage SnitchDNS resources using Terraform. SnitchDNS is a database-driven DNS server with a web UI for managing zones and records.

## Features

- Manage DNS zones with tags and configuration
- Manage DNS records with full support for all standard DNS record types
- Import existing resources from SnitchDNS
- Automatic detection and removal of externally deleted resources
- Environment variable support for credentials

## Example Usage

```terraform
terraform {
  required_providers {
    snitchdns = {
      source = "EinDev/snitchdns"
      version = "~> 1.0"
    }
  }
}

provider "snitchdns" {
  api_url = "http://localhost:8000"
  api_key = "your-api-key-here"
}

# Create a DNS zone
resource "snitchdns_zone" "example" {
  domain = "example.com"
  tags   = ["production", "web"]
}

# Create a DNS record
resource "snitchdns_record" "www" {
  zone_id = snitchdns_zone.example.id
  name    = "www"
  type    = "A"
  cls     = "IN"
  ttl     = 3600
  data    = "192.168.1.100"
}
```

## Configuration

The provider can be configured using either the provider block or environment variables.

### Provider Block

```terraform
provider "snitchdns" {
  api_url = "http://localhost:8000"
  api_key = "your-api-key-here"
}
```

### Environment Variables

```bash
export SNITCHDNS_API_URL="http://localhost:8000"
export SNITCHDNS_API_KEY="your-api-key-here"
```

When using environment variables, you can omit the provider configuration:

```terraform
provider "snitchdns" {}
```

## Schema

### Required

Note: At least one of the following must be provided, either directly or via environment variables.

- `api_url` (String) - SnitchDNS API URL. Can also be set via `SNITCHDNS_API_URL` environment variable.
  - Example: `http://localhost:8000` or `https://dns.example.com`

- `api_key` (String, Sensitive) - SnitchDNS API Key for authentication. Can also be set via `SNITCHDNS_API_KEY` environment variable.
  - Obtain this from your SnitchDNS web UI under Settings > API

## Authentication

To obtain an API key:

1. Log in to your SnitchDNS web interface
2. Navigate to Settings > API
3. Generate a new API key
4. Copy the key and use it in your provider configuration

**Security Note:** The API key is marked as sensitive and will not appear in Terraform logs or output. Consider using environment variables or secret management tools instead of hardcoding keys in your Terraform files.

## Getting Started

1. Install and configure SnitchDNS server
2. Generate an API key from the SnitchDNS web UI
3. Configure the Terraform provider with your API URL and key
4. Start managing your DNS infrastructure as code

## Resources

- [snitchdns_zone](resources/zone.md) - Manage DNS zones
- [snitchdns_record](resources/record.md) - Manage DNS records

## Support

For issues or questions:
- Provider issues: [GitHub Issues](https://github.com/EinDev/snitchdns-tf/issues)
- SnitchDNS documentation: [SnitchDNS Docs](https://github.com/ctxis/SnitchDNS)
