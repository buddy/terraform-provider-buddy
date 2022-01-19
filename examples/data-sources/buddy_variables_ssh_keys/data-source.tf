data "buddy_variables_ssh_keys" "workspace" {
  domain = "mydomain"
}

data "buddy_variables_ssh_keys" "search" {
  domain = "mydomain"
  key_regex = "^mykey"
}

data "buddy_variables_ssh_keys" "project" {
  domain = "mydomain"
  project_name = "myproject"
}

data "buddy_variables_ssh_keys" "pipeline" {
  domain = "mydomain"
  pipeline_id = 123456
}




