data "buddy_rulesets" "all" {
  domain       = "mydomain"
  project_name = "myproject"
}

data "buddy_rulesets" "protect" {
  domain           = "mydomain"
  project_name     = "myproject"
  identifier_regex = "^protect-"
}
