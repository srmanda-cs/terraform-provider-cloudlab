data "cloudlab_experiment" "existing" {
  id = "550e8400-e29b-41d4-a716-446655440000"
}

output "experiment_name" {
  value = data.cloudlab_experiment.existing.name
}

output "experiment_status" {
  value = data.cloudlab_experiment.existing.status
}
