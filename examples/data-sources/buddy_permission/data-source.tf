data "buddy_permission" "by_name" {
  domain = "mydomain"
  name   = "my perm"
}

data "buddy_permission" "by_id" {
  domain        = "mydomain"
  permission_id = 123456
}




