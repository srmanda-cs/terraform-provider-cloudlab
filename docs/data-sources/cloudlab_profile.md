---
page_title: "cloudlab_profile Data Source - terraform-provider-cloudlab"
description: |-
  Queries an existing CloudLab profile by its UUID or project,name identifier.
---

# cloudlab_profile (Data Source)

Queries an existing CloudLab profile by its UUID or `project,name` composite identifier. Use this data source to reference profiles that were created outside of Terraform, in a separate Terraform state, or managed manually through the Portal UI.

## Example Usage

### Look up a profile by UUID

```terraform
data "cloudlab_profile" "existing" {
  id = "a194e2be-1e5b-4617-84de-c4966cb5c578"
}
```

### Look up a profile by project and name

```terraform
data "cloudlab_profile" "existing" {
  id = "CloudLab,small-lan"
}
```

### Use to instantiate an experiment from an existing profile

```terraform
data "cloudlab_profile" "topology" {
  id = "MyProject,my-topology"
}

resource "cloudlab_experiment" "run" {
  name            = "experiment-1"
  project         = "MyProject"
  profile_name    = data.cloudlab_profile.topology.name
  profile_project = data.cloudlab_profile.topology.project
  duration        = 24
}
```

## Schema

### Required

- `id` (String) — The unique identifier of the profile to look up. Accepts either a UUID (`a194e2be-1e5b-4617-84de-c4966cb5c578`) or a `project,name` composite (`MyProject,my-topology`).

### Read-Only

- `name` (String) — The name of the profile.
- `project` (String) — The CloudLab project that owns the profile.
- `creator` (String) — The CloudLab username who created the profile.
- `version` (Number) — The current version number of the profile.
- `created_at` (String) — The timestamp when the profile was created (ISO 8601).
- `updated_at` (String) — The timestamp when the profile was last updated (ISO 8601, nullable).
- `repository_url` (String) — The URL of the git repository (for repository-backed profiles, nullable).
- `repository_refspec` (String) — The current refspec (for repository-backed profiles, nullable).
- `repository_hash` (String) — The current commit hash (for repository-backed profiles, nullable).
- `repository_githook` (String) — The Portal webhook URL for the repository (for repository-backed profiles, nullable).
- `public` (Boolean) — Whether the profile can be instantiated by any CloudLab user.
- `project_writable` (Boolean) — Whether other members of the project can modify this profile.
