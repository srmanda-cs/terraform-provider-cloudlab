---
page_title: "cloudlab_profile Resource - terraform-provider-cloudlab"
description: |-
  Manages a CloudLab experiment profile (topology template).
---

# cloudlab_profile (Resource)

Manages a CloudLab experiment profile. A profile defines the topology template — including node types, hardware specifications, and network configuration — used to instantiate experiments via `cloudlab_experiment`.

Profiles can be defined using a [geni-lib](https://docs.cloudlab.us/geni-lib/intro/intro.html) Python script, an RSpec XML document, or backed by a git repository. Once created, the `script`, `public`, and `project_writable` attributes can be updated **in-place** without recreating the profile (all other attributes force replacement).

## Example Usage

### Inline geni-lib Python script

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
  public         = true
}
```

### Public profile writable by the whole project

```terraform
resource "cloudlab_profile" "shared" {
  name             = "shared-topology"
  project          = "MyProject"
  script           = file("${path.module}/topology.py")
  public           = true
  project_writable = true
}
```

### Updating a script in-place

Changing `script` does not require destroying the profile — it performs a `PATCH /profiles/{id}`:

```terraform
resource "cloudlab_profile" "topology" {
  name    = "my-topology"
  project = "MyProject"
  script  = file("${path.module}/topology-v2.py")  # Updated script
}
```

### Use the profile in an experiment

```terraform
resource "cloudlab_profile" "topology" {
  name    = "my-topology"
  project = "MyProject"
  script  = file("${path.module}/profile.py")
}

resource "cloudlab_experiment" "run" {
  name            = "my-run"
  project         = "MyProject"
  profile_name    = cloudlab_profile.topology.name
  profile_project = cloudlab_profile.topology.project
  duration        = 24
}
```

## Schema

### Required

- `name` (String) — The name of the profile. Must be unique within the project. **Forces new resource.**
- `project` (String) — The CloudLab project that owns this profile. **Forces new resource.**

### Optional

- `script` (String) — A geni-lib Python script that defines the experiment topology. Mutually exclusive with `repository_url`. **Mutable in-place** via `PATCH /profiles/{id}`.
- `repository_url` (String) — URL of a git repository containing the profile definition. Mutually exclusive with `script`. **Forces new resource.**
- `public` (Boolean) — If `true`, the profile can be instantiated by any CloudLab user. Defaults to `false`. **Mutable in-place.**
- `project_writable` (Boolean) — If `true`, other members of the project can modify this profile. Defaults to `false`. **Mutable in-place.**

### Read-Only

- `id` (String) — The unique identifier (UUID) of the profile assigned by CloudLab.
- `creator` (String) — The CloudLab username who created the profile.
- `version` (Number) — The current version number of the profile. Increments each time the script is updated.
- `created_at` (String) — The timestamp when the profile was created (ISO 8601).
- `updated_at` (String) — The timestamp when the profile was last updated (ISO 8601, nullable).
- `repository_url` (String) — The repository URL (for repository-backed profiles, nullable).
- `repository_refspec` (String) — The current refspec of the profile (for repository-backed profiles, nullable).
- `repository_hash` (String) — The current commit hash of the profile (for repository-backed profiles, nullable).
- `repository_githook` (String) — The Portal webhook URL for the repository (for repository-backed profiles, nullable). Configure this as a webhook in your git repository to trigger automatic profile updates on push.

## Notes

- Either `script` or `repository_url` must be provided (but not both).
- For repository-backed profiles, the profile is updated by triggering a pull via `PUT /profiles/{id}`. This is not exposed as a Terraform action — use `make api-update` or the Portal UI to trigger a repo refresh.
- The `version` counter increments each time `script` is changed via `PATCH`.

## Import

Profiles can be imported using their UUID or `project,name`:

```shell
terraform import cloudlab_profile.topology <profile-uuid>
terraform import cloudlab_profile.topology "MyProject,my-topology"
```
