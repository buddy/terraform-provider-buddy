data "buddy_domains" "all" {
  workspace_domain = "myworkspace"
}

data "buddy_domains" "private" {
  workspace_domain = "myworkspace"
  type             = "PRIVATE"
}

data "buddy_domains" "filter" {
  workspace_domain = "myworkspace"
  domain_regex     = "\\.com$"
}
