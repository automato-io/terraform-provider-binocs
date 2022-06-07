# terraform-provider-binocs

Terraform provider for [Binocs](https://binocs.sh)

## Docs

https://registry.terraform.io/providers/automato-io/binocs/latest/docs

## Resources

| TYPE | NAME | DESCRIPTION |
|---|---|---|
| **resource** |`binocs_check`| HTTP(S) or TCP resource check |
| **resource** |`binocs_channel`| Notification channel |

## Example usage

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
  interval = 60
  target   = 1.2
  up_codes = "200-302"
  
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
