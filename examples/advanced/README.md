# Advanced Example

This example demonstrates advanced usage patterns with the SnitchDNS Terraform provider, including dynamic resource creation, regex zones, conditional records, and multi-environment setups.

## What This Example Demonstrates

### Multi-Environment Setup
- Creates separate zones for dev, staging, and prod
- Dynamic zone creation using `for_each`
- Environment-specific tagging

### Load Balancing
- Multiple A records per environment for round-robin DNS
- Different numbers of backend servers per environment
- Low TTL values for faster failover

### Advanced Zone Types
- **Regex Zone**: Matches patterns like `*.test.example.com`
- **Catch-All Zone**: Honeypot for security monitoring
- Both zones with appropriate tags

### Canary Deployment
- Conditional DNS record that switches after N queries
- Useful for gradual rollouts
- Low TTL for quick updates

### Service Discovery
- SRV records for service location
- Multiple MX records with priorities
- Dynamic subdomain creation from variables

### Certificate Authority Authorization
- CAA records for Let's Encrypt
- Wildcard certificate control
- Multiple CAA records for different purposes

### Dynamic Configuration
- Service subdomains created from a map variable
- Flexible infrastructure-as-code patterns
- Easy to add new services

## Prerequisites

1. A running SnitchDNS server
2. An API key from your SnitchDNS instance
3. Understanding of advanced DNS concepts

## Usage

1. Set environment variables:
   ```bash
   export SNITCHDNS_API_URL="http://localhost:8000"
   export SNITCHDNS_API_KEY="your-api-key"
   ```

2. (Optional) Customize the services map in `main.tf` or create a `terraform.tfvars`:
   ```hcl
   services = {
     "api"   = "192.168.1.50"
     "admin" = "192.168.1.51"
     "cdn"   = "192.168.1.52"
   }
   ```

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Review the planned changes:
   ```bash
   terraform plan
   ```

5. Apply the configuration:
   ```bash
   terraform apply
   ```

## Key Patterns

### Load Balancing with Round-Robin DNS

Multiple A records for the same hostname distribute traffic across backends:

```hcl
lb_ips = {
  prod = ["192.168.1.30", "192.168.1.31", "192.168.1.32"]
}
```

### Regex Zones for Pattern Matching

Catch all queries matching a regex pattern:

```hcl
resource "snitchdns_zone" "regex_zone" {
  domain = ".*\\.test\\.example\\.com"
  regex  = true
}
```

### Conditional Records for Canary Deployments

Gradually shift traffic from old to new version:

```hcl
resource "snitchdns_record" "canary" {
  data              = { address = "192.168.1.100" }  # Old
  is_conditional    = true
  conditional_limit = 100
  conditional_data  = { address = "192.168.1.200" }  # New
}
```

### Dynamic Resource Creation

Create resources from maps for easy management:

```hcl
variable "services" {
  type = map(string)
}

resource "snitchdns_record" "services" {
  for_each = var.services
  # ...
}
```

## Security Considerations

### Honeypot Zone

The catch-all honeypot zone can help detect:
- DNS exfiltration attempts
- Subdomain enumeration
- Malware communication patterns

Monitor queries to this zone for security insights.

### CAA Records

CAA records restrict which certificate authorities can issue certificates:
- Prevents unauthorized certificate issuance
- Protects against some phishing attacks
- Industry best practice

## Performance Tips

### TTL Strategy

- **Low TTL (60-300s)**: Load balanced records, canary deployments
- **Medium TTL (3600s)**: Standard records, regular services
- **High TTL (86400s)**: Rarely changing records (NS, CAA)

### Resource Organization

Use tags consistently for:
- Environment identification
- Service categorization
- Management automation

## Cleanup

To destroy all resources:

```bash
terraform destroy
```

**Note**: This will destroy all zones and records created by this configuration across all environments.

## Further Reading

- [Terraform for_each documentation](https://www.terraform.io/language/meta-arguments/for_each)
- [DNS record types reference](https://en.wikipedia.org/wiki/List_of_DNS_record_types)
- [CAA records explained](https://letsencrypt.org/docs/caa/)
- [SRV records for service discovery](https://en.wikipedia.org/wiki/SRV_record)
