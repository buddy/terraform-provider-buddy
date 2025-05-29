# Get all targets in workspace
data "buddy_targets" "all" {
  domain = "myworkspace"
}

# Get all targets in project
data "buddy_targets" "project_targets" {
  domain       = "myworkspace"
  project_name = "my-project"
}

# Get targets matching a regex pattern
data "buddy_targets" "prod_targets" {
  domain     = "myworkspace"
  name_regex = "^prod-"
}