---
page_title: "cloudlab_resgroup Data Source - terraform-provider-cloudlab"
description: |-
  Queries an existing CloudLab reservation group by its UUID.
---

# cloudlab_resgroup (Data Source)

Queries an existing CloudLab reservation group by its UUID. Use this data source to reference reservation groups that were created outside of Terraform or in a separate Terraform state.

## Example Usage

### Look up a reservation group

```terraform
data "cloudlab_resgroup" "existing" {
  id = "a194e2be-1e5b-4617-84de-c4966cb5c578"
}

output "resgroup_expires_at" {
  value = data.cloudlab_resgroup.existing.expires_at
}

output "resgroup_node_types" {
  value = data.cloudlab_resgroup.existing.node_types
}
```

## Schema

### Required

- `id` (String) — The unique identifier (UUID) of the reservation group to look up.

### Read-Only

- `project` (String) — The CloudLab project this reservation group belongs to.
- `group` (String) — The project subgroup (nullable).
- `reason` (String) — The reason the reservation was created.
- `creator` (String) — The CloudLab username who created the reservation group.
- `created_at` (String) — The timestamp when the reservation group was created (ISO 8601, nullable).
- `start_at` (String) — The time the reservation starts (ISO 8601, nullable).
- `expires_at` (String) — The time the reservation expires (ISO 8601, nullable).
- `powder_zones` (String) — The Powder zone for radio resource reservations (nullable).
- `node_types` (List of Object) — The list of node type reservations. Each entry contains:
  - `urn` (String) — The aggregate URN of the reservation.
  - `node_type` (String) — The hardware node type reserved.
  - `count` (Number) — The number of nodes reserved.
- `ranges` (List of Object) — The list of frequency range reservations. Each entry contains:
  - `min_freq` (Number) — The start of the frequency range (inclusive) in MHz.
  - `max_freq` (Number) — The end of the frequency range (inclusive) in MHz.
- `routes` (List of Object) — The list of named route reservations. Each entry contains:
  - `name` (String) — The route name reserved.
