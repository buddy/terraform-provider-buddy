resource "buddy_project_group" "devs_in_test" {
  domain        = "mydomain"
  project_name  = "test"
  group_id      = 1234
  permission_id = 5678
}