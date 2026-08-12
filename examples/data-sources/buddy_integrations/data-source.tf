data "buddy_integrations" "all" {
  domain = "mydomain"
}

data "buddy_integrations" "amazon" {
  domain = "mydomain"
  type   = "AMAZON"
}

data "buddy_integrations" "project" {
  domain       = "mydomain"
  project_name = "myproject"
}
