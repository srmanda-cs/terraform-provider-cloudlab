---
page_title: "cloudlab_snapshot Resource - terraform-provider-cloudlab"
description: |-
  Takes a disk image snapshot of a node in a running CloudLab experiment.
---

# cloudlab_snapshot (Resource)

Takes a disk image snapshot of a running node in a CloudLab experiment. The snapshot creates (or updates) a named disk image that can then be used as a custom base image in future experiments.

By default, the provider polls until the snapshot reaches a terminal status (`ready` or `failed`), polling every 15 seconds with a 60-minute timeout.

**Important:** Destroying this resource removes it from Terraform state only. The created image persists in CloudLab and must be deleted manually through the Portal UI. All snapshot attributes force replacement.

## Example Usage

### Basic snapshot

```terraform
resource "cloudlab_experiment" "my_exp" {
  name            = "my-experiment"
  project         = "MyProject"
  profile_name    = "small-lan"
  profile_project = "CloudLab"
}

resource "cloudlab_snapshot" "node1_image" {
  experiment_id = cloudlab_experiment.my_exp.id
  client_id     = "node1"
  image_name    = "my-custom-ubuntu-image"
}

output "image_urn" {
  value = cloudlab_snapshot.node1_image.image_urn
}
```

### Whole-disk snapshot

```terraform
resource "cloudlab_snapshot" "full_disk" {
  experiment_id = cloudlab_experiment.my_exp.id
  client_id     = "node1"
  image_name    = "my-whole-disk-image"
  whole_disk    = true
}
```

### Non-blocking snapshot (fire and forget)

```terraform
resource "cloudlab_snapshot" "async" {
  experiment_id     = cloudlab_experiment.my_exp.id
  client_id         = "node1"
  image_name        = "my-async-image"
  wait_for_complete = false
}

output "snapshot_id" {
  value = cloudlab_snapshot.async.id
}
```

## Schema

### Required

- `experiment_id` (String) — The UUID of the running experiment containing the node to snapshot. **Forces new resource.**
- `client_id` (String) — The logical name (client ID) of the node to snapshot, as defined in the profile (e.g., `node1`, `server`, `worker-0`). **Forces new resource.**
- `image_name` (String) — The name of the image to create or update. If an image with this name already exists in the project, it will be updated. **Forces new resource.**

### Optional

- `whole_disk` (Boolean) — If `true`, take a whole-disk image. If `false` (default), take a partition image. **Forces new resource.**
- `wait_for_complete` (Boolean) — If `true` (default), Terraform blocks until the snapshot reaches a terminal status before completing the `apply`. Set to `false` to return immediately after initiating the snapshot. The provider polls every 15 seconds with a 60-minute timeout.

### Read-Only

- `id` (String) — The unique identifier (UUID) of the snapshot request assigned by CloudLab.
- `status` (String) — The current status of the snapshot operation.
- `status_timestamp` (String) — The timestamp of the last status update (ISO 8601, nullable).
- `image_size` (Number) — The current size of the image in KB (nullable).
- `image_urn` (String) — The URN of the created image. Use this to reference the image in future profile scripts.
- `error_message` (String) — Error message if the snapshot failed (nullable).

## Timeouts

When `wait_for_complete = true` (the default), the provider polls every **15 seconds** with a maximum timeout of **60 minutes**.

## Import

Snapshots can be imported using a composite ID of the form `<experiment_uuid>/<snapshot_uuid>`:

```shell
terraform import cloudlab_snapshot.my_image <experiment-uuid>/<snapshot-uuid>
```

Both UUIDs are required because the CloudLab API needs the experiment ID to look up snapshot status. The snapshot UUID is the value of the `id` attribute on the resource.

**Note:** Importing a snapshot does not re-create the image. It only brings an existing snapshot request into Terraform state so its status can be tracked.
