# Get target by name
data "buddy_target" "by_name" {
  domain = "myworkspace"
  name   = "Production Server"
}

# Get target by ID
data "buddy_target" "by_id" {
  domain    = "myworkspace"
  target_id = 12345
}

# Get target from project
data "buddy_target" "project_target" {
  domain       = "myworkspace"
  project_name = "my-project"
  name         = "Staging Server"
}