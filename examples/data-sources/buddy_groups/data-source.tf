data "buddy_groups" "all" {
  domain = "mydomain"
}

data "buddy_groups" "filter" {
  domain = "mydomain"
  name_regex = "devs"
}




