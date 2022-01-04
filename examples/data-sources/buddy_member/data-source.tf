data "buddy_member" "by_name" {
  domain = "mydomain"
  name   = "John Doe"
}

data "buddy_member" "by_email" {
  domain = "mydomain"
  email  = "john@doe.com"
}

data "buddy_member" "by_id" {
  domain    = "mydomain"
  member_id = 123456
}





