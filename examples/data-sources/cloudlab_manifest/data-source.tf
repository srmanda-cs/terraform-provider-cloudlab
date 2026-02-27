data "cloudlab_manifest" "cluster_nodes" {
  experiment_id = cloudlab_experiment.cluster.id
}

output "node_hostnames" {
  value = [
    for manifest in data.cloudlab_manifest.cluster_nodes.manifests :
    [for node in manifest.nodes : node.hostname]
  ]
}
