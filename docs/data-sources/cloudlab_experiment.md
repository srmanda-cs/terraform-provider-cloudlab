---
page_title: "cloudlab_experiment Data Source - terraform-provider-cloudlab"
description: |-
  Queries an existing CloudLab experiment by its UUID.
---

# cloudlab_experiment (Data Source)

Queries an existing CloudLab experiment by its UUID. Use this data source to reference experiments that were created outside of Terraform or in a separate Terraform state.

## Example Usage

```terraform
data "cloudlab_experiment" "existing" {
  id = "550e8400-e29b-41d4-a716-446655440000"
}

output "experiment_status" {
  value = data.cloudlab_experiment.existing.status
}
```

## Schema

### Required

- `id` (String) — The unique identifier (UUID) of the experiment to look up.

### Read-Only

- `name` (String) — The human-readable name of the experiment.
- `project` (String) — The CloudLab project the experiment belongs to.
- `profile_name` (String) — The name of the profile used to create the experiment.
- `profile_project` (String) — The project that owns the profile.
- `creator` (String) — The CloudLab username who created the experiment.
- `status` (String) — The current status of the experiment.
- `created_at` (String) — The timestamp when the experiment was created.
- `expires_at` (String) — The timestamp when the experiment is scheduled to expire.
