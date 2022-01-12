data "buddy_permissions" "all" {
  domain = "mydomain"
}

data "buddy_permissions" "filter" {
  domain = "mydomain"
  type   = "DEVELOPER"
}




