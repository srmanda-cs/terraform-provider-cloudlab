---
page_title: "cloudlab_experiment Resource - terraform-provider-cloudlab"
description: |-
  Manages a CloudLab experiment. Creating an experiment provisions physical or virtual machines on CloudLab.
---

# cloudlab_experiment (Resource)

Manages a CloudLab experiment. Creating an experiment provisions physical or virtual machines on the CloudLab testbed using a specified profile (topology template). Deleting the experiment terminates and releases all associated resources.

In CloudLab terminology, an "experiment" is a running instantiation of a profile — it represents actual allocated hardware (or VMs) with assigned IP addresses and hostnames.

## Example Usage

```terraform
resource "cloudlab_experiment" "cluster" {
  name            = "my-cluster"
  project         = "MyProject"
  profile_name    = "small-lan"
  profile_project = "CloudLab"
  duration        = 24
}

output "experiment_id" {
  value = cloudlab_experiment.cluster.id
}
```

### Full example with custom profile

```terraform
resource "cloudlab_profile" "topology" {
  name    = "my-topology"
  project = "MyProject"
  script  = file("${path.module}/profile.py")
}

resource "cloudlab_experiment" "run" {
  name            = "experiment-run1"
  project         = "MyProject"
  profile_name    = cloudlab_profile.topology.name
  profile_project = cloudlab_profile.topology.project
  duration        = 48
}
```

## Schema

### Required

- `name` (String) — A human-readable name for the experiment. Must be unique within the project.
- `project` (String) — The CloudLab project to instantiate the experiment in.
- `profile_name` (String) — The name of the profile (topology template) used to create the experiment.
- `profile_project` (String) — The project that owns the profile.

### Optional

- `duration` (Number) — Initial experiment duration in hours.
- `start_at` (String) — Schedule the experiment to start at a future time (RFC3339 format).
- `stop_at` (String) — Schedule the experiment to stop at a future time (RFC3339 format).
- `wait_for_ready` (Boolean) — If `true` (default), Terraform will wait until the experiment reaches `ready` status before completing. Set to `false` to return immediately after creation is submitted.

### Read-Only

- `id` (String) — The unique identifier (UUID) of the experiment assigned by CloudLab.
- `creator` (String) — The CloudLab username who created the experiment.
- `status` (String) — The current status of the experiment (e.g., `created`, `ready`, `failed`).
- `created_at` (String) — The timestamp when the experiment was created.
- `expires_at` (String) — The timestamp when the experiment is scheduled to expire.

## Timeouts

When `wait_for_ready` is `true`, the provider polls for experiment readiness every 15 seconds with a maximum timeout of 30 minutes.

## Import

Experiments can be imported using their UUID:

```shell
terraform import cloudlab_experiment.cluster <experiment-uuid>
```
