# Get all targets in workspace
data "buddy_targets" "all" {
  domain = "myworkspace"
}

# Get targets in a specific project
data "buddy_targets" "project" {
  domain       = "myworkspace"
  project_name = "my-project"
}

# Get targets in a specific pipeline
data "buddy_targets" "pipeline" {
  domain       = "myworkspace"
  project_name = "my-project"
  pipeline_id  = 12345
}

# Get targets in a specific environment and filter by name
data "buddy_targets" "environment" {
  domain         = "myworkspace"
  project_name   = "my-project"
  environment_id = "env123"
  name_regex     = "^myname"
}