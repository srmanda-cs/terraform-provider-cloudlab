---
page_title: "cloudlab_resgroup Resource - terraform-provider-cloudlab"
description: |-
  Manages a CloudLab reservation group for pre-reserving hardware resources.
---

# cloudlab_resgroup (Resource)

Manages a CloudLab reservation group. Reservation groups allow you to pre-reserve specific hardware resources on CloudLab for a defined time window, ensuring the hardware is available when you need to run experiments.

## Example Usage

```terraform
resource "cloudlab_resgroup" "hw_reservation" {
  project    = "MyProject"
  reason     = "Weekly ML experiment requiring guaranteed xl170 hardware"
  expires_at = "2026-03-01T00:00:00Z"

  node_types = [
    {
      node_type = "xl170"
      aggregate = "utah.cloudlab.us"
      count     = 4
    }
  ]
}
```

## Schema

### Required

- `project` (String) — The CloudLab project for this reservation group.
- `reason` (String) — A description of why you need to reserve these resources.

### Optional

- `start_at` (String) — The time the reservation should start (RFC3339 format). If omitted, the reservation starts immediately.
- `expires_at` (String) — The time the reservation expires (RFC3339 format).
- `node_types` (List of Object) — The list of node types and counts to reserve. Each entry has:
  - `node_type` (String) — The hardware node type to reserve (e.g., `xl170`, `m400`, `r320`).
  - `aggregate` (String) — The CloudLab site/aggregate (e.g., `utah.cloudlab.us`, `clemson.cloudlab.us`).
  - `count` (Number) — The number of nodes of this type to reserve.

### Read-Only

- `id` (String) — The unique identifier (UUID) of the reservation group assigned by CloudLab.
- `creator` (String) — The CloudLab username who created the reservation group.
- `status` (String) — The current status of the reservation group.
- `created_at` (String) — The timestamp when the reservation group was created.
