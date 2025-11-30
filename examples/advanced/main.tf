terraform {
  required_providers {
    snitchdns = {
      source = "EinDev/snitchdns"
      version = "~> 1.0"
    }
  }
}

provider "snitchdns" {
  # Using environment variables for security
  # export SNITCHDNS_API_URL="http://localhost:8000"
  # export SNITCHDNS_API_KEY="your-api-key"
}

# Local variables for dynamic configuration
locals {
  environments = ["dev", "staging", "prod"]
  base_domain  = "example.com"

  # Load balancer IPs for each environment
  lb_ips = {
    dev     = ["192.168.1.10", "192.168.1.11"]
    staging = ["192.168.1.20", "192.168.1.21"]
    prod    = ["192.168.1.30", "192.168.1.31", "192.168.1.32"]
  }
}

# Create zones for each environment
resource "snitchdns_zone" "env_zones" {
  for_each = toset(local.environments)

  domain     = "${each.key}.${local.base_domain}"
  active     = true
  catch_all  = false
  forwarding = false
  regex      = false
  tags       = [each.key, "auto-managed", "terraform"]
}

# Create multiple A records for load balancing (round-robin DNS)
resource "snitchdns_record" "lb_records" {
  for_each = merge([
    for env in local.environments : {
      for idx, ip in local.lb_ips[env] :
      "${env}-${idx}" => {
        zone_id = snitchdns_zone.env_zones[env].id
        ip      = ip
        env     = env
      }
    }
  ]...)

  zone_id = each.value.zone_id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 300  # Low TTL for load balancing

  data = {
    address = each.value.ip
  }
}

# Regex zone for catching subdomain patterns
resource "snitchdns_zone" "regex_zone" {
  domain     = ".*\\.test\\.${local.base_domain}"  # Matches *.test.example.com
  active     = true
  catch_all  = true
  forwarding = false
  regex      = true
  tags       = ["regex", "test", "wildcard"]
}

# Catch-all zone for security monitoring
resource "snitchdns_zone" "honeypot" {
  domain     = "honeypot.${local.base_domain}"
  active     = true
  catch_all  = true  # Catch all subdomains
  forwarding = false
  regex      = false
  tags       = ["security", "monitoring", "honeypot"]
}

# Conditional record for canary deployment
resource "snitchdns_record" "canary" {
  zone_id = snitchdns_zone.env_zones["prod"].id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 60  # Very low TTL for canary

  # Start with old version
  data = {
    address = "192.168.1.100"
  }

  # After 100 queries, switch to new version
  is_conditional    = true
  conditional_limit = 100
  conditional_reset = false

  conditional_data = {
    address = "192.168.1.200"
  }
}

# SRV records for service discovery
resource "snitchdns_record" "sip_srv" {
  zone_id = snitchdns_zone.env_zones["prod"].id
  active  = true
  cls     = "IN"
  type    = "SRV"
  ttl     = 3600

  data = {
    priority = "10"
    weight   = "60"
    port     = "5060"
    target   = "sip.${local.base_domain}."
  }
}

# Multiple MX records with priorities
resource "snitchdns_record" "mx_primary" {
  zone_id = snitchdns_zone.env_zones["prod"].id
  active  = true
  cls     = "IN"
  type    = "MX"
  ttl     = 3600

  data = {
    priority = "10"
    hostname = "mail1.${local.base_domain}."
  }
}

resource "snitchdns_record" "mx_backup" {
  zone_id = snitchdns_zone.env_zones["prod"].id
  active  = true
  cls     = "IN"
  type    = "MX"
  ttl     = 3600

  data = {
    priority = "20"
    hostname = "mail2.${local.base_domain}."
  }
}

# CAA records for multiple certificate authorities
resource "snitchdns_record" "caa_letsencrypt" {
  zone_id = snitchdns_zone.env_zones["prod"].id
  active  = true
  cls     = "IN"
  type    = "CAA"
  ttl     = 86400

  data = {
    flags = "0"
    tag   = "issue"
    value = "letsencrypt.org"
  }
}

resource "snitchdns_record" "caa_wildcard" {
  zone_id = snitchdns_zone.env_zones["prod"].id
  active  = true
  cls     = "IN"
  type    = "CAA"
  ttl     = 86400

  data = {
    flags = "0"
    tag   = "issuewild"
    value = "letsencrypt.org"
  }
}

# Dynamic subdomain creation from a map
variable "services" {
  description = "Map of service names to IP addresses"
  type        = map(string)
  default = {
    "api"    = "192.168.1.50"
    "admin"  = "192.168.1.51"
    "cdn"    = "192.168.1.52"
    "status" = "192.168.1.53"
  }
}

resource "snitchdns_record" "services" {
  for_each = var.services

  zone_id = snitchdns_zone.env_zones["prod"].id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 300

  data = {
    address = each.value
  }
}

# Outputs
output "environment_zones" {
  value = {
    for env, zone in snitchdns_zone.env_zones :
    env => {
      id     = zone.id
      domain = zone.domain
    }
  }
  description = "Map of environment zones"
}

output "canary_record_id" {
  value       = snitchdns_record.canary.id
  description = "ID of the canary deployment record"
}

output "service_domains" {
  value = {
    for name, record in snitchdns_record.services :
    name => "${name}.prod.${local.base_domain}"
  }
  description = "Map of service names to their FQDNs"
}
