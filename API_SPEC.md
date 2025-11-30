# SnitchDNS API Specification

Based on analysis of SnitchDNS v1 API (https://github.com/sadreck/SnitchDNS)

## Authentication

All API requests require authentication via API key.

**Header:** `X-SnitchDNS-Auth: <api_key>`

Authentication is managed through the API Key system:
- Users can have multiple API keys
- Each key has a name, enabled status, and belongs to a specific user
- Keys can be created, listed, and deleted
- Admin users can access all keys

## Base URL

`/api/v1`

## Core Resources

### 1. Zones (DNS Domains)

Zones represent DNS domains managed by SnitchDNS. Each zone can be queried by either ID or domain name.

#### Zone Properties
- `id` (integer) - Unique zone identifier
- `user_id` (integer) - Owner user ID
- `domain` (string) - Domain name
- `active` (boolean) - Whether zone is active
- `catch_all` (boolean) - Catch-all DNS queries
- `forwarding` (boolean) - Forward queries to upstream DNS
- `regex` (boolean) - Use regex matching for domain
- `master` (boolean) - Master zone (unique subdomain for non-admin users)
- `tags` (string) - Comma-separated tags
- `created_at` (timestamp)
- `updated_at` (timestamp)

#### Endpoints

**GET /zones**
- List all zones user has access to
- Query parameters:
  - `page` (integer) - Page number (default: 1)
  - `per_page` (integer) - Items per page (default: 50)
  - `search` (string) - Search domains
  - `tags` (string) - Comma-separated tags to filter
- Returns: Paginated list of zones

**POST /zones**
- Create new zone
- Required fields: `domain`, `active`, `catch_all`, `forwarding`, `regex`, `master`, `tags`
- Returns: Created zone object

**GET /zones/{zone}**
- Get specific zone by ID or domain
- `{zone}` can be zone ID (integer) or domain name (string)
- Returns: Zone object

**POST /zones/{zone}**
- Update zone
- Optional fields: `domain`, `active`, `catch_all`, `forwarding`, `regex`, `tags`
- Note: `master` property cannot be updated via API
- Returns: Updated zone object

**DELETE /zones/{zone}**
- Delete zone
- Returns: Success response

---

### 2. Records (DNS Records)

DNS records within zones. Records can have conditional responses based on query count.

#### Record Properties
- `id` (integer) - Unique record identifier
- `zone_id` (integer) - Parent zone ID
- `active` (boolean) - Whether record is active
- `cls` (string) - DNS class (e.g., "IN")
- `type` (string) - DNS record type (A, AAAA, CNAME, MX, TXT, etc.)
- `ttl` (integer) - Time to live in seconds
- `data` (object) - Record data (type-specific properties)
- `is_conditional` (boolean) - Enable conditional responses
- `conditional_count` (integer) - Current count for conditional logic
- `conditional_limit` (integer) - Limit for conditional responses
- `conditional_reset` (boolean) - Reset counter on limit
- `conditional_data` (object) - Data for conditional responses

#### Endpoints

**GET /zones/{zone}/records**
- Get all records for a zone
- Returns: Array of record objects

**POST /zones/{zone}/records**
- Create new record
- Required fields: `class`, `type`, `ttl`, `active`, `data`, `is_conditional`, `conditional_count`, `conditional_limit`, `conditional_reset`, `conditional_data`
- `data` object structure depends on record type
- Returns: Created record object

**GET /zones/{zone}/records/{id}**
- Get specific record
- Returns: Record object

**POST /zones/{zone}/records/{id}**
- Update record
- All fields optional (will use existing values if not provided)
- Returns: Updated record object

**DELETE /zones/{zone}/records/{id}**
- Delete record
- Returns: Success response

**GET /records/classes**
- Get supported DNS classes
- Returns: Array of class names

**GET /records/types**
- Get supported DNS record types
- Returns: Array of type names

---

### 3. Restrictions (IP Access Control)

IP-based access control for zones (allow/block lists).

#### Restriction Properties
- `id` (integer) - Unique restriction identifier
- `ip` (string) - IP address or CIDR range
- `type` (string) - "allow" or "block"
- `enabled` (boolean) - Whether restriction is active

#### Endpoints

**GET /zones/{zone}/restrictions**
- Get all restrictions for zone
- Returns: Array of restriction objects

**POST /zones/{zone}/restrictions**
- Create new restriction
- Required fields: `type`, `enabled`, `ip_or_range`
- `type` must be "allow" or "block"
- `ip_or_range` can be single IP or CIDR (e.g., "192.168.0.0/24")
- Returns: Created restriction object

**GET /zones/{zone}/restrictions/{id}**
- Get specific restriction
- Returns: Restriction object

**POST /zones/{zone}/restrictions/{id}**
- Update restriction
- Optional fields: `type`, `enabled`, `ip_or_range`
- Returns: Updated restriction object

**DELETE /zones/{zone}/restrictions/{id}**
- Delete restriction
- Returns: Success response

---

### 4. Notifications

Notification subscriptions for zones (email, webhooks, etc.).

#### Notification Provider Properties
- `id` (integer) - Provider type ID
- `name` (string) - Provider name (e.g., "email", "webhook")
- `enabled` (boolean) - Whether provider is enabled globally

#### Notification Subscription Properties
- `zone_id` (integer) - Zone ID
- `type_id` (integer) - Provider type ID
- `type` (string) - Provider name
- `enabled` (boolean) - Whether subscription is enabled
- `data` (mixed) - Provider-specific configuration data
  - For email: array of email addresses
  - For other providers: string or object

#### Endpoints

**GET /notifications/providers**
- Get all available notification providers
- Returns: Array of notification provider objects

**GET /zones/{zone}/notifications**
- Get all notification subscriptions for zone
- Returns: Array of notification subscription objects

**GET /zones/{zone}/notifications/{provider}**
- Get specific notification subscription
- `{provider}` is the provider name (e.g., "email")
- Returns: Notification subscription object

**POST /zones/{zone}/notifications/{provider}**
- Update notification subscription
- Optional fields: `enabled`, `data`
- `data` format depends on provider type
- Returns: Updated notification subscription object

---

### 5. Search (DNS Query Logs)

Search historical DNS query logs.

#### Search Result Properties
- `page` (integer) - Current page
- `pages` (integer) - Total pages
- `count` (integer) - Total results
- `results` (array) - Array of search result items

#### Search Result Item Properties
- `id` (integer) - Log entry ID
- `domain` (string) - Queried domain
- `source_ip` (string) - Client IP address
- `type` (string) - DNS query type
- `matched` (boolean) - Whether query matched a record
- `forwarded` (boolean) - Whether query was forwarded
- `blocked` (boolean) - Whether query was blocked
- `date` (string) - Timestamp
- `zone_id` (integer) - Matched zone ID (if any)
- `record_id` (integer) - Matched record ID (if any)

#### Endpoints

**GET /search**
- Search DNS query logs
- Query parameters:
  - `domain` (string) - Filter by domain
  - `source_ip` (string) - Filter by source IP
  - `type` (string) - Filter by query type
  - `class` (string) - Filter by DNS class
  - `matched` (boolean) - Filter by match status
  - `forwarded` (boolean) - Filter by forward status
  - `blocked` (boolean) - Filter by block status
  - `user_id` (integer) - Filter by user
  - `tags` (string) - Filter by zone tags
  - `alias` (string) - Filter by alias
  - `date_from` (string) - Start date filter
  - `time_from` (string) - Start time filter
  - `date_to` (string) - End date filter
  - `time_to` (string) - End time filter
  - `page` (integer) - Page number
  - `per_page` (integer) - Items per page
- Returns: Paginated search results

---

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "OK"
}
```

### Error Response
```json
{
  "success": false,
  "code": 5000,
  "message": "Error message",
  "details": "Additional details"
}
```

### Standard Error Codes
- `5000` - Missing required fields
- `5001` - Domain cannot be empty
- `5002` - Could not save tags
- `5003` - Domain already exists / Could not create/save zone
- `5004` - Invalid incoming data
- `5005` - Invalid field value (class, type, TTL, etc.)
- `5006` - Invalid notification type
- `5007` - Invalid notification subscription
- `5008` - No data sent
- `5009` - Notification provider is disabled

### HTTP Status Codes
- `200` - Success
- `401` - Access Denied (invalid or missing API key)
- `404` - Not Found
- `500` - Server Error

---

## Notes for Terraform Provider Implementation

### Resource Hierarchy
1. **Zone** (parent resource) - Can be imported by ID or domain
2. **Record** (child of zone) - Requires zone reference
3. **Restriction** (child of zone) - Requires zone reference
4. **Notification** (child of zone) - Requires zone reference

### Zone Parameter Flexibility
The API accepts zones by either ID or domain name in URL paths. This provides flexibility:
- Use `zone_id` when working with known IDs
- Use `domain` for more readable configurations

### Record Data Structure
Record `data` field is dynamic based on record type. The API endpoint `/records/types` provides available types, and each type has specific required properties. This will need special handling in the Terraform schema.

### Conditional Records
Records support conditional responses based on query count:
- Can serve different responses based on how many times queried
- Counter can reset after hitting limit
- Useful for detection/testing scenarios

### Master Zones
Non-admin users have a "master" zone (unique subdomain) that cannot be deleted or have its `master` flag changed. Admin users don't have master zones.

### Pagination
Zone listing and search results are paginated:
- Default page size: 50
- Returns: `page`, `pages`, `per_page`, `total`

### Tag System
Zones support comma-separated tags for organization and filtering.

### Restrictions
IP restrictions use CIDR notation for ranges and support both allow and block modes (type 1 = allow, type 2 = block internally, but API uses string values "allow"/"block").

### Notification Data Types
Notification `data` field format varies by provider:
- Email: Array of email addresses
- Others: Provider-specific format (string or object)
