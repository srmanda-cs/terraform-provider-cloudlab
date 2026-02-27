---
page_title: "cloudlab_vlan_connection Resource - terraform-provider-cloudlab"
description: |-
  Manages a shared VLAN connection between two running CloudLab experiments.
---

# cloudlab_vlan_connection (Resource)

Manages a shared VLAN connection between two running CloudLab experiments. Creates a layer-2 connection between a LAN in one experiment and a LAN in another, enabling direct network communication between nodes across experiment boundaries.

Both experiments must be running and their profiles must declare the LAN as a shared VLAN. Destroying this resource calls the disconnect endpoint to remove the connection.

All attributes force replacement — there is no in-place update for VLAN connections.

## Example Usage

### Connect two experiments via a shared LAN

```terraform
resource "cloudlab_experiment" "exp_a" {
  name            = "experiment-a"
  project         = "MyProject"
  profile_name    = "profile-with-shared-vlan"
  profile_project = "MyProject"
}

resource "cloudlab_experiment" "exp_b" {
  name            = "experiment-b"
  project         = "MyProject"
  profile_name    = "profile-with-shared-vlan"
  profile_project = "MyProject"
}

resource "cloudlab_vlan_connection" "link" {
  experiment_id = cloudlab_experiment.exp_a.id
  source_lan    = "shared-lan"
  target_id     = cloudlab_experiment.exp_b.id
  target_lan    = "shared-lan"
}
```

### Reference an existing experiment via data source

```terraform
data "cloudlab_experiment" "existing" {
  id = "550e8400-e29b-41d4-a716-446655440000"
}

resource "cloudlab_experiment" "new_exp" {
  name            = "new-experiment"
  project         = "MyProject"
  profile_name    = "profile-with-shared-vlan"
  profile_project = "MyProject"
}

resource "cloudlab_vlan_connection" "link" {
  experiment_id = cloudlab_experiment.new_exp.id
  source_lan    = "shared-lan"
  target_id     = data.cloudlab_experiment.existing.id
  target_lan    = "shared-lan"
}
```

## Schema

### Required

- `experiment_id` (String) — The UUID of the source experiment. **Forces new resource.**
- `source_lan` (String) — The client ID of the shared LAN in the source experiment (as defined in the profile). **Forces new resource.**
- `target_id` (String) — The UUID or `project,name` of the target experiment to connect to. **Forces new resource.**
- `target_lan` (String) — The client ID of the shared LAN in the target experiment. **Forces new resource.**

### Read-Only

- `id` (String) — A synthetic identifier for tracking this connection in Terraform state, formatted as `experiment_id/source_lan`.

## Notes

- Both experiments must be in `ready` status for the connection to succeed.
- The LAN must be declared as a shared VLAN in both profile definitions.
- The CloudLab API does not provide a query endpoint for VLAN connection state, so this resource cannot be imported. If a connection already exists, it must be managed by Terraform from the time it is created.

## Import

VLAN connections cannot be imported — the CloudLab API does not expose a query endpoint for connection state.
