---
page_title: "snitchdns_record Resource"
subcategory: ""
description: |-
  Manages a DNS record within a SnitchDNS zone.
---

# snitchdns_record

Manages a DNS record within a SnitchDNS zone. Records define the actual DNS responses for queries and support all standard DNS record types (A, AAAA, CNAME, MX, TXT, etc.) as well as conditional responses.

## Example Usage

### A Record (IPv4 Address)

```terraform
resource "snitchdns_record" "web" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 3600

  data = {
    address = "192.168.1.100"
  }
}
```

### AAAA Record (IPv6 Address)

```terraform
resource "snitchdns_record" "web_ipv6" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "AAAA"
  ttl     = 3600

  data = {
    address = "2001:0db8::1"
  }
}
```

### CNAME Record

```terraform
resource "snitchdns_record" "www" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "CNAME"
  ttl     = 3600

  data = {
    name = "web.example.com."
  }
}
```

### MX Record (Mail Exchange)

```terraform
resource "snitchdns_record" "mail" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "MX"
  ttl     = 3600

  data = {
    priority = "10"
    hostname = "mail.example.com."
  }
}
```

### TXT Record

```terraform
resource "snitchdns_record" "spf" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "TXT"
  ttl     = 3600

  data = {
    data = "v=spf1 include:_spf.example.com ~all"
  }
}
```

### NS Record (Name Server)

```terraform
resource "snitchdns_record" "nameserver" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "NS"
  ttl     = 86400

  data = {
    name = "ns1.example.com."
  }
}
```

### SRV Record (Service)

```terraform
resource "snitchdns_record" "sip" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "SRV"
  ttl     = 3600

  data = {
    priority = "10"
    weight   = "60"
    port     = "5060"
    target   = "sip.example.com."
  }
}
```

### CAA Record (Certificate Authority Authorization)

```terraform
resource "snitchdns_record" "caa" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "CAA"
  ttl     = 3600

  data = {
    flags = "0"
    tag   = "issue"
    value = "letsencrypt.org"
  }
}
```

### Conditional Record

```terraform
resource "snitchdns_record" "conditional" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 60

  # Return this IP for the first 10 queries
  data = {
    address = "192.168.1.100"
  }

  # Conditional response after 10 queries
  is_conditional   = true
  conditional_limit = 10
  conditional_reset = true

  # Return this IP after 10 queries
  conditional_data = {
    address = "192.168.1.200"
  }
}
```

### Complete Web Infrastructure Example

```terraform
resource "snitchdns_zone" "example" {
  domain     = "example.com"
  active     = true
  catch_all  = false
  forwarding = false
  regex      = false
  tags       = ["production"]
}

# Root A record
resource "snitchdns_record" "root" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 3600
  data    = { address = "192.168.1.100" }
}

# WWW CNAME
resource "snitchdns_record" "www" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "CNAME"
  ttl     = 3600
  data    = { name = "example.com." }
}

# Mail MX record
resource "snitchdns_record" "mx" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "MX"
  ttl     = 3600
  data = {
    priority = "10"
    hostname = "mail.example.com."
  }
}

# Mail server A record
resource "snitchdns_record" "mail" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 3600
  data    = { address = "192.168.1.101" }
}

# SPF TXT record
resource "snitchdns_record" "spf" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "TXT"
  ttl     = 3600
  data    = { data = "v=spf1 mx -all" }
}
```

## Schema

### Required

- `zone_id` (String) - ID of the zone this record belongs to. Records are always associated with a specific zone. **Note:** Changing this requires resource replacement.

- `active` (Boolean) - Whether the record is active and will respond to DNS queries. Set to `false` to temporarily disable without deleting.

- `cls` (String) - DNS class for the record. Must be one of: `IN` (Internet), `CH` (Chaos), or `HS` (Hesiod). In most cases, use `IN`.

- `type` (String) - DNS record type. Supported types: `A`, `AAAA`, `AFSDB`, `CAA`, `CNAME`, `DNAME`, `HINFO`, `MX`, `NAPTR`, `NS`, `PTR`, `RP`, `SOA`, `SPF`, `SRV`, `SSHFP`, `TSIG`, `TXT`. **Note:** Changing this requires resource replacement.

- `ttl` (Number) - Time to live in seconds (1 to 2,147,483,647). Determines how long DNS resolvers should cache this record. Common values:
  - 60: 1 minute (dynamic/testing)
  - 300: 5 minutes (frequently changing)
  - 3600: 1 hour (standard)
  - 86400: 1 day (stable)

