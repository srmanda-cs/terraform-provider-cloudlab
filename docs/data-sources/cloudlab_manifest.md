---
page_title: "cloudlab_manifest Data Source - terraform-provider-cloudlab"
description: |-
  Retrieves the manifests for a running CloudLab experiment, including node hostnames and IP addresses.
---

# cloudlab_manifest (Data Source)

Retrieves the manifests for a running CloudLab experiment. The manifest contains the assigned node hostnames, IP addresses, and network interfaces for all nodes in the experiment.

Use this data source to obtain the hostnames or IPs of your provisioned nodes so you can, for example, configure them with a subsequent provisioner or pass them to another tool.

## Example Usage

```terraform
resource "cloudlab_experiment" "cluster" {
  name            = "my-cluster"
  project         = "MyProject"
  profile_name    = "small-lan"
  profile_project = "CloudLab"
}

data "cloudlab_manifest" "nodes" {
  experiment_id = cloudlab_experiment.cluster.id
}

output "node_hostnames" {
  value = flatten([
    for manifest in data.cloudlab_manifest.nodes.manifests :
    [for node in manifest.nodes : node.hostname]
  ])
}
```

## Schema

### Required

- `experiment_id` (String) — The UUID of the running experiment to retrieve manifests for.

### Read-Only

- `manifests` (List of Object) — The list of manifests, one per CloudLab aggregate/site. Each entry contains:
  - `aggregate` (String) — The CloudLab aggregate (site) this manifest applies to.
  - `nodes` (List of Object) — The list of nodes provisioned at this aggregate. Each node contains:
    - `client_id` (String) — The client-assigned node identifier from the profile.
    - `hostname` (String) — The fully qualified hostname of the node.
    - `interfaces` (List of Object) — The network interfaces on this node. Each interface contains:
      - `name` (String) — The interface name (e.g., `eth0`).
      - `address` (String) — The IP address assigned to this interface.
