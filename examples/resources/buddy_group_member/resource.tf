resource "buddy_group_member" "john_in_devs" {
  domain    = "mydomain"
  group_id  = 1234
  member_id = 5678
}

resource "buddy_group_member" "sylvia_manager_in_devs" {
  domain    = "mydomain"
  group_id  = 1234
  member_id = 9999
  status    = "MANAGER"
}