resource "buddy_permission" "testers" {
  domain                  = "mydomain"
  name                    = "testers"
  pipeline_access_level   = "RUN_ONLY"
  repository_access_level = "READ_ONLY"
  sandbox_access_level    = "DENIED"
}