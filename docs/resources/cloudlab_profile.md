---
page_title: "cloudlab_profile Resource - terraform-provider-cloudlab"
description: |-
  Manages a CloudLab experiment profile (topology template).
---

# cloudlab_profile (Resource)

Manages a CloudLab experiment profile. A profile defines the topology template — including node types, hardware specifications, and network configuration — used to instantiate experiments.

Profiles can be defined using a [geni-lib](https://docs.cloudlab.us/geni-lib/intro/intro.html) Python script or backed by a git repository.

## Example Usage

### Inline geni-lib script

```terraform
resource "cloudlab_profile" "two_nodes" {
  name    = "two-node-lan"
  project = "MyProject"
  script  = file("${path.module}/profile.py")
}
```

### Repository-backed profile

```terraform
resource "cloudlab_profile" "from_repo" {
  name           = "repo-backed-profile"
  project        = "MyProject"
  repository_url = "https://github.com/example/cloudlab-profiles.git"
}
```

## Schema

### Required

- `name` (String) — The name of the profile. Must be unique within the project.
- `project` (String) — The CloudLab project that owns this profile.

### Optional

- `script` (String) — A geni-lib Python script that defines the experiment topology. Mutually exclusive with `repository_url`.
- `repository_url` (String) — URL of a git repository containing the profile definition. Mutually exclusive with `script`.
- `public` (Boolean) — If `true`, the profile can be instantiated by any CloudLab user. Defaults to `false`.
- `project_writable` (Boolean) — If `true`, other members of the project can modify this profile. Defaults to `false`.

### Read-Only

- `id` (String) — The unique identifier (UUID) of the profile assigned by CloudLab.
- `creator` (String) — The CloudLab username who created the profile.
- `version` (Number) — The current version number of the profile.
- `created_at` (String) — The timestamp when the profile was created.
- `updated_at` (String) — The timestamp when the profile was last updated.
