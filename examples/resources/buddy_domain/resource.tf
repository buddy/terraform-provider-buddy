resource "buddy_domain" "dev" {
  workspace_domain = "myworkspace"
  domain           = "test.com"
}

resource "buddy_domain" "internal" {
  workspace_domain = "myworkspace"
  domain           = "internal.lan"
  type             = "PRIVATE"
}
