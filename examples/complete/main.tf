terraform {
  required_providers {
    snitchdns = {
      source  = "EinDev/snitchdns"
      version = "~> 1.0"
    }
  }
}

provider "snitchdns" {
  api_url = var.snitchdns_url
  api_key = var.snitchdns_key
}

# Create a zone for the main domain
resource "snitchdns_zone" "main" {
  domain     = var.domain
  active     = true
  catch_all  = false
  forwarding = false
  regex      = false
  tags       = ["production", "web"]
}

# Root A record
resource "snitchdns_record" "root" {
  zone_id = snitchdns_zone.main.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 3600

  data = {
    address = var.web_server_ip
  }
}

# WWW CNAME
resource "snitchdns_record" "www" {
  zone_id = snitchdns_zone.main.id
  active  = true
  cls     = "IN"
  type    = "CNAME"
  ttl     = 3600

  data = {
    name = "${var.domain}."
  }
}

# Mail server A record
resource "snitchdns_record" "mail" {
  zone_id = snitchdns_zone.main.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 3600

  data = {
    address = var.mail_server_ip
  }
}

# MX record
resource "snitchdns_record" "mx" {
  zone_id = snitchdns_zone.main.id
  active  = true
  cls     = "IN"
  type    = "MX"
  ttl     = 3600

  data = {
    priority = "10"
    hostname = "mail.${var.domain}."
  }
}

# SPF TXT record
resource "snitchdns_record" "spf" {
  zone_id = snitchdns_zone.main.id
  active  = true
  cls     = "IN"
  type    = "TXT"
  ttl     = 3600

  data = {
    data = "v=spf1 mx -all"
  }
}

# DMARC TXT record
resource "snitchdns_record" "dmarc" {
  zone_id = snitchdns_zone.main.id
  active  = true
  cls     = "IN"
  type    = "TXT"
  ttl     = 3600

  data = {
    data = "v=DMARC1; p=quarantine; rua=mailto:dmarc@${var.domain}"
  }
}

# API subdomain
resource "snitchdns_record" "api" {
  zone_id = snitchdns_zone.main.id
  active  = true
  cls     = "IN"
  type    = "A"
  ttl     = 300 # Lower TTL for API

  data = {
    address = var.api_server_ip
  }
}

# CAA record to restrict certificate issuance
resource "snitchdns_record" "caa" {
  zone_id = snitchdns_zone.main.id
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

# Name servers
resource "snitchdns_record" "ns1" {
  zone_id = snitchdns_zone.main.id
  active  = true
  cls     = "IN"
  type    = "NS"
  ttl     = 86400

  data = {
    name = "ns1.${var.domain}."
  }
}

resource "snitchdns_record" "ns2" {
  zone_id = snitchdns_zone.main.id
  active  = true
  cls     = "IN"
  type    = "NS"
  ttl     = 86400

  data = {
    name = "ns2.${var.domain}."
  }
}

# Outputs
output "zone_id" {
  value       = snitchdns_zone.main.id
  description = "The ID of the created zone"
}

output "zone_domain" {
  value       = snitchdns_zone.main.domain
  description = "The domain name of the zone"
}

output "name_servers" {
  value = [
    snitchdns_record.ns1.data["name"],
    snitchdns_record.ns2.data["name"]
  ]
  description = "Name servers for the zone"
}
