data "buddy_workspaces" "all" {
}

data "buddy_workspaces" "name_started" {
  name_regex = "^started"
}

data "buddy_workspaces" "domain_ended" {
  domain_regex = "ended$"
}