- `data` (Map of String) - Record-specific data as key-value pairs. The required fields depend on the record type. See [Data Field Formats](#data-field-formats) below.

### Optional

- `is_conditional` (Boolean) - Enable conditional responses based on query count. When enabled, the record can return different data based on how many times it has been queried.

- `conditional_limit` (Number) - Query limit for conditional responses. When `conditional_count` reaches this limit, the `conditional_data` is returned instead.

- `conditional_reset` (Boolean) - Reset the query counter when the limit is reached. If `true`, the counter resets to 0; if `false`, it remains at the limit.

- `conditional_data` (Map of String) - Alternative data to return when conditional limit is reached. Uses the same format as the `data` attribute.

### Read-Only

- `id` (String) - Unique identifier for the DNS record. Assigned by the API upon creation.

- `conditional_count` (Number) - Current query count for conditional logic. Automatically incremented by SnitchDNS when the record is queried.

## Data Field Formats

The `data` attribute format varies by record type. Here are the required fields for each type:

### A Record
```terraform
data = {
  address = "192.168.1.100"
}
```

### AAAA Record
```terraform
data = {
  address = "2001:0db8::1"
}
```

### CNAME Record
```terraform
data = {
  name = "target.example.com."  # Must end with a dot
}
```

### MX Record
```terraform
data = {
  priority = "10"
  hostname = "mail.example.com."  # Must end with a dot
}
```

### TXT Record
```terraform
data = {
  data = "your text content here"
}
```

### NS Record
```terraform
data = {
  name = "ns1.example.com."  # Must end with a dot
}
```

### SRV Record
```terraform
data = {
  priority = "10"
  weight   = "60"
  port     = "5060"
  target   = "service.example.com."  # Must end with a dot
}
```

### CAA Record
```terraform
data = {
  flags = "0"      # Usually 0 or 128
  tag   = "issue"  # issue, issuewild, or iodef
  value = "letsencrypt.org"
}
```

### SOA Record
```terraform
data = {
  mname   = "ns1.example.com."
  rname   = "admin.example.com."
  serial  = "2024010101"
  refresh = "3600"
  retry   = "600"
  expire  = "86400"
  minimum = "3600"
}
```

### PTR Record
```terraform
data = {
  name = "host.example.com."  # Must end with a dot
}
```

### SPF Record
```terraform
data = {
  data = "v=spf1 mx include:_spf.example.com ~all"
}
```

## Import

Records can be imported using the format `zone_id:record_id`:

```bash
terraform import snitchdns_record.example 123:456
```

Where:
- `123` is the zone ID
- `456` is the record ID

To find these IDs:
1. Check the SnitchDNS web UI
2. Use the SnitchDNS API to list zones and records
3. Check the Terraform state of existing resources

## Notes

### General

- **External Deletion**: If a record is deleted outside of Terraform (e.g., through the SnitchDNS web UI), Terraform will automatically detect this during the next `terraform plan` or `terraform apply` and remove it from the state.

- **Zone Dependency**: Records must belong to a zone. If the zone is destroyed, all associated records will be deleted by SnitchDNS.

- **Immutable Fields**: The `zone_id` and `type` fields cannot be changed after creation. Modifying them will destroy and recreate the record.

### DNS Best Practices

- **Trailing Dots**: For hostname/domain fields in data (CNAME, MX, NS, etc.), always include the trailing dot (`.`) to indicate a fully qualified domain name (FQDN). Without the dot, the zone domain will be appended.

- **TTL Values**:
  - Use lower TTL (60-300) during testing or before DNS changes
  - Use higher TTL (3600-86400) for stable records to reduce DNS load
  - Consider the impact on propagation time when planning changes

- **CNAME Limitations**: CNAME records cannot coexist with other record types for the same name. Don't create a CNAME for `@` (zone apex) if you have other records there.

### Conditional Records

Conditional records are useful for:
- DNS-based canary deployments
- Rotating IP addresses after N queries
- Tracking DNS query patterns
- Implementing DNS-based rate limiting

Example use case: Return IP A for the first 100 queries, then switch to IP B:

```terraform
resource "snitchdns_record" "canary" {
  zone_id           = snitchdns_zone.example.id
  active            = true
  cls               = "IN"
  type              = "A"
  ttl               = 60
  data              = { address = "192.168.1.100" }
  is_conditional    = true
  conditional_limit = 100
  conditional_reset = false
  conditional_data  = { address = "192.168.1.200" }
}
```

## Common Patterns

### Load Balancing with Multiple A Records

```terraform
locals {
  servers = ["192.168.1.101", "192.168.1.102", "192.168.1.103"]
}

resource "snitchdns_record" "lb" {
  for_each = toset(local.servers)

  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 300
  data    = { address = each.value }
}
```

### Dynamic Record Creation

```terraform
variable "subdomains" {
  type = map(string)
  default = {
    "www"  = "192.168.1.100"
    "api"  = "192.168.1.101"
    "blog" = "192.168.1.102"
  }
}

resource "snitchdns_record" "subdomains" {
  for_each = var.subdomains

  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 3600
  data    = { address = each.value }
}
```
