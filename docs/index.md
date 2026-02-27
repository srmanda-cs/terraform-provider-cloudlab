---
page_title: "CloudLab Provider"
description: |-
  The CloudLab provider manages resources on CloudLab (cloudlab.us), the academic cloud and network testbed.
---

# CloudLab Provider

The **CloudLab provider** allows [Terraform](https://terraform.io) to manage resources on [CloudLab](https://www.cloudlab.us/) — the academic cloud and network testbed operated by the University of Utah, Clemson University, and the University of Wisconsin.

## Authentication

The provider authenticates to the CloudLab Portal API using an API token.

To obtain a token:
1. Log in to [cloudlab.us](https://www.cloudlab.us/)
2. Navigate to your profile page
3. Generate an API token under **API Tokens**

The token can be provided via:
- The `token` provider attribute (recommended: use a variable or secret)
- The `CLOUDLAB_TOKEN` environment variable

## Example Usage

```terraform
terraform {
  required_providers {
    cloudlab = {
      source  = "srmanda-cs/cloudlab"
      version = "~> 0.1"
    }
  }
}

provider "cloudlab" {
  token = var.cloudlab_token
}
```

## Schema

### Required

- `token` (String, Sensitive) — CloudLab Portal API token. Can also be set via the `CLOUDLAB_TOKEN` environment variable.

### Optional

- `portal_url` (String) — CloudLab portal base URL. Defaults to `https://www.cloudlab.us`. Can also be set via the `CLOUDLAB_PORTAL_URL` environment variable.
