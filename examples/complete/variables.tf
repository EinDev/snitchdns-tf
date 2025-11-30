variable "snitchdns_url" {
  description = "SnitchDNS API URL"
  type        = string
  default     = "http://localhost:8000"
}

variable "snitchdns_key" {
  description = "SnitchDNS API Key"
  type        = string
  sensitive   = true
}

variable "domain" {
  description = "Domain name for the zone"
  type        = string
  default     = "example.com"
}

variable "web_server_ip" {
  description = "IP address of the web server"
  type        = string
  default     = "192.168.1.100"
}

variable "mail_server_ip" {
  description = "IP address of the mail server"
  type        = string
  default     = "192.168.1.101"
}

variable "api_server_ip" {
  description = "IP address of the API server"
  type        = string
  default     = "192.168.1.102"
}
