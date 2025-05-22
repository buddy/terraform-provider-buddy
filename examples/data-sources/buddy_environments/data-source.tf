data "buddy_environments" "all" {
  domain       = "mydomain"
  project_name = "myproject"
}

data "buddy_environments" "with_name_dev" {
  domain       = "mydomain"
  project_name = "myproject"
  name_regex   = "^dev"
}
