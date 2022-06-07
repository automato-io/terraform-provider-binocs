---
page_title: "binocs_check Resource - terraform-provider-binocs"
subcategory: "Monitoring"
description: |-
  Use `binocs_check` Resource to manage a [Binocs](https://binocs.sh) uptime monitoring check.
---

# binocs_check (Resource)

Use `binocs_check` Resource to manage a [Binocs](https://binocs.sh) uptime monitoring check.

## Schema

### Required

- `resource` (String) The resource to check; a URL in case of a HTTP(S) resource, or HOSTNAME:PORT in case of a TCP resource.

### Optional

- `name` (String) The name (alias) of this check.
- `method` (String) The HTTP method (one of GET, HEAD, POST, PUT or DELETE). Only required for HTTP(S) resources.
- `interval` (Number) How often Binocs checks this resource (in seconds).
- `regions` (Set of String) From where in the world Binocs checks this resource.
- `target` (Number) The response time that accommodates Apdex=1.0 (in seconds with up to 3 decimal places).
- `up_codes` (String) The good ("up") HTTP(S) response codes, e.g. 2xx or `200-302`, or `200,301`
- `up_confirmations_threshold` (String) How many subsequent "up" responses need to occur before Binocs creates an incident and triggers notifications.
- `down_confirmations_threshold` (String) How many subsequent "down" responses need to occur before Binocs closes an incident and triggers "recovery" notifications.

### Read-Only

- `id` (String) The ID of this resource.


