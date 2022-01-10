resource "buddy_project_member" "john_in_test" {
  domain        = "mydomain"
  project_name  = "test"
  member_id     = 1234
  permission_id = 5678
}