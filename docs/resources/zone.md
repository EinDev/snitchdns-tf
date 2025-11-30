---
page_title: "snitchdns_zone Resource"
subcategory: ""
description: |-
  Manages a DNS zone in SnitchDNS.
---

# snitchdns_zone

Manages a DNS zone in SnitchDNS. Zones are containers for DNS records and can be configured with various options like catch-all, forwarding, and regex matching.

## Example Usage

### Basic Zone

```terraform
resource "snitchdns_zone" "example" {
  domain     = "example.com"
  active     = true
  catch_all  = false
  forwarding = false
  regex      = false
}
```

### Zone with Tags

```terraform
resource "snitchdns_zone" "production" {
  domain     = "prod.example.com"
  active     = true
  catch_all  = false
  forwarding = false
  regex      = false
  tags       = ["production", "web", "critical"]
}
```

### Catch-All Zone

```terraform
resource "snitchdns_zone" "catchall" {
  domain     = "wildcard.example.com"
  active     = true
  catch_all  = true  # Responds to any subdomain query
  forwarding = false
  regex      = false
  tags       = ["wildcard"]
}
```

### Forwarding Zone

```terraform
resource "snitchdns_zone" "forwarding" {
  domain     = "forward.example.com"
  active     = true
  catch_all  = false
  forwarding = true  # Forwards unmatched queries to upstream
  regex      = false
}
```

### Regex Zone

```terraform
resource "snitchdns_zone" "regex" {
  domain     = ".*\\.test\\.example\\.com"  # Matches *.test.example.com
  active     = true
  catch_all  = false
  forwarding = false
  regex      = true  # Enables regex pattern matching
}
```

### Disabled Zone

```terraform
resource "snitchdns_zone" "disabled" {
  domain     = "inactive.example.com"
  active     = false  # Zone exists but doesn't respond to queries
  catch_all  = false
  forwarding = false
  regex      = false
}
```

## Schema

### Required

- `domain` (String) - The domain name for this zone (e.g., `example.com`). Must be between 1 and 255 characters. When `regex` is enabled, this can be a regular expression pattern.

- `active` (Boolean) - Whether the zone is active and will respond to DNS queries. Set to `false` to disable the zone without deleting it.

- `catch_all` (Boolean) - Enable catch-all DNS queries for this zone. When enabled, the zone will respond to queries for any subdomain, even if no specific record exists.

- `forwarding` (Boolean) - Enable DNS forwarding to upstream DNS servers. When enabled, unmatched queries will be forwarded to a configured upstream resolver.

- `regex` (Boolean) - Use regular expression matching for the domain name. When enabled, the domain field can contain a regex pattern instead of a literal domain.

### Optional

- `tags` (List of String) - List of tags to organize and categorize zones. Tags can be used for filtering and grouping zones in the SnitchDNS UI.

### Read-Only

- `id` (String) - Unique identifier for the zone. Assigned by the API upon creation.

- `user_id` (Number) - ID of the user who owns this zone. Automatically set by the API based on authentication.

- `master` (Boolean) - Indicates if this is a master zone. Master zones have special privileges and cannot be modified via the API.

- `created_at` (String) - Timestamp when the zone was created in RFC3339 format.

- `updated_at` (String) - Timestamp when the zone was last updated in RFC3339 format.

## Import

Zones can be imported using their ID:

```bash
terraform import snitchdns_zone.example 123
```

To find the zone ID, you can:
1. Check the SnitchDNS web UI
2. Use the SnitchDNS API to list zones
3. Check the Terraform state of an existing zone

## Notes

- **Regex Zones**: When using regex patterns, ensure the pattern is properly escaped for Terraform strings. Use double backslashes (`\\`) for regex escape sequences.

- **Catch-All Behavior**: Catch-all zones will respond to any subdomain query, even if no specific record exists. This can be useful for capturing DNS exfiltration attempts or providing wildcard functionality.

- **Forwarding**: The forwarding feature requires upstream DNS servers to be configured in the SnitchDNS server settings.

- **Master Zones**: Master zones are created automatically by SnitchDNS and cannot be modified or deleted through the API. The `master` attribute is read-only.

- **External Deletion**: If a zone is deleted outside of Terraform (e.g., through the SnitchDNS web UI), Terraform will automatically detect this during the next `terraform plan` or `terraform apply` and remove it from the state.

- **Tags**: Tags are purely organizational and do not affect DNS functionality. They are useful for managing large numbers of zones.

## Common Patterns

### Managing Related Zones

```terraform
locals {
  environments = ["dev", "staging", "prod"]
}

resource "snitchdns_zone" "env_zones" {
  for_each = toset(local.environments)

  domain     = "${each.key}.example.com"
  active     = true
  catch_all  = false
  forwarding = false
  regex      = false
  tags       = [each.key, "auto-managed"]
}
```

### Zone with Records

```terraform
resource "snitchdns_zone" "example" {
  domain     = "example.com"
  active     = true
  catch_all  = false
  forwarding = false
  regex      = false
  tags       = ["production"]
}

resource "snitchdns_record" "www" {
  zone_id = snitchdns_zone.example.id
  name    = "www"
  type    = "A"
  cls     = "IN"
  ttl     = 3600
  data    = "192.168.1.100"
}

resource "snitchdns_record" "mx" {
  zone_id = snitchdns_zone.example.id
  name    = "@"
  type    = "MX"
  cls     = "IN"
  ttl     = 3600
  data    = "10 mail.example.com."
}
```
