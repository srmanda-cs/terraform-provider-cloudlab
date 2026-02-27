---
page_title: "cloudlab_node Data Source - terraform-provider-cloudlab"
description: |-
  Queries a specific node in a running CloudLab experiment by its client ID.
---

# cloudlab_node (Data Source)

Queries a specific node in a running CloudLab experiment by its client ID. Returns detailed node status including hostname, IPv4 address, operational state, and startup service status.

Use this data source to obtain a specific node's hostname or IP address — for example, to use as a connection target in a `null_resource` provisioner, or to pass to another system.

For retrieving all nodes at once, use [`cloudlab_manifest`](cloudlab_manifest.md) instead.

## Example Usage

### Get a node's hostname and IP

```terraform
resource "cloudlab_experiment" "cluster" {
  name            = "my-cluster"
  project         = "MyProject"
  profile_name    = "small-lan"
  profile_project = "CloudLab"
}

data "cloudlab_node" "server" {
  experiment_id = cloudlab_experiment.cluster.id
  client_id     = "server"
}

output "server_hostname" {
  value = data.cloudlab_node.server.hostname
}

output "server_ipv4" {
  value = data.cloudlab_node.server.ipv4
}
```

### Use hostname in a remote-exec provisioner

```terraform
data "cloudlab_node" "node1" {
  experiment_id = cloudlab_experiment.cluster.id
  client_id     = "node1"
}

resource "null_resource" "configure" {
  connection {
    type  = "ssh"
    host  = data.cloudlab_node.node1.hostname
    user  = "your-cloudlab-username"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get update",
      "sudo apt-get install -y docker.io",
    ]
  }
}
```

## Schema

### Required

- `experiment_id` (String) — The UUID of the running experiment containing the node.
- `client_id` (String) — The logical name of the node within the experiment, as defined in the profile (e.g., `node1`, `server`, `worker-0`).

### Read-Only

- `urn` (String) — The URN of the physical node assigned by CloudLab.
- `hostname` (String) — The fully qualified hostname of the node (e.g., `node1.my-cluster.MyProject.cloudlab.us`).
- `ipv4` (String) — The management IPv4 address of the node.
- `status` (String) — The current status of the node (e.g., `ready`, `booting`).
- `state` (String) — The current state of the node within CloudLab's state machine.
- `rawstate` (String) — The raw low-level state of the node.
- `startup_status` (String) — The current status of the startup script execution service on the node.
