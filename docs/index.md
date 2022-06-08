---
page_title: "Binocs Provider"
subcategory: "Monitoring"
description: |-
    This Terraform Provider manages [Binocs](https://binocs.sh) uptime monitoring checks and notification channels. Binocs is a CLI-first uptime and performance monitoring tool for websites, applications and APIs.
---

# Binocs Provider

This Terraform Provider manages [Binocs](https://binocs.sh) uptime monitoring checks and notification channels. Binocs is a CLI-first uptime and performance monitoring tool for websites, applications and APIs.

## Example Usage

```hcl
terraform {
  required_providers {
    binocs = {
      source = "automato-io/binocs"
    }
  }
}

provider "binocs" {
  access_key = "<YOUR_BINOCS_ACCESS_KEY>"
  secret_key = "<YOUR_BINOCS_SECRET_KEY>"
}

# configure a check

resource "binocs_check" "my_website" {
  name     = "My website"
  resource = "https://example.com"
  method   = "GET"
  up_codes = "200-302"
  interval = 60
  target   = 1.2
  
  regions = [
      "us-east-1",
      "us-west-1",
      "us-central-1",
      "ap-southeast-1",
  ]

  up_confirmations_threshold   = 3
  down_confirmations_threshold = 2
}

# configure a notification channel

resource "binocs_channel" "olivia_email" {
  alias  = "Olivia (NL)"
  handle = "olivia@example.com"
  type   = "email"
  checks = [
      binocs_check.my_website.id,
  ]
}
```

## Schema

### Optional

See [Getting Started with Binocs](https://binocs.sh/?#get-started) to learn how to obtain your Access and Secret Keys.

- `access_key` (String) Access Key required to communicate with Binocs API.
- `secret_key` (String) Secret Key required to communicate with Binocs API.
