data "buddy_projects" "all" {
  domain = "mydomain"
}

data "buddy_projects" "only_mine" {
  domain     = "mydomain"
  membership = true
}

data "buddy_projects" "only_active" {
  domain = "mydomain"
  status = "ACTIVE"
}

data "buddy_projects" "with_name_started" {
  domain     = "mydomain"
  name_regex = "^started"
}

data "buddy_projects" "with_display_name_ended" {
  domain             = "mydomain"
  display_name_regex = "ended$"
}



