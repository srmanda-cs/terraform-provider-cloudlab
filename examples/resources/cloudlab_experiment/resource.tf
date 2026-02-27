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

output "experiment_status" {
  value = cloudlab_experiment.cluster.status
}
