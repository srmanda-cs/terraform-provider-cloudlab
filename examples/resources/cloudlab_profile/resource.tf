resource "cloudlab_profile" "my_cluster" {
  name    = "my-cluster-profile"
  project = "MyProject"
  script  = file("${path.module}/profile.py")
}

# Or using a git repository:
# resource "cloudlab_profile" "from_repo" {
#   name           = "repo-profile"
#   project        = "MyProject"
#   repository_url = "https://github.com/example/cloudlab-profile.git"
# }
