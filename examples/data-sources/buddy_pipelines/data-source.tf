data "buddy_pipelines" "all" {
  domain       = "mydomain"
  project_name = "myproject"
}

data "buddy_pipelines" "with_name_started" {
  domain       = "mydomain"
  project_name = "myproject"
  name_regex   = "^started"
}

data "buddy_pipelines" "with_name_ended" {
  domain       = "mydomain"
  project_name = "myproject"
  name_regex   = "ended$"
}
