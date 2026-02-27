---
page_title: "cloudlab_resgroup Resource - terraform-provider-cloudlab"
description: |-
  Manages a CloudLab reservation group for pre-reserving hardware resources.
---

# cloudlab_resgroup (Resource)

Manages a CloudLab reservation group. Reservation groups allow you to pre-reserve specific hardware resources on CloudLab for a defined time window, ensuring the hardware is available when you need to run experiments.

Reservation groups support three types of reservations (which can be combined):
- **Node type reservations** — Reserve a number of specific node types at a specific aggregate (cluster)
- **Frequency range reservations** — Reserve RF spectrum (Powder/WIRELESS testbed only)
- **Named route reservations** — Reserve specific network routes

**Mutable without replacement:** `reason`, `start_at`, `expires_at`, `duration`, `powder_zones`, `node_types`, `ranges`, and `routes`. The `project` and `group` attributes force replacement.

## Example Usage

### Reserve xl170 nodes at Utah CloudLab

```terraform
resource "cloudlab_resgroup" "compute" {
  project    = "MyProject"
  reason     = "Weekly ML experiment requiring guaranteed xl170 hardware"
  expires_at = "2026-04-01T00:00:00Z"

  node_types = [
    {
      urn       = "urn:publicid:IDN+utah.cloudlab.us+authority+cm"
      node_type = "xl170"
      count     = 4
    }
  ]
}
```

### Reserve nodes at multiple clusters

```terraform
resource "cloudlab_resgroup" "multi_site" {
  project    = "MyProject"
  reason     = "Cross-cluster network experiment"
  expires_at = "2026-04-15T00:00:00Z"

  node_types = [
    {
      urn       = "urn:publicid:IDN+utah.cloudlab.us+authority+cm"
      node_type = "xl170"
      count     = 2
    },
    {
      urn       = "urn:publicid:IDN+clemson.cloudlab.us+authority+cm"
      node_type = "d430"
      count     = 2
    }
  ]
}
```

### Use duration instead of expires_at

```terraform
resource "cloudlab_resgroup" "timed" {
  project  = "MyProject"
  reason   = "24-hour experiment window"
  start_at = "2026-03-15T08:00:00Z"
  duration = 24
}
```

### Powder frequency range reservation (WIRELESS testbed)

```terraform
resource "cloudlab_resgroup" "spectrum" {
  project      = "MyProject"
  reason       = "5G NR experiment requiring dedicated spectrum"
  expires_at   = "2026-04-01T00:00:00Z"
  powder_zones = "Outdoor"

  ranges = [
    {
      min_freq = 3550.0
      max_freq = 3600.0
    }
  ]
}
```

### Combined node and spectrum reservation

```terraform
resource "cloudlab_resgroup" "powder_full" {
  project      = "MyProject"
  reason       = "Full Powder testbed experiment"
  expires_at   = "2026-04-01T00:00:00Z"
  powder_zones = "Indoor OTA Lab"

  node_types = [
    {
      urn       = "urn:publicid:IDN+emulab.net+authority+cm"
      node_type = "nuc5300"
      count     = 2
    }
  ]

  ranges = [
    {
      min_freq = 2400.0
      max_freq = 2500.0
    }
  ]
}
```

## Schema

### Required

- `project` (String) — The CloudLab project for this reservation group. **Forces new resource.**
- `reason` (String) — A description of why you need to reserve these resources. Mutable in-place.

### Optional

- `group` (String) — The project subgroup for this reservation group. **Forces new resource.**
- `start_at` (String) — The time the reservation should start (RFC3339 format, e.g. `2026-03-15T08:00:00Z`). Validated at plan time. If omitted, the reservation starts immediately. Mutable until the reservation is approved.
- `expires_at` (String) — The time the reservation expires (RFC3339 format). Validated at plan time. Mutually exclusive with `duration`. Mutable: once approved, can only be moved earlier, not later.
- `duration` (Number) — Duration of the reservation in hours, as an alternative to `expires_at`. Passed as a query parameter on create/update.
- `powder_zones` (String) — Powder zone for radio resource reservations. Valid values: `Outdoor`, `Indoor OTA Lab`, `Flux`. Only applicable on the Powder/WIRELESS testbed.
- `node_types` (List of Object) — Node type reservations. Each entry has:
  - `urn` (String, Required) — The aggregate URN identifying the cluster (e.g., `urn:publicid:IDN+utah.cloudlab.us+authority+cm`).
  - `node_type` (String, Required) — The hardware node type to reserve (e.g., `xl170`, `d430`, `m400`, `r320`, `d710`).
  - `count` (Number, Required) — The number of nodes of this type to reserve.
- `ranges` (List of Object) — Frequency range reservations (Powder testbed only). Each entry has:
  - `min_freq` (Number, Required) — The start of the frequency range (inclusive) in MHz.
  - `max_freq` (Number, Required) — The end of the frequency range (inclusive) in MHz.
- `routes` (List of Object) — Named route reservations. Each entry has:
  - `name` (String, Required) — The route name to reserve.

### Read-Only

- `id` (String) — The unique identifier (UUID) of the reservation group assigned by CloudLab.
- `creator` (String) — The CloudLab username who created the reservation group.
- `created_at` (String) — The timestamp when the reservation group was created (ISO 8601, nullable).

## Common Aggregate URNs

| Cluster | URN |
|---------|-----|
| Utah CloudLab | `urn:publicid:IDN+utah.cloudlab.us+authority+cm` |
| Clemson CloudLab | `urn:publicid:IDN+clemson.cloudlab.us+authority+cm` |
| Wisconsin CloudLab | `urn:publicid:IDN+wisc.cloudlab.us+authority+cm` |
| Powder/WIRELESS | `urn:publicid:IDN+emulab.net+authority+cm` |

## Import

Reservation groups can be imported using their UUID:

```shell
terraform import cloudlab_resgroup.compute <resgroup-uuid>
```
