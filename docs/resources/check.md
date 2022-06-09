---
page_title: "binocs_check Resource - terraform-provider-binocs"
subcategory: "Monitoring"
description: |-
  Use `binocs_check` Resource to manage a Binocs (https://binocs.sh) uptime monitoring check.
---

# binocs_check (Resource)

Use `binocs_check` Resource to manage a [Binocs](https://binocs.sh) uptime monitoring check.

## Schema

### Required

- `resource` (String) The resource to check; a valid URL in case of a HTTP(S) resource, or tcp://{HOSTNAME}:{PORT} in case of a TCP resource.

### Optional

- `name` (String) The name (alias) of this check. Maximum length is 25 characters.
- `method` (String) The HTTP method (one of "GET", "HEAD", "POST", "PUT", "DELETE"). Only used and required with HTTP(S) resources.
- `interval` (Number) How often Binocs checks this resource, in seconds. Minimum is 5 and maximum is 900 seconds.
- `regions` (Set of String) From where in the world Binocs checks this resource. At least one region is required. Valid values: af-south-1, ap-east-1, ap-northeast-1, ap-south-1, ap-southeast-1, ap-southeast-2, eu-central-1, eu-west-1, sa-east-1, us-east-1, us-west-1.
- `target` (Number) The response time that accommodates Apdex=1.0 (in seconds with up to 3 decimal places). Valid target is a value between 0.01 and 10.0 seconds.
- `up_codes` (String) The good ("up") HTTP(S) response codes, e.g. 2xx or `200-302`, or `200,301`. Only used with HTTP(S) resources.
- `up_confirmations_threshold` (String) How many subsequent "up" responses need to occur before Binocs creates an incident and triggers notifications. Minimum is 1 and maximum is 10.
- `down_confirmations_threshold` (String) How many subsequent "down" responses need to occur before Binocs closes an incident and triggers "recovery" notifications. Minimum is 1 and maximum is 10.

### Read-Only

- `id` (String) The ID of this resource.


