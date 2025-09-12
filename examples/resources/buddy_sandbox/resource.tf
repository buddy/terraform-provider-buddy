resource "buddy_sandbox" "sb" {
  domain              = "mydomain"
  project_name        = "test"
  name                = "sb"
  install_commands    = "apt-get update && apt-get install -y curl"
  wait_for_running    = true
  wait_for_configured = true
  endpoints = {
    "local" = {
      endpoint = "localhost:3333"
      type     = "HTTP"
    }
  }
}