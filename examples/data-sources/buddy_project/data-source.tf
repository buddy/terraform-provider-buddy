data "buddy_project" "by_name" {
  domain = "mydomain"
  name   = "myproject"
}

data "buddy_project" "by_display_name" {
  domain       = "mydomain"
  display_name = "My project"
}



