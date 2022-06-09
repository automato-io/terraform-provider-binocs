---
page_title: "binocs_channel Resource - terraform-provider-binocs"
subcategory: "Monitoring"
description: |-
  Use `binocs_channel` Resource to manage a Binocs (https://binocs.sh) uptime monitoring notification channel.
---

# binocs_channel (Resource)

Use `binocs_channel` Resource to manage a [Binocs](https://binocs.sh) uptime monitoring notification channel.

## Schema

### Required

- `type` (String) The only supported channel is currently "email", and it requires e-mail address verification. All other notification channels ("slack", "telegram") currently require interactive creation using Binocs CLI. All notification channels can be imported to Terraform.
- `handle` (String) The e-mail address for a channel of `type = email`.

### Optional

- `alias` (String) The alias (name) of this notification channel. Maximum length is 25 characters.
- `checks` (Set of String) The checks to associate with this notifications channel.

### Read-Only

- `id` (String) The ID of this resource. 


