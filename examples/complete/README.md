# Complete Example

This example demonstrates a complete web infrastructure setup with the SnitchDNS Terraform provider, including web servers, mail servers, and various DNS record types.

## What This Example Creates

- A DNS zone for your domain
- Root A record (example.com → web server)
- WWW CNAME (www.example.com → example.com)
- Mail server A record (mail.example.com → mail server)
- MX record for email routing
- SPF TXT record for email authentication
- DMARC TXT record for email policy
- API subdomain A record
- CAA record for certificate authority authorization
- NS records for name servers

## Prerequisites

1. A running SnitchDNS server
2. An API key from your SnitchDNS instance
3. IP addresses for your web, mail, and API servers

## Usage

1. Copy the example variables file:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Edit `terraform.tfvars` with your actual values:
   - SnitchDNS API URL and key
   - Your domain name
   - IP addresses for your servers

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

## Variables

See `variables.tf` for all configurable options. Key variables:

- `snitchdns_url` - Your SnitchDNS API endpoint
- `snitchdns_key` - Your API key (sensitive)
- `domain` - Your domain name
- `web_server_ip` - IP for web traffic
- `mail_server_ip` - IP for mail server
- `api_server_ip` - IP for API server

## Security Note

Never commit `terraform.tfvars` to version control as it contains sensitive information. The `.gitignore` file should exclude this file.

## Cleanup

To destroy all resources:

```bash
terraform destroy
```
