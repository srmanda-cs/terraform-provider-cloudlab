resource "cloudlab_resgroup" "hw_reservation" {
  project    = "MyProject"
  reason     = "Weekly experiment run requiring guaranteed xl170 hardware"
  expires_at = "2026-03-01T00:00:00Z"

  node_types = [
    {
      node_type = "xl170"
      aggregate = "utah.cloudlab.us"
      count     = 4
    }
  ]
}
