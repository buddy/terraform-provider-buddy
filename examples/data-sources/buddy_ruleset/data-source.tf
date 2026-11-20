data "buddy_ruleset" "by_id" {
  domain       = "mydomain"
  project_name = "myproject"
  ruleset_id   = "4GP6nQPa"
}

data "buddy_ruleset" "by_identifier" {
  domain       = "mydomain"
  project_name = "myproject"
  identifier   = "protect-main"
}
