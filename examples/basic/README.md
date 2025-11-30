# Basic Example

This example demonstrates the basic usage of the SnitchDNS Terraform provider.

## What This Example Creates

- A DNS zone for `example.com`
- An A record for the root domain pointing to 192.168.1.100
- A CNAME record for www pointing to the root domain

## Prerequisites

1. A running SnitchDNS server
2. An API key from your SnitchDNS instance (Settings > API)

## Usage

1. Update the `api_url` and `api_key` in `main.tf` with your actual values
2. Initialize Terraform:
   ```bash
   terraform init
   ```
3. Review the planned changes:
   ```bash
   terraform plan
   ```
4. Apply the configuration:
   ```bash
   terraform apply
   ```

## Using Environment Variables

Instead of hardcoding credentials in `main.tf`, you can use environment variables:

```bash
export SNITCHDNS_API_URL="http://localhost:8000"
export SNITCHDNS_API_KEY="your-api-key"
terraform apply
```

## Cleanup

To destroy all resources created by this example:

```bash
terraform destroy
```
