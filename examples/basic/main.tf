terraform {
  required_providers {
    snitchdns = {
      source  = "EinDev/snitchdns"
      version = "~> 1.0"
    }
  }
}

# Configure the SnitchDNS Provider
# You can also use environment variables:
# export SNITCHDNS_API_URL="http://localhost:8000"
# export SNITCHDNS_API_KEY="your-api-key"
provider "snitchdns" {
  api_url = "http://localhost:8000"
  api_key = "your-api-key-here" # Replace with your actual API key
}

# Create a simple DNS zone
resource "snitchdns_zone" "example" {
  domain     = "example.com"
  active     = true
  catch_all  = false
  forwarding = false
  regex      = false
  tags       = ["example", "basic"]
}

# Create an A record for the root domain
resource "snitchdns_record" "root" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 3600

  data = {
    address = "192.168.1.100"
  }
}

# Create a CNAME record for www
resource "snitchdns_record" "www" {
  zone_id = snitchdns_zone.example.id
  active  = true
  cls     = "IN"
  type    = "CNAME"
  ttl     = 3600

  data = {
    name = "example.com."
  }
}

# Output the zone ID for reference
output "zone_id" {
  value       = snitchdns_zone.example.id
  description = "The ID of the created zone"
}

output "zone_domain" {
  value       = snitchdns_zone.example.domain
  description = "The domain name of the zone"
}
