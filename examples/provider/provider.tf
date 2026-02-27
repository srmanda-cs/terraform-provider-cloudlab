terraform {
  required_providers {
    cloudlab = {
      source  = "srmanda-cs/cloudlab"
      version = "~> 0.1"
    }
  }
}

provider "cloudlab" {
  # API token from your CloudLab portal profile.
  # Can also be set via CLOUDLAB_TOKEN environment variable.
  token = var.cloudlab_token

  # Optional: override the portal URL (defaults to https://www.cloudlab.us)
  # portal_url = "https://www.cloudlab.us"
}

variable "cloudlab_token" {
  description = "CloudLab Portal API token."
  type        = string
  sensitive   = true
}
