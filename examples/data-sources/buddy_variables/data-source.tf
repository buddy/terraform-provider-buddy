data "buddy_variables" "workspace" {
  domain = "mydomain"
}

data "buddy_variables" "search" {
  domain    = "mydomain"
  key_regex = "^mykey"
}

data "buddy_variables" "project" {
  domain       = "mydomain"
  project_name = "myproject"
}

data "buddy_variables" "pipeline" {
  domain      = "mydomain"
  pipeline_id = 123456
}




