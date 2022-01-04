data "buddy_group" "devs" {
  domain   = "mydomain"
  group_id = 1234
}

data "buddy_group" "admins" {
  domain = "mydomain"
  name   = "admins"
